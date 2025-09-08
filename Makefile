.PHONY: run stop logs test test-specific sync-vendor benchmark

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
        		-e DB_HOST=mysql \
        		-e DB_PORT=3306 \
        		-e DB_USER=root \
        		-e DB_PASSWORD=password \
        		-e DB_NAME=todos \
        		--network=host \
        		$(GO_IMAGE) go test ./... -v

test-specific:
	docker run --rm -v $(PWD):/app -w /app \
    		-e DB_HOST=mysql \
    		-e DB_PORT=3306 \
    		-e DB_USER=root \
    		-e DB_PASSWORD=password \
    		-e DB_NAME=todos \
    		--network=host \
    		$(GO_IMAGE) go test ./... -v -run $(name)


benchmark:
	docker run --rm -v /Users/mozhde/Projects/gocast/ice-todo:/app -w /app \
      -e DB_HOST=host.docker.internal \
      -e DB_PORT=3306 \
      -e DB_USER=root \
      -e DB_PASSWORD=password \
      -e DB_NAME=todos \
      --network=host \
      $(GO_IMAGE) go test -bench=. -benchmem ./internal/usecase

sync-vendor:
	docker run --rm -v $(PWD):/app -w /app $(GO_IMAGE) sh -c "go mod tidy && go mod vendor"


