package history

import (
	"github.com/vinothyadav-777/chat-app/entities"
	"github.com/vinothyadav-777/chat-app/entities/repositories/chat_history"
	"context"
)

type Provider interface {
	SaveHistory(ctx context.Context, chatHistoryRepo chat_history.ChatHistoryRepo, entities []entities.ChatHistory) error
}
