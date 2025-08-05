package tests

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/stretchr/testify/assert"

	sqsadapter "github.com/mozhdekzm/heli-task/internal/adapters/sqs"
	"github.com/mozhdekzm/heli-task/internal/domain"
)

type mockSQSClient struct {
	sentBodies []string
	fail       bool
}

func (m *mockSQSClient) SendMessage(ctx context.Context, in *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
	if m.fail {
		return nil, errors.New("mock send failed")
	}
	m.sentBodies = append(m.sentBodies, *in.MessageBody)
	return &sqs.SendMessageOutput{}, nil
}

func TestSQSAdapter_Publish(t *testing.T) {
	tests := []struct {
		name        string
		todo        domain.TodoItem
		expectError bool
		mockFail    bool
	}{
		{
			name:        "valid publish",
			todo:        domain.TodoItem{Description: "Send to SQS", DueDate: time.Now()},
			expectError: false,
			mockFail:    false,
		},
		{
			name:        "client fails",
			todo:        domain.TodoItem{Description: "Failing SQS", DueDate: time.Now()},
			expectError: true,
			mockFail:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockSQSClient{fail: tt.mockFail}
			adapter := &sqsadapter.SQSAdapter{
				Client:   mockClient,
				QueueURL: "mock://queue",
			}

			err := adapter.Publish(tt.todo)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				var parsed domain.TodoItem
				json.Unmarshal([]byte(mockClient.sentBodies[0]), &parsed)
				assert.Equal(t, tt.todo.Description, parsed.Description)
			}
		})
	}
}
