package queue

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/vinothyadav-777/chat-app/constants"
	"github.com/vinothyadav-777/chat-app/models"
	"github.com/vinothyadav-777/chat-app/services/consumer"
	config "github.com/vinothyadav-777/chat-app/utils"
	log "github.com/sirupsen/logrus"
)

func BeginProcessing(queueService *consumer.QueueService, consumer SQSConsumer) {
	consumerType := os.Getenv(constants.ConsumerType)
	bufferLengthMap, err := config.GetClient().GetMap(constants.ApplicationConfig, constants.ConsumerBufferLength)
	if err != nil {
		log.Fatalln("Error while getting ConsumerBufferLength config: err", err)
	}
	bufferLength := int(bufferLengthMap[consumerType].(float64))
	buffer := make(chan models.Message, bufferLength)
	log.WithField("bufferLenght", bufferLength).WithField("consumerType", consumerType).Infoln("Buffer Length")
	if bufferLength != 0 {
		BeginProcessingBuffer(buffer, queueService, consumer, int64(bufferLength))
	}
	for {
		payloads, err := queueService.Receive()
		if err != nil {
			log.Errorln(constants.QueueReceiveError, err, payloads)
			return
		}
		if len(payloads) == 0 {
			log.Debugln(constants.NoMessageInQueue)
			continue
		}

		for _, payload := range payloads {
			go ProcessMessage(payload, queueService, consumer)
		}
	}
}

func BeginBulkProcessing(queueService *consumer.QueueService, consumer SQSConsumer) {

	consumerType := os.Getenv(constants.ConsumerType)
	bufferLengthMap, err := config.GetClient().GetMap(constants.ApplicationConfig, constants.ConsumerBufferLength)
	if err != nil {
		log.Fatalln("Error while getting ConsumerBufferLength config: err", err)
	}
	bufferLength := int(bufferLengthMap[consumerType].(float64))
	buffer := make(chan []models.Message, bufferLength)
	if bufferLength != 0 {
		BeginBulkProcessingBuffer(buffer, queueService, consumer, int64(bufferLength))
	}
	for {
		payloads, err := queueService.Receive()
		if err != nil {
			log.Errorln(constants.QueueReceiveError, err, payloads)
			return
		}
		if len(payloads) == 0 {
			log.Debugln(constants.NoMessageInQueue)
			continue
		}
		go ProcessBulkMessage(payloads, queueService, consumer)
	}
}

func BeginProcessingBuffer(buffer chan models.Message, queueService *consumer.QueueService, consumer SQSConsumer, bufferLength int64) {
	go ConsumeMessage(buffer, queueService, consumer, bufferLength)
	for {
		payloads, err := queueService.Receive()
		if err != nil {
			log.Errorln(constants.QueueReceiveError, err, payloads)
			return
		}
		if len(payloads) == 0 {
			log.Debugln(constants.NoMessageInQueue)
			continue
		}
		for _, payload := range payloads {
			buffer <- payload
		}
	}
}

func BeginBulkProcessingBuffer(buffer chan []models.Message, queueService *consumer.QueueService, consumer SQSConsumer, bufferLength int64) {
	go ConsumeBulkMessage(buffer, queueService, consumer, bufferLength)
	for {
		payloads, err := queueService.Receive()
		if err != nil {
			log.Errorln(constants.QueueReceiveError, err, payloads)
			return
		}
		if len(payloads) == 0 {
			log.Debugln(constants.NoMessageInQueue)
			continue
		}
		buffer <- payloads

	}
}

func ConsumeMessage(buffer chan models.Message, queueService *consumer.QueueService, consumerService SQSConsumer, bufferLength int64) {
	for {
		wg := &sync.WaitGroup{}
		var i int64
		for i = 0; i < bufferLength; i++ {
			payload := <-buffer
			wg.Add(1)
			go parallelProcessing(payload, queueService, consumerService, wg)
		}
		wg.Wait()
	}
}

func ConsumeBulkMessage(buffer chan []models.Message, queueService *consumer.QueueService, consumerService SQSConsumer, bufferLength int64) {
	for {
		wg := &sync.WaitGroup{}
		var i int64
		for i = 0; i < bufferLength; i++ {
			payload := <-buffer
			wg.Add(1)
			go parallelBulkProcessing(payload, queueService, consumerService, wg)
		}
		wg.Wait()
	}
}

func parallelProcessing(payloads models.Message, queueService *consumer.QueueService, consumer SQSConsumer, wg *sync.WaitGroup) {
	defer wg.Done()
	timeout := config.GetClient().GetIntD(constants.ApplicationConfig, constants.MessageProcessingTimeoutInMilli, 5000)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Millisecond)
	defer cancel()
	result := make(chan string)
	go func() {
		ProcessMessage(payloads, queueService, consumer)
		result <- "completed"
	}()
	select {
	case <-ctx.Done():
		err := queueService.Delete(payloads.ReceiptHandle)
		if err != nil {
			log.Errorln("error in delete sqs", err)
		}
		msg := fmt.Sprintf("task aborted %+v", payloads)
		log.Info("Go Routine processing timeout with payload", msg)
		return
	case <-result:
		return
	}
}

func parallelBulkProcessing(payloads []models.Message, queueService *consumer.QueueService, consumer SQSConsumer, wg *sync.WaitGroup) {
	defer wg.Done()
	ProcessBulkMessage(payloads, queueService, consumer)
}

func ProcessMessage(payload models.Message, queueService *consumer.QueueService, consumer SQSConsumer) {
	ctx := context.Background()
	deleteFromQueue, err := consumer.ProcessMessage(ctx, payload)
	if deleteFromQueue {
		err := queueService.Delete(payload.ReceiptHandle)
		if err != nil {
			log.Errorln(constants.QueueDeleteError, err, payload)
		}
	} else {
		log.WithField("payload", payload).WithError(err).Warningln("Not Deleting from Queue Retrying")
	}
}

func ProcessBulkMessage(payloads []models.Message, queueService *consumer.QueueService, consumer SQSConsumer) {
	ctx := context.Background()
	deleteFromQueue, _, err := consumer.ProcessBulkMessage(ctx, payloads)
	if deleteFromQueue {
		err := queueService.DeleteBatch(payloads)
		if err != nil {
			log.Errorln(constants.QueueDeleteError, err, payloads)
		}
	} else {
		log.WithField("payload", payloads).WithError(err).Warningln("Not Deleting from Queue Retrying")
	}
}
