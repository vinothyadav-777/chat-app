package queue

import (
	"chat-app/models"
	"context"
)

type SQSConsumer interface {
	ProcessMessage(ctx context.Context, message models.Message) (bool, error)
	ProcessBulkMessage(ctx context.Context, message []models.Message) (bool, interface{}, error)
	GetConsumerName() string
}
