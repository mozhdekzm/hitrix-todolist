FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build static binary for Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/hitrix-todo-app ./cmd/main.go

# ---- Runtime image ----
FROM alpine:3.19

# Install certificates and tzdata for TLS/time correctness
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /out/hitrix-todo-app ./hitrix-todo-app
COPY migrations ./migrations

EXPOSE 8080

# Run the app directly; rely on DB retry logic or compose healthchecks
CMD ["./hitrix-todo-app"]