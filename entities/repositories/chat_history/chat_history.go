package chat_history

import (
	"github.com/vinothyadav-777/chat-app/entities"
	"github.com/vinothyadav-777/chat-app/entities/repositories"
	"context"
)

type ChatHistoryRepoImpl struct {
}

func GetChatHistoryRepoImpl() ChatHistoryRepo {
	return &ChatHistoryRepoImpl{}
}

func (chr *ChatHistoryRepoImpl) Save(ctx context.Context, saveHistory entities.ChatHistory) error {

	if err := repositories.DB().PutItem(ctx, &saveHistory); err != nil {
		return err
	}

	return nil
}

func (chr *ChatHistoryRepoImpl) SaveBulk(ctx context.Context, bulkHistoryEntities []entities.ChatHistory) error {
	return nil
}

func (chr *ChatHistoryRepoImpl) DeleteExpiredData(ctx context.Context, expiredAT int64) error {
	return nil
}
