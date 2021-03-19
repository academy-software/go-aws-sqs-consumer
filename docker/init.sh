export AWS_DEFAULT_REGION=us-east-1
export LOCALSTACK_HOST=http://localhost
export LOCALSTACK_PORT=4566

aws --endpoint-url=${LOCALSTACK_HOST}:${LOCALSTACK_PORT} sqs create-queue --queue-name my-queue