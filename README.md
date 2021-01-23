# Go SQS Consumer

A package to easily create an AWS SQS message consumer

## Features

* Configurable number of goroutines that poll messages from SQS
* Backpressure
  - You can dinamically control the poll delay period by using the SetPollDelay function

## Usage

```golang
func main() {
	q := "http://localhost:4566/000000000000/my-queue"
	c := consumer.New(q, handle,
		&consumer.Config{
			Receivers:                   1,
			SqsMaxNumberOfMessages:      10,
			SqsMessageVisibilityTimeout: 20,
			PollDelayInMilliseconds:     100,
		})

	c.Start()
}

func handle(m *sqs.Message) error {
  fmt.Println("Message Body:", *(m.Body))
  //emulate processing time
  time.Sleep(time.Second * 2) 
  return nil
}
