test: docker-compose-down docker-compose-up test-execution docker-compose-down

docker-compose-up:
	cd ./docker && docker-compose up -d
	bash ./docker/wait-for-it localhost:4566

test-execution:
	sh -x ./docker/init.sh
	go test

docker-compose-down:
	cd ./docker && docker-compose down
