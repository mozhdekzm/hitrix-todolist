.PHONY: run stop logs test test-specific sync-vendor benchmark

GO_IMAGE=golang:1.23.4
WORKDIR=/app

run:
	docker compose up --build -d

stop:
	docker compose down -v

logs:
	docker compose logs -f

sync-vendor:
	docker run --rm -v $(PWD):/app -w /app $(GO_IMAGE) sh -c "go mod tidy && go mod vendor"




