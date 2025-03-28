package queue

import (
	"github.com/vinothyadav-777/chat-app/externals/queues"
	"github.com/vinothyadav-777/chat-app/models"
	"github.com/vinothyadav-777/chat-app/utils"
)

type QueueService struct {
	consumer queues.Queues
}

func GetQueueService(c queues.Queues) *QueueService {
	return &QueueService{c}
}

func (queueService *QueueService) Delete(messageId string) error {
	return queueService.consumer.Delete(messageId)
}

func (queueService *QueueService) DeleteBatch(payloads []models.Message) error {
	mapIdToHandle := make(map[string]string)
	for _, payload := range payloads {
		mapIdToHandle[payload.ID] = payload.ReceiptHandle
		return queueService.consumer.Delete(payload.ReceiptHandle)
	}
	return nil
}

func (queueService *QueueService) Send(payload string, delay int64) error {
	return queueService.consumer.Send(payload, delay)
}

func (queueService *QueueService) Receive() ([]models.Message, error) {
	return queueService.consumer.Receive()
}

func (queueService *QueueService) SendBatch(payloads []string, delay int64) error {
	//Splitting the queue payload because SQS doesn't allow more than 10 batch publish
	splittedPayloads := utils.SplitIntoSizedChunks(10, payloads)
	for _, sizedPayload := range splittedPayloads {
		err := queueService.consumer.SendBatch(sizedPayload, delay)
		if err != nil {
			return err
		}
	}
	return nil
}
