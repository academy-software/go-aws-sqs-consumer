package consumer

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
	"time"
)

//Receiver defines the struct that polls messages from AWS SQS
type Receiver struct {
	queueURL                string
	channel                 chan *sqs.Message
	sess                    *session.Session
	visibilityTimeout       int64
	maxNumberOfMessages     int64
	pollDelayInMilliseconds int
}

func (r *Receiver) applyBackPressure() {
	time.Sleep(time.Millisecond * time.Duration(r.pollDelayInMilliseconds))
}

func (r *Receiver) receiveMessages() {
	queue := sqs.New(r.sess)
	for {
		msgResult, err := queue.ReceiveMessage(&sqs.ReceiveMessageInput{
			AttributeNames: []*string{
				aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
			},
			MessageAttributeNames: []*string{
				aws.String(sqs.QueueAttributeNameAll),
			},
			QueueUrl:            aws.String(r.queueURL),
			MaxNumberOfMessages: aws.Int64(r.maxNumberOfMessages),
			VisibilityTimeout:   aws.Int64(r.visibilityTimeout),
		})

		if err != nil {
			log.Println("Could not read from queue", err)
			return
		}

		if len(msgResult.Messages) > 0 {
			for _, m := range msgResult.Messages {
				r.channel <- m
			}
		}

		r.applyBackPressure()
	}
}
