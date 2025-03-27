package chat_history

import (
	"chat-app/entities"
	"context"
)

type ChatHistoryRepo interface {
	Save(ctx context.Context, saveHistory entities.ChatHistory) error
	SaveBulk(ctx context.Context, bulkHistoryEntities []entities.ChatHistory) error
	DeleteExpiredData(ctx context.Context, expiredAT int64) error
}
