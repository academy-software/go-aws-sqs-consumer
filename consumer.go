package consumer

import (
	"log"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

//Consumer holds the consumer data
type Consumer struct {
	queueURL string
	channel  chan *sqs.Message
	handler  func(m *sqs.Message) error
	config   *Config
}

//Config holds the configuration for consuming and processing the queue
type Config struct {
	MaxNumberOfMessages      int64
	NumberOfMessageReceivers int
	VisibilityTimeout        int64
	ReceiverTimeout          int
}

var sess *session.Session

func init() {
	sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
}

//New creates a new Queue consumer
func New(queueURL string, handler func(m *sqs.Message) error, config *Config) Consumer {
	c := make(chan *sqs.Message)
	return Consumer{
		queueURL,
		c,
		handler,
		config,
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
	r := Receiver{
		queueURL:            c.queueURL,
		channel:             c.channel,
		sess:                sess,
		visibilityTimeout:   c.config.VisibilityTimeout,
		receiverTimeout:     c.config.ReceiverTimeout,
		maxNumberOfMessages: c.config.MaxNumberOfMessages,
	}

	for i := 0; i < c.config.NumberOfMessageReceivers; i++ {
		go r.receiveMessages()
	}
}

// startProcessor starts a goroutine to handle each message from channel
func (c *Consumer) startProcessor() {
	p := Processor{
		queueURL: c.queueURL,
		channel:  c.channel,
		sess:     sess,
		handler:  c.handler,
	}

	for m := range c.channel {
		go p.processMessage(m)
	}
}
