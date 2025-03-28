package main

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/vinothyadav-777/chat-app/constants"
	"github.com/vinothyadav-777/chat-app/externals/queues/rabbitmq"
	_ "github.com/vinothyadav-777/chat-app/externals/queues/sqs"
	"github.com/vinothyadav-777/chat-app/services/consumer"
	"github.com/vinothyadav-777/chat-app/services/queue"
	"github.com/vinothyadav-777/chat-app/websocket"

	log "github.com/sirupsen/logrus"
)

func main() {
	_ = context.Background()

	//logger.InitializeLogger(constants.DefaultLogLevel)
	consumerType := os.Getenv(constants.ConsumerType)
	constants.Consumer = consumerType
	constants.LocalLocation = loadLocalLocation()
	log.Infoln("Initializing consumer for !!!!!!!!!!!!!!!!!!!!!!!!!", consumerType)
	switch consumerType {
	case "message_consumer":
		initQueueConsumer(consumerType)
	case "history":
		panic("history queue not yet implemented")
	case "retry":
		panic("retry queue not yet implemented")
	default:
		log.Fatal("In correct Consumer type ", consumerType)
	}
}

func initQueueConsumer(consumerType string) {
	queueService, consumerService := initQueueConsumers(consumerType)
	consumer.BeginProcessing(queueService, consumerService)
}

func initQueueConsumers(consumerType string) (*queue.QueueService, consumer.QueueConsumer) {

	queueURL := os.Getenv("SQS_" + strings.ToUpper(consumerType))
	retryQueueUrl := os.Getenv("SQS_RETRY")

	if consumerType == constants.Empty || queueURL == constants.Empty || retryQueueUrl == constants.Empty {
		log.Fatal("Error in initializing queue ", retryQueueUrl, queueURL, consumerType)
	}

	switch consumerType {
	case "message_consumer":
		queueProvider, err := rabbitmq.GetRabbitMQClient(queueURL, consumerType)
		if err != nil {
			log.Fatal("Error in initializing queue ", queueURL, consumerType)
		}

		consumerService := queue.GetQueueService(queueProvider)
		messageConsumer := consumer.GetMessageConsumer(nil, consumerService, consumerService)
		return consumerService, messageConsumer

	default:
		{
		}
	}
	return nil, nil
}

func loadLocalLocation() *time.Location {
	loc, err := time.LoadLocation(constants.TimeZone)
	if err != nil {
		log.Fatalln("Error in Load Local Location")
	}
	return loc
}

func initQueue() {
	queueName := ""
	queueUrl := os.Getenv(constants.MessageQueue)
	queueProvider, _ := rabbitmq.GetRabbitMQClient(queueUrl, queueName) //queueProvider
	websocket.QueueService = queue.GetQueueService(queueProvider)
}
