export AWS_ACCESS_KEY_ID=foo
export AWS_SECRET_ACCESS_KEY=bar
export AWS_DEFAULT_REGION=us-east-1
export LOCALSTACK_HOST=http://localhost
export LOCALSTACK_PORT=4566

aws --endpoint-url=${LOCALSTACK_HOST}:${LOCALSTACK_PORT} sqs create-queue --queue-name my-queue --cli-connect-timeout 0