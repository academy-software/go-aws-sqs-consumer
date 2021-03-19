package consumer

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"math/rand"
	"syscall"
	"testing"
	"time"
)

func TestConsumeQueueEventSuccessfully(t *testing.T) {
	queueURL := "http://localhost:4566/000000000000/my-queue"

	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials("not", "empty", ""),
		DisableSSL:  aws.Bool(true),
		Region:      aws.String(endpoints.UsEast1RegionID),
		Endpoint:    aws.String("http://localhost:4566"),
	}))

	q := sqs.New(sess)
	sendMessage(q, queueURL)

	c := New(queueURL, func(m *sqs.Message) error {
		_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		return nil
	}, &Config{
		AwsSession:                  sess,
		SqsMaxNumberOfMessages:      10,
		SqsMessageVisibilityTimeout: 10,
		Receivers:                   1,
		PollDelayInMilliseconds:     100,
	})

	c.Start()
}

func sendMessage(queue *sqs.SQS, queueURL string) {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	_, err := queue.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(0),
		MessageBody:  aws.String(fmt.Sprintf("Message %d", r1.Intn(100))),
		QueueUrl:     &queueURL,
	})

	if err != nil {
		fmt.Println("Could not send message:", err)
	}
}
