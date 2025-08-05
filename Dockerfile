FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o todo-app ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/todo-app .
COPY migrations ./migrations

EXPOSE 8080

CMD ["./todo-app"]
