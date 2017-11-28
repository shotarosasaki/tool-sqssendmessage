package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"time"

	"flag"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
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
	argThroughput := flag.Int("throughput", 10, "１秒当たりのSQS投入メッセージ数")
	argDuration := flag.Int("duration", 5, "継続時間（単位：秒）")

	argEndpoint := flag.String("endpoint", "http://localhost:9324/queue/", "SQSエンドポイント(ローカルでElasticMQ起動した際のエンドポイント指定用)")

	flag.Parse()

	fmt.Printf("[%v][ACID:%v]ARGS(userID:%v)\n", time.Now().Format(tformat), *argACID, *argUserID)
	fmt.Printf("[%v][ACID:%v]ARGS(queueName:%v)\n", time.Now().Format(tformat), *argACID, *argQueueName)
	fmt.Printf("[%v][ACID:%v]ARGS(throughput:%v)\n", time.Now().Format(tformat), *argACID, *argThroughput)
	fmt.Printf("[%v][ACID:%v]ARGS(duration:%v)\n", time.Now().Format(tformat), *argACID, *argDuration)

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

	loopCount := *argThroughput / 10
	fmt.Printf("[%v][ACID:%v]loopCount:%v\n", time.Now().Format(tformat), *argACID, loopCount)

	// コストのかかるものは事前に準備
	entriesentriesentries := make(map[int]map[int][]*sqs.SendMessageBatchRequestEntry)
	for k := 0; k < *argDuration; k++ {
		entriesentries := make(map[int][]*sqs.SendMessageBatchRequestEntry)
		for i := 0; i < loopCount; i++ {
			var entries []*sqs.SendMessageBatchRequestEntry
			for j := 0; j < 10; j++ {
				entries = append(entries, &sqs.SendMessageBatchRequestEntry{
					Id: aws.String(fmt.Sprintf("%v%v%v", k, i, j)),
					MessageAttributes: map[string]*sqs.MessageAttributeValue{
						"uniqueMessageID": &sqs.MessageAttributeValue{DataType: aws.String("String"), StringValue: aws.String(uuid.NewV4().String())},
						"acid":            &sqs.MessageAttributeValue{DataType: aws.String("String"), StringValue: aws.String(*argACID)},
						"sendType":        &sqs.MessageAttributeValue{DataType: aws.String("Number"), StringValue: aws.String("2")},
					},
					MessageBody: aws.String(
						fmt.Sprintf(`{"to":"U9999999999aaaaaaaaaa9999999999bbbb","messages":[{"type":"text","text":"SQS投入負荷ツール：No.%v-%v-%v"}]}`, k, i, j),
					),
				})
			}
			entriesentries[i] = entries
		}
		entriesentriesentries[k] = entriesentries
	}

	cancelCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	d, err := time.ParseDuration("1s")
	if err != nil {
		panic(err)
	}

	kCnt := 0
	ticker := time.NewTicker(d)
	for {
		select {
		case t := <-ticker.C:
			if kCnt >= *argDuration {
				fmt.Printf("[%v][ACID:%v]END\n", time.Now().Format(tformat), *argACID)
				return
			}
			fmt.Printf("[%v][ACID:%v]tick\n", t.Format(tformat), *argACID)

			for i := 0; i < loopCount; i++ {
				fmt.Printf("[%v][ACID:%v]No.%v-%v\n", time.Now().Format(tformat), *argACID, kCnt, i)
				go func(k, i int, entriesentriesentries map[int]map[int][]*sqs.SendMessageBatchRequestEntry, cancel context.CancelFunc) {
					res, err := svc.SendMessageBatch(&sqs.SendMessageBatchInput{
						QueueUrl: queueUrl,
						Entries:  entriesentriesentries[k][i],
					})
					if err != nil {
						fmt.Printf("[%v][ACID:%v]SendMessageBatch Error:%#v\n", time.Now().Format(tformat), *argACID, err.Error())
						cancel()
					}
					if res == nil {
						cancel()
					}
				}(kCnt, i, entriesentriesentries, cancel)
			}

			kCnt = kCnt + 1
		case <-cancelCtx.Done():
			fmt.Printf("[%v][ACID:%v]END2\n", time.Now().Format(tformat), *argACID)
			return
		}
	}

	fmt.Printf("[%v][ACID:%v]END\n", time.Now().Format(tformat), *argACID)
}
