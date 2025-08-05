.PHONY: run stop logs test test-specific sync-vendor list-queues receive-messages

GO_IMAGE=golang:1.23.4
WORKDIR=/app

run:
	docker compose up --build -d

stop:
	docker compose down -v

logs:
	docker compose logs -f

test:
	docker run --rm -v $(PWD):/app -w /app \
		-e DB_HOST=todotask_postgres \
		--network=host \
		$(GO_IMAGE) go test ./... -v

test-specific:
	docker run --rm -v $(PWD):/app -w /app \
		-e DB_HOST=todotask_postgres \
		--network=host \
		$(GO_IMAGE) go test ./... -v -run $(name)

sync-vendor:
	docker run --rm -v $(PWD):/app -w /app $(GO_IMAGE) sh -c "go mod tidy && go mod vendor"

AWS_ENDPOINT=http://localhost:4566
AWS_REGION=us-east-1
QUEUE_URL=$(AWS_ENDPOINT)/000000000000/todo-queue

list-queues:
	aws --endpoint-url=$(AWS_ENDPOINT) sqs list-queues --region $(AWS_REGION)

receive-messages:
	aws --endpoint-url=$(AWS_ENDPOINT) sqs receive-message \
		--queue-url $(QUEUE_URL) \
		--region $(AWS_REGION) \
		--max-number-of-messages 10 \
		--visibility-timeout 0 \
		--wait-time-seconds 1
