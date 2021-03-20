package consumer

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
	"os"
	"time"
)

//SqsReceiver defines the struct that polls messages from AWS SQS
type SqsReceiver struct {
	queueURL                string
	channel                 chan *sqs.Message
	shutdown                chan os.Signal
	sess                    *session.Session
	visibilityTimeout       int64
	maxNumberOfMessages     int64
	pollDelayInMilliseconds int
}

func (r *SqsReceiver) applyBackPressure() {
	time.Sleep(time.Millisecond * time.Duration(r.pollDelayInMilliseconds))
}

func (r *SqsReceiver) receiveMessages() {
	queue := sqs.New(r.sess)
	for {

		select {
		case <-r.shutdown:
			log.Println("Shutting down message receiver")
			close(r.channel)
			return
		default:
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
}
