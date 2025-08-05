.PHONY: run stop logs test migrate

run:
	docker compose up --build -d

stop:
	docker compose down -v

logs:
	docker compose logs -f

test:
	go test ./... -v

test-specific:
	go test ./... -v -run $(name)

sync-vendor:
	go mod tidy
	go mod vendor


AWS_ENDPOINT=http://localhost:4566
AWS_REGION=us-east-1

list-queues:
	aws --endpoint-url=$(AWS_ENDPOINT) sqs list-queues --region $(AWS_REGION)

receive-messages:
	aws --endpoint-url=http://localhost:4566 sqs receive-message \
        --queue-url http://sqs.us-east-1.localhost.localstack.cloud:4566/000000000000/todo_queue \
        --region us-east-1 \
        --max-number-of-messages 10 \
        --visibility-timeout 0 \
        --wait-time-seconds 1
