package consumer

import (
	"context"
	"github.com/vinothyadav-777/chat-app/models"
)

type QueueConsumer interface {
	ProcessMessage(ctx context.Context, message models.Message) (bool, error)
	ProcessBulkMessage(ctx context.Context, message []models.Message) (bool, interface{}, error)
	GetConsumerName() string
}
