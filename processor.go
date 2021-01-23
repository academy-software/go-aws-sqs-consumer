package consumer

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
)

//Processor defines the struct that processes messages from AWS SQS
type Processor struct {
	queueURL string
	channel  chan *sqs.Message
	sess     *session.Session
	handler  func(*sqs.Message) error
}

func (p *Processor) processMessage(m *sqs.Message) {
	err := p.handler(m)
	if err != nil {
		log.Println("Error while handling message:", err)
		return
	}

	queue := sqs.New(sess)

	_, derr := queue.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(p.queueURL),
		ReceiptHandle: m.ReceiptHandle,
	})

	if derr != nil {
		log.Println("Failed while trying to delete message:", derr)
	}
}
