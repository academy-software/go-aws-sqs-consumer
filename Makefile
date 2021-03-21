test: docker-compose-down docker-compose-up test-execution docker-compose-down

docker-compose-up:
	cd ./docker && docker-compose up -d
	bash ./docker/wait-for-it localhost:4566
	sh -x ./docker/init.sh

test-execution:
	go test

docker-compose-down:
	cd ./docker && docker-compose down
