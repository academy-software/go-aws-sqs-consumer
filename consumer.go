package consumer

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
	"os"
	"os/signal"
	"syscall"
)

//Consumer holds the consumer data
type Consumer struct {
	queueURL        string
	messagesChannel chan []*sqs.Message
	handler         func(m *sqs.Message) error
	config          *Config
	receiver        SqsReceiver
}

//Config holds the configuration for consuming and processing the queue
type Config struct {
	AwsSession                  *session.Session
	SqsMaxNumberOfMessages      int64
	SqsMessageVisibilityTimeout int64
	Receivers                   int
	PollDelayInMilliseconds     int
}

var sess *session.Session

//New creates a new Queue consumer
func New(queueURL string, handler func(m *sqs.Message) error, config *Config) Consumer {
	sess = config.AwsSession
	c := make(chan []*sqs.Message)
	shutdown := make(chan os.Signal, 1)

	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	r := SqsReceiver{
		queueURL:                queueURL,
		messagesChannel:         c,
		shutdown:                shutdown,
		sess:                    sess,
		visibilityTimeout:       config.SqsMessageVisibilityTimeout,
		maxNumberOfMessages:     config.SqsMaxNumberOfMessages,
		pollDelayInMilliseconds: config.PollDelayInMilliseconds,
	}

	return Consumer{
		queueURL:        queueURL,
		messagesChannel: c,
		handler:         handler,
		config:          config,
		receiver:        r,
	}
}

//Start initiates the queue consumption process
func (c *Consumer) Start() {
	log.Println("Starting to consume", c.queueURL)
	c.startReceivers()
	c.startProcessor()
}

// startReceivers starts N (defined in NumberOfMessageReceivers) goroutines to poll messages from SQS
func (c *Consumer) startReceivers() {
	for i := 0; i < c.config.Receivers; i++ {
		go c.receiver.receiveMessages()
	}
}

// startProcessor starts a goroutine to handle each message from messagesChannel
func (c *Consumer) startProcessor() {
	queue := sqs.New(sess)

	p := Processor{
		queueURL: c.queueURL,
		queue:    queue,
		handler:  c.handler,
	}

	for messages := range c.messagesChannel {
		go p.processMessages(messages)
	}
}

//SetPollDelay increases time between message poll
func (c *Consumer) SetPollDelay(delayBetweenPoolsInMilliseconds int) {
	c.receiver.pollDelayInMilliseconds = delayBetweenPoolsInMilliseconds
}
