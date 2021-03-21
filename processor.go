package consumer

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/uuid"
	"log"
	"sync"
)

//Processor defines the struct that processes messages from AWS SQS
type Processor struct {
	queueURL string
	queue    *sqs.SQS
	handler  func(*sqs.Message) error
}

func (p *Processor) processMessages(messages []*sqs.Message) {
	nMessages := len(messages)
	deleteChannel := make(chan *string, nMessages)
	wg := sync.WaitGroup{}
	wg.Add(nMessages)

	for _, m := range messages {
		go func(message *sqs.Message) {
			defer wg.Done()
			err := p.handler(message)
			if err != nil {
				log.Println("Error while handling message:", err)
				return
			}
			deleteChannel <- message.ReceiptHandle
		}(m)
	}

	wg.Wait()

	close(deleteChannel)
	entries := make([]*sqs.DeleteMessageBatchRequestEntry, 0, nMessages)

	for receipt := range deleteChannel {
		entries = append(entries, &sqs.DeleteMessageBatchRequestEntry{
			Id:            aws.String(uuid.NewString()),
			ReceiptHandle: receipt,
		})
	}

	if len(entries) > 0 {
		_, dErr := p.queue.DeleteMessageBatch(&sqs.DeleteMessageBatchInput{
			QueueUrl: &p.queueURL,
			Entries:  entries,
		})
		if dErr != nil {
			log.Println("Failed while trying to delete message:", dErr)
		}
	}
}
