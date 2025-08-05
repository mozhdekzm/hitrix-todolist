package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	sqsconf "github.com/aws/aws-sdk-go-v2/config"
	awssqs "github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/mozhdekzm/heli-task/config"
	"github.com/mozhdekzm/heli-task/internal/adapters/httpserver"
	"github.com/mozhdekzm/heli-task/internal/adapters/postgres"
	"github.com/mozhdekzm/heli-task/internal/adapters/sqs"
	"github.com/mozhdekzm/heli-task/internal/application"
	"log"
)

func main() {
	fmt.Println("hi this is a todo:")
	// 1. Load Config
	cfg := config.LoadConfig()
	mg := postgres.NewMigrator(cfg)
	mg.Up()

	// 2. Initialize DB and Repository
	db := postgres.NewGormDB(cfg)
	todoRepo := postgres.NewTodoRepo(db)

	// 3. Initialize Queue (SQS mocked via LocalStack)
	awsCfg, err := sqsconf.LoadDefaultConfig(context.TODO(),
		sqsconf.WithRegion("us-east-1"),
		sqsconf.WithEndpointResolver(
			aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           "http://localstack:4566",
					SigningRegion: "us-east-1",
				}, nil
			}),
		),
	)
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	sqsClient := awssqs.NewFromConfig(awsCfg)
	queueAdapter := sqs.NewSQSAdapter(sqsClient, "todo_queue")

	// 4. Initialize Service
	todoService := application.NewTodoService(todoRepo, queueAdapter)

	// 5. Initialize HTTP Server
	server := httpserver.New(cfg, *todoService)

	// 6. Start Server
	server.Serv()
}
