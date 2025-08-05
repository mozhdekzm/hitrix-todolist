package sqs

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/mozhdekzm/heli-task/internal/domain"
	"github.com/mozhdekzm/heli-task/internal/ports"
	"log"
)

type SQSAdapter struct {
	Client   ports.SQSClient
	QueueURL string
}

func NewSQSAdapter(client *sqs.Client, queueName string) *SQSAdapter {
	out, err := client.CreateQueue(context.TODO(), &sqs.CreateQueueInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		log.Fatal("failed to create queue:", err)
	}

	return &SQSAdapter{
		Client:   client,
		QueueURL: *out.QueueUrl,
	}
}

func (q *SQSAdapter) Publish(todo domain.TodoItem) error {
	data, _ := json.Marshal(todo)
	_, err := q.Client.SendMessage(context.TODO(), &sqs.SendMessageInput{
		QueueUrl:    &q.QueueURL,
		MessageBody: aws.String(string(data)),
	})
	return err
}
