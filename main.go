package main

import (
	"fmt"

	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	uuid "github.com/satori/go.uuid"
	"flag"
)

func main() {
	tformat := time.StampMicro
	region := "ap-northeast-1"

	argACID := flag.String("acid", "1234abcd5678efgh", "ACID")
	argUserID := flag.String("userID", "U1234567890abcdefgh1234567890abcd", "ユーザID")
	argEndpoint := flag.String("endpoint", "http://localhost:9324/queue/", "SQSエンドポイント")
	argQueueName := flag.String("queueName", "local_line_messages", "SQSキュー名")
	argLoopCount := flag.Int("loopCount", 100, "ループカウント")
	argDelay := flag.String("delay", "1s", "１ループ毎の遅延時間（「500ms」のようなtime.Duration形式）")
	argEntryCount := flag.Int("entryCount", 10, "１ループ毎の送信メッセージ数")
	flag.Parse()

	fmt.Printf("[%v][ACID:%v]ARGS(userID:%v)\n", time.Now().Format(tformat), *argACID, *argUserID)
	fmt.Printf("[%v][ACID:%v]ARGS(endpoint:%v)\n", time.Now().Format(tformat), *argACID, *argEndpoint)
	fmt.Printf("[%v][ACID:%v]ARGS(queueName:%v)\n", time.Now().Format(tformat), *argACID, *argQueueName)
	fmt.Printf("[%v][ACID:%v]ARGS(loopCount:%v)\n", time.Now().Format(tformat), *argACID, *argLoopCount)
	fmt.Printf("[%v][ACID:%v]ARGS(delay:%v)\n", time.Now().Format(tformat), *argACID, *argDelay)
	fmt.Printf("[%v][ACID:%v]ARGS(entryCount:%v)\n", time.Now().Format(tformat), *argACID, *argEntryCount)

	// Credentialは環境変数セット済の前提
	awsCfg := &aws.Config{
		Region: aws.String(region),
		Endpoint: argEndpoint,
	}

	s, err := session.NewSession(awsCfg)
	if err != nil {
		panic(err)
	}
	svc := sqs.New(s)

	bodyTmpl := `{
      "replyToken": "%v%v",
      "type": "message",
      "timestamp": 1462629479859,
      "source": {
        "type": "user",
        "userId": "%v"
      },
      "message": {
        "id": "325708",
        "type": "text",
        "text": "Hello, world"
      }
    }`

    fmt.Printf("[%v][ACID:%v]START\n", time.Now().Format(tformat), *argACID)

	for i := 0; i < *argLoopCount; i++ {
		fmt.Printf("[%v][ACID:%v]No.%v\n", time.Now().Format(tformat), *argACID, i)

		var entries []*sqs.SendMessageBatchRequestEntry
		for j := 0; j < *argEntryCount; j++ {
			entries = append(entries, &sqs.SendMessageBatchRequestEntry{
				Id: aws.String(fmt.Sprintf("%v%v", i, j)),
				MessageAttributes: map[string]*sqs.MessageAttributeValue{
					"uniqueMessageID": &sqs.MessageAttributeValue{DataType: aws.String("String"), StringValue: aws.String(uuid.NewV4().String())},
					"acid":            &sqs.MessageAttributeValue{DataType: aws.String("String"), StringValue: aws.String(*argACID)},
					"sendType":        &sqs.MessageAttributeValue{DataType: aws.String("Number"), StringValue: aws.String("1")},
				},
				MessageBody: aws.String(fmt.Sprintf(bodyTmpl, i, j, *argUserID)),
			},)
		}

		res, err := svc.SendMessageBatch(&sqs.SendMessageBatchInput{
			QueueUrl: aws.String(fmt.Sprintf("%s%s", *argEndpoint, *argQueueName)),
			Entries: entries,
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
