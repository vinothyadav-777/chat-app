package queue

import (
	"chat-app/constants"
	"chat-app/provider/history"
	"chat-app/services/consumer"
	"context"
	"errors"
	"time"

	"chat-app/models"

	"chat-app/utils"

	log "github.com/sirupsen/logrus"
)

type MessageConsumer struct {
	provider      *history.Provider
	ConsumerQueue *consumer.QueueService
}

func GetMessageConsumer(provider *history.Provider, consumerQueue *consumer.QueueService, retryQueue *consumer.QueueService) *MessageConsumer {
	return &MessageConsumer{provider, consumerQueue}
}

func (consumer *MessageConsumer) ProcessMessage(ctx context.Context, payload models.Message) (bool, error) {

	event := models.MessageEvent{}
	bindErr := utils.JsonUnmarshal(payload.Body, &event)
	if bindErr != nil {
		log.WithField("error", bindErr.Error()).WithField("payload.Body", payload.Body).Errorln("binding error")
		return true, bindErr
	}
	ctx = context.WithValue(ctx, constants.UserId, event.MessageRequest.UserId)
	ctx = context.WithValue(ctx, constants.RequestTimestamp, time.Now())

	defer func() {
		if r := recover(); r != nil {
			//utils.LogError(ctx, "recovered from panic", *constants.ErrPanicRecovered.WithMessage(string(debug.Stack())))
		}
	}()
	log.WithField("path", "MessageConsumer").WithField("user_id", event.MessageRequest.UserId).WithField("type", event.MessageRequest.MessageType).WithField("payload.Body", event.MessageRequest).Infoln("Request received")
	event.QReceiveTime = time.Now().Unix()
	if event.QPublishTime > 0 {
		log.WithField("delaySec", time.Since(time.Unix(event.QPublishTime, 0)).Seconds()).Infoln("SQS-Delay")
	}
	request := event.MessageRequest
	if err := request.Validate(); err != nil {
		//utils.LogWarning(ctx, "ValidateError", err.Error())
		return true, err
	}

	return false, nil
}

func (consumer *MessageConsumer) ProcessBulkMessage(ctx context.Context, message []models.Message) (bool, interface{}, error) {
	return false, nil, errors.New("ProcessBulkMessage not supported for Message Consumer")

}

func (consumer *MessageConsumer) GetConsumerName() string {
	return "MessageConsumer"
}
