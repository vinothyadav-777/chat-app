package history

import (
	"context"
	"time"

	"github.com/vinothyadav-777/chat-app/entities"
	"github.com/vinothyadav-777/chat-app/entities/repositories/chat_history"

	log "github.com/sirupsen/logrus"
	"github.com/vinothyadav-777/chat-app/constants"
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
		timeout = time.Duration(constants.HistorySaveRecordsDBCallTimeoutInMilli) * time.Millisecond
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
