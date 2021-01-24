# Go SQS Consumer

A package to easily create an AWS SQS message consumer

## Features

- Configurable number of goroutines that poll messages from SQS
- Backpressure
  - You can dinamically control the poll delay period by using the SetPollDelay function

## Usage

```golang
func main() {
	q := "http://localhost:4566/000000000000/my-queue"

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	c := consumer.New(q, handle,
		&consumer.Config{
			AwsSession:                  sess,
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
```

## Testing locally

A local instance of SQS can be created using [LocalStack](https://github.com/localstack/localstack) and docker-compose:

```yaml
version: '3'
services:
  localstack:
    image: localstack/localstack
    environment:
      - SERVICES=sqs
      - DOCKER_HOST=unix:///var/run/docker.sock
    ports:
      - "4566:4566" # sqs
```

Run it with:
```docker-compose up```

### Creating a queue
With the aws-cli, run the following command to create a new queue:

```aws sqs --endpoint=http://localhost:4566 create-queue --queue-name my-queue```

And the output will look like this:
```
{
    "QueueUrl": "http://localhost:4566/000000000000/my-queue"
}
```

### Sending messages
After creating a queue, messages can be sent to it by running:

```aws sqs send-message --endpoint=http://localhost:4566 --queue-url http://localhost:4566/000000000000/my-queue --message-body 'Message body'```

And the output will look like this:

```
{
    "MessageId": "8b563f4c-7e38-3cae-b6e5-47dfa7f2358e",
    "MD5OfMessageBody": "78b28efdf19c153fe37474bcd69abfbd"
}
```

### Creating an AWS session to LocalStack
To connect an application to the LocalStack queue, a session should be configured properly. Here is an example:

```
sess := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials("not", "empty", ""),
		DisableSSL:  aws.Bool(true),
		Region:      aws.String(endpoints.UsEast1RegionID),
		Endpoint:    aws.String("http://localhost:4566"),
	}))
```

