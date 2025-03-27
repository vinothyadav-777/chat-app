package main

import (
	"chat-app/constants"
	"chat-app/queues/rabbitmq"
	_ "chat-app/queues/sqs"
	"chat-app/services/consumer"
	"chat-app/services/consumer/queue"
	config "chat-app/utils"
	"context"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()

	initConfigs(ctx)
	//logger.InitializeLogger(constants.DefaultLogLevel)
	consumerType := os.Getenv(constants.ConsumerType)
	constants.Consumer = consumerType
	constants.LocalLocation = loadLocalLocation()
	log.Infoln("Initializing consumer for !!!!!!!!!!!!!!!!!!!!!!!!!", consumerType)
	switch consumerType {
	case "message_consumer":
		initSqsConsumer(consumerType)
	case "history":
	case "retry":
		initSqsConsumer(consumerType)
	default:
		log.Fatal("In correct Consumer type ", consumerType)
	}

}

func initConfigs(ctx context.Context) {
	err := config.InitConfigs(flags.BaseConfigPath(), constants.ApplicationConfig, constants.DatabaseConfig)
	if err != nil {
		log.Fatal("error loading configs", err)
	}
}

func initSqsConsumer(consumerType string) {
	queueService, consumerService := initSQSConsumer(consumerType)
	queue.BeginProcessing(queueService, consumerService)
}

func initSQSConsumer(consumerType string) (*consumer.QueueService, queue.SQSConsumer) {

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

		consumerService := consumer.GetQueueService(queueProvider)
		messageConsumer := queue.GetMessageConsumer(nil, consumerService, consumerService)
		return consumerService, messageConsumer

	default:
		{
		}
	}
	return nil, nil
}
