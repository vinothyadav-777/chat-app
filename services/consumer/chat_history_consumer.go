package consumer

import (
	"context"
	"encoding/json"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/vinothyadav-777/chat-app/entities"
	"github.com/vinothyadav-777/chat-app/entities/repositories/chat_history"
	"github.com/vinothyadav-777/chat-app/models"
	"github.com/vinothyadav-777/chat-app/provider/history"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

type HistoryConsumer struct {
	provider        history.Provider
	chatHistoryRepo chat_history.ChatHistoryRepo
	redis           redis.Client
}

func GetHistoryConsumer(provider history.Provider, chatHistoryRepo chat_history.ChatHistoryRepo, redis redis.Client) *HistoryConsumer {
	return &HistoryConsumer{provider: provider, chatHistoryRepo: chatHistoryRepo, redis: redis}
}

func (consumer *HistoryConsumer) ProcessMessage(ctx context.Context, payload models.Message) (bool, error) {
	panic("ProcessMessage Not supported for HistoryConsumer")
}

func (consumer *HistoryConsumer) ProcessBulkMessage(ctx context.Context, payloads []models.Message) (bool, interface{}, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Errorln("recovered from exception", string(debug.Stack()))
		}
	}()
	chatHistory := getBulkHistoryEntity(payloads)
	log.WithField("path", "HistoryConsumer").WithField("historyEntities", chatHistory).Infoln("Request received")

	err := consumer.provider.SaveHistory(ctx, consumer.chatHistoryRepo, chatHistory)
	if err != nil {
		return false, nil, err
	}

	return true, nil, nil
}

func (consumer *HistoryConsumer) GetConsumerName() string {
	return "HistoryConsumer"
}

func getBulkHistoryEntity(payloads []models.Message) []entities.ChatHistory {
	historyEntities := make([]entities.ChatHistory, 0)
	for _, payload := range payloads {
		unix, err := strconv.ParseInt(payload.Attributes["sent_timestamp"], 10, 64)
		if err == nil {
			delay := time.Since(time.Unix(unix, 0)).Seconds()
			if delay > 0 {
				log.WithField("delaySec", delay).Infoln("SQS-Delay")
			}
		}

		request := entities.ChatHistory{}
		bindErr := json.Unmarshal([]byte(payload.Body), &request)
		if bindErr != nil {
			log.WithField("error", bindErr.Error()).WithField("payload.Body", payload.Body).Errorln("binding error")
			continue
		}

		historyEntities = append(historyEntities, request)
	}
	return historyEntities
}
