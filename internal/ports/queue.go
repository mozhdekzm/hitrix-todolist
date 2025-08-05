package ports

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/mozhdekzm/heli-task/internal/domain"
)

type Queue interface {
	Publish(todo domain.TodoItem) error
}

type SQSClient interface {
	SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
}
