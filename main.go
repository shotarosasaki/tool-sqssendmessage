package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"time"

	"flag"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/line/line-bot-sdk-go/linebot"
	uuid "github.com/satori/go.uuid"
)

func createLocalSqsCli(argQueueName, argEndpoint *string) (*sqs.SQS, *string) {

	cfg := &aws.Config{
		Region:   aws.String("ap-northeast-1"),
		Endpoint: argEndpoint,
	}

	s, err := session.NewSession(cfg)
	if err != nil {
		panic(err)
	}

	svc := sqs.New(
		s,
		&aws.Config{
			MaxRetries: aws.Int(10)}) // Exponential Backoff リトライ回数

	queueUrl := aws.String(fmt.Sprintf("%v%v", *argEndpoint, *argQueueName))

	return svc, queueUrl
}

func createSqsCli(argQueueName *string) (*sqs.SQS, *string) {

	cfg := &aws.Config{
		Region: aws.String("ap-northeast-1"),
	}

	s, err := session.NewSession(cfg)
	if err != nil {
		panic(err)
	}

	svc := sqs.New(
		s,
		&aws.Config{
			MaxRetries: aws.Int(10)}) // Exponential Backoff リトライ回数

	resp, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: argQueueName,
	})
	if err != nil {
		panic(err)
	}
	queueUrl := resp.QueueUrl

	return svc, queueUrl
}

func main() {
	tformat := time.StampMicro

	argACID := flag.String("acid", "1234abcd5678efgh", "ACID")
	argUserID := flag.String("userID", "U1234567890abcdefgh1234567890abcd", "ユーザID")
	argQueueName := flag.String("queueName", "local_line_messages", "SQSキュー名")
	argLoopCount := flag.Int("loopCount", 1, "ループカウント")
	argDelay := flag.String("delay", "1s", "１ループ毎の遅延時間（「500ms」のようなtime.Duration形式）")
	argEntryCount := flag.Int("entryCount", 1, "１ループ毎の送信メッセージ数")

	argEndpoint := flag.String("endpoint", "http://localhost:9324/queue/", "SQSエンドポイント(ローカルでElasticMQ起動した際のエンドポイント指定用)")

	flag.Parse()

	fmt.Printf("[%v][ACID:%v]ARGS(userID:%v)\n", time.Now().Format(tformat), *argACID, *argUserID)
	fmt.Printf("[%v][ACID:%v]ARGS(queueName:%v)\n", time.Now().Format(tformat), *argACID, *argQueueName)
	fmt.Printf("[%v][ACID:%v]ARGS(loopCount:%v)\n", time.Now().Format(tformat), *argACID, *argLoopCount)
	fmt.Printf("[%v][ACID:%v]ARGS(delay:%v)\n", time.Now().Format(tformat), *argACID, *argDelay)
	fmt.Printf("[%v][ACID:%v]ARGS(entryCount:%v)\n", time.Now().Format(tformat), *argACID, *argEntryCount)

	fmt.Printf("[%v][ACID:%v]ARGS(endpoint:%v)\n", time.Now().Format(tformat), *argACID, *argEndpoint)

	if strings.HasPrefix(*argQueueName, "production_") {
		fmt.Println("本番用のキューには流せません。")
		os.Exit(-1)
	}

	var svc *sqs.SQS
	var queueUrl *string
	if strings.HasPrefix(*argQueueName, "local_") {
		svc, queueUrl = createLocalSqsCli(argQueueName, argEndpoint)
	} else {
		svc, queueUrl = createSqsCli(argQueueName)
	}

	fmt.Printf("[%v][ACID:%v]START\n", time.Now().Format(tformat), *argACID)
	fmt.Printf("[%v][ACID:%v]queueUrl:%v\n", time.Now().Format(tformat), *argACID, *queueUrl)

	// ディレイかけながら複数回のSQSバッチ送信を行う
	for i := 0; i < *argLoopCount; i++ {
		fmt.Printf("[%v][ACID:%v]No.%v\n", time.Now().Format(tformat), *argACID, i)

		var entries []*sqs.SendMessageBatchRequestEntry

		// SQSへのバッチ送信に渡すエントリーを指定件数分、生成
		for j := 0; j < *argEntryCount; j++ {

			// LINE-API[Push](https://developers.line.me/ja/docs/messaging-api/reference/#anchor-0c00cb0f42b970892f7c3382f92620dca5a110fc)のメッセージ形式に依存
			pushStruct := &struct {
				To       string            `json:"to"`
				Messages []linebot.Message `json:"messages"`
			}{
				To:       *argUserID,
				Messages: []linebot.Message{linebot.NewTextMessage(fmt.Sprintf("SQS投入負荷ツール：No.%v-%v", i, j))},
			}
			btEvent, err := json.Marshal(pushStruct)
			if err != nil {
				panic(err)
			}
			bodyTmpl := string(btEvent)

			entries = append(entries, &sqs.SendMessageBatchRequestEntry{
				Id: aws.String(fmt.Sprintf("%v%v", i, j)),
				MessageAttributes: map[string]*sqs.MessageAttributeValue{
					"uniqueMessageID": &sqs.MessageAttributeValue{DataType: aws.String("String"), StringValue: aws.String(uuid.NewV4().String())},
					"acid":            &sqs.MessageAttributeValue{DataType: aws.String("String"), StringValue: aws.String(*argACID)},
					"sendType":        &sqs.MessageAttributeValue{DataType: aws.String("Number"), StringValue: aws.String("2")},
				},
				MessageBody: aws.String(bodyTmpl),
			})
		}

		res, err := svc.SendMessageBatch(&sqs.SendMessageBatchInput{
			QueueUrl: queueUrl,
			Entries:  entries,
		})
		if err != nil {
			fmt.Printf("[%v][ACID:%v]SendMessageBatch Error:%#v\n", time.Now().Format(tformat), *argACID, err.Error())
			panic(err)
		}
		if res == nil {
			panic("res is nil")
		}

		delay, err := time.ParseDuration(*argDelay)
		if err != nil {
			fmt.Printf("[%v][ACID:%v]ParseDuration Error:%#v\n", time.Now().Format(tformat), *argACID, err.Error())
			panic(err)
		}

		time.Sleep(delay)
	}

	fmt.Printf("[%v][ACID:%v]END\n", time.Now().Format(tformat), *argACID)
}
