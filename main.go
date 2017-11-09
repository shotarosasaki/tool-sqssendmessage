package main

import (
	"time"

	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func main() {
	s, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1"),
	})
	if err != nil {
		panic(err)
	}
	svc := sqs.New(s, aws.NewConfig().WithEndpoint("http://localhost:9324/queue/").WithRegion("ap-northeast-1"))

	bodyTmpl := `{
    "replyToken": "%v%v",
    "type": "message",
    "timestamp": 1462629479859,
    "source": {
      "type": "user",
      "userId": "U206d25c2ea6bd87c17655609a1c37cb8"
    },
    "message": {
      "id": "325708",
      "type": "text",
      "text": "Hello, world"
    }
  }`

	for i := 0; i < 1000; i++ {
		fmt.Println(fmt.Sprintf("No.%v", i))
		res, err := svc.SendMessageBatch(&sqs.SendMessageBatchInput{
			QueueUrl: aws.String("http://localhost:9324/queue/local_line_messages"),
			Entries: []*sqs.SendMessageBatchRequestEntry{

				&sqs.SendMessageBatchRequestEntry{
					Id: aws.String(fmt.Sprintf("%v%v", "A", i)),
					MessageAttributes: map[string]*sqs.MessageAttributeValue{
						"acid":     &sqs.MessageAttributeValue{DataType: aws.String("String"), StringValue: aws.String("736afbdc388752eb")},
						"sendType": &sqs.MessageAttributeValue{DataType: aws.String("Number"), StringValue: aws.String("1")},
					},
					MessageBody: aws.String(fmt.Sprintf(bodyTmpl, "A", i)),
				},

				&sqs.SendMessageBatchRequestEntry{
					Id: aws.String(fmt.Sprintf("%v%v", "B", i)),
					MessageAttributes: map[string]*sqs.MessageAttributeValue{
						"acid":     &sqs.MessageAttributeValue{DataType: aws.String("String"), StringValue: aws.String("736afbdc388752eb")},
						"sendType": &sqs.MessageAttributeValue{DataType: aws.String("Number"), StringValue: aws.String("1")},
					},
					MessageBody: aws.String(fmt.Sprintf(bodyTmpl, "B", i)),
				},

				&sqs.SendMessageBatchRequestEntry{
					Id: aws.String(fmt.Sprintf("%v%v", "C", i)),
					MessageAttributes: map[string]*sqs.MessageAttributeValue{
						"acid":     &sqs.MessageAttributeValue{DataType: aws.String("String"), StringValue: aws.String("736afbdc388752eb")},
						"sendType": &sqs.MessageAttributeValue{DataType: aws.String("Number"), StringValue: aws.String("1")},
					},
					MessageBody: aws.String(fmt.Sprintf(bodyTmpl, "C", i)),
				},

				&sqs.SendMessageBatchRequestEntry{
					Id: aws.String(fmt.Sprintf("%v%v", "D", i)),
					MessageAttributes: map[string]*sqs.MessageAttributeValue{
						"acid":     &sqs.MessageAttributeValue{DataType: aws.String("String"), StringValue: aws.String("736afbdc388752eb")},
						"sendType": &sqs.MessageAttributeValue{DataType: aws.String("Number"), StringValue: aws.String("1")},
					},
					MessageBody: aws.String(fmt.Sprintf(bodyTmpl, "D", i)),
				},

				&sqs.SendMessageBatchRequestEntry{
					Id: aws.String(fmt.Sprintf("%v%v", "E", i)),
					MessageAttributes: map[string]*sqs.MessageAttributeValue{
						"acid":     &sqs.MessageAttributeValue{DataType: aws.String("String"), StringValue: aws.String("736afbdc388752eb")},
						"sendType": &sqs.MessageAttributeValue{DataType: aws.String("Number"), StringValue: aws.String("1")},
					},
					MessageBody: aws.String(fmt.Sprintf(bodyTmpl, "E", i)),
				},

				&sqs.SendMessageBatchRequestEntry{
					Id: aws.String(fmt.Sprintf("%v%v", "F", i)),
					MessageAttributes: map[string]*sqs.MessageAttributeValue{
						"acid":     &sqs.MessageAttributeValue{DataType: aws.String("String"), StringValue: aws.String("736afbdc388752eb")},
						"sendType": &sqs.MessageAttributeValue{DataType: aws.String("Number"), StringValue: aws.String("1")},
					},
					MessageBody: aws.String(fmt.Sprintf(bodyTmpl, "F", i)),
				},

				&sqs.SendMessageBatchRequestEntry{
					Id: aws.String(fmt.Sprintf("%v%v", "G", i)),
					MessageAttributes: map[string]*sqs.MessageAttributeValue{
						"acid":     &sqs.MessageAttributeValue{DataType: aws.String("String"), StringValue: aws.String("736afbdc388752eb")},
						"sendType": &sqs.MessageAttributeValue{DataType: aws.String("Number"), StringValue: aws.String("1")},
					},
					MessageBody: aws.String(fmt.Sprintf(bodyTmpl, "G", i)),
				},

				&sqs.SendMessageBatchRequestEntry{
					Id: aws.String(fmt.Sprintf("%v%v", "H", i)),
					MessageAttributes: map[string]*sqs.MessageAttributeValue{
						"acid":     &sqs.MessageAttributeValue{DataType: aws.String("String"), StringValue: aws.String("736afbdc388752eb")},
						"sendType": &sqs.MessageAttributeValue{DataType: aws.String("Number"), StringValue: aws.String("1")},
					},
					MessageBody: aws.String(fmt.Sprintf(bodyTmpl, "H", i)),
				},

				&sqs.SendMessageBatchRequestEntry{
					Id: aws.String(fmt.Sprintf("%v%v", "I", i)),
					MessageAttributes: map[string]*sqs.MessageAttributeValue{
						"acid":     &sqs.MessageAttributeValue{DataType: aws.String("String"), StringValue: aws.String("736afbdc388752eb")},
						"sendType": &sqs.MessageAttributeValue{DataType: aws.String("Number"), StringValue: aws.String("1")},
					},
					MessageBody: aws.String(fmt.Sprintf(bodyTmpl, "I", i)),
				},

				&sqs.SendMessageBatchRequestEntry{
					Id: aws.String(fmt.Sprintf("%v%v", "J", i)),
					MessageAttributes: map[string]*sqs.MessageAttributeValue{
						"acid":     &sqs.MessageAttributeValue{DataType: aws.String("String"), StringValue: aws.String("736afbdc388752eb")},
						"sendType": &sqs.MessageAttributeValue{DataType: aws.String("Number"), StringValue: aws.String("1")},
					},
					MessageBody: aws.String(fmt.Sprintf(bodyTmpl, "J", i)),
				},
			},
		})
		if err != nil {
			panic(err)
		}
		if res == nil {
			panic("res is nil")
		}
		time.Sleep(10 * time.Millisecond)
		fmt.Println("Next!")
	}
}
