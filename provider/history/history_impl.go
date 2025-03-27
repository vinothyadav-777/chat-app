package history

import (
	"chat-app/entities"
	"chat-app/entities/repositories/chat_history"
	"context"
	"time"

	"chat-app/constants"
	config "chat-app/utils"
	log "github.com/sirupsen/logrus"
)

type HistoryProvider struct {
}

func GetHistoryProvider() Provider {
	return &HistoryProvider{}
}

func (provider *HistoryProvider) SaveHistory(ctx context.Context, ChatHistoryRepo chat_history.ChatHistoryRepo, entities []entities.ChatHistory) error {

	// Check if a deadline is set in the parent context
	deadline, ok := ctx.Deadline()
	var timeout time.Duration

	if ok {
		// Calculate the remaining time until the deadline
		timeout = time.Until(deadline)
	} else {
		// Fallback to a default timeout if no deadline is set
		timeout = time.Duration(config.GetClient().GetIntD(constants.DatabaseConfig, constants.HistorySaveRecordsDBCallTimeout, 2000)) * time.Millisecond
	}

	// Create a new context with the timeout
	dbCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	err := ChatHistoryRepo.SaveBulk(dbCtx, entities)
	if err != nil {
		log.WithField("error", err).WithField("historyEntities", entities).Warnln("error in save history")
		return err
	}
	return nil
}
