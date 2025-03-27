package main

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/vinothyadav-777/chat-app/config/flags"
	"github.com/vinothyadav-777/chat-app/constants"
	"github.com/vinothyadav-777/chat-app/externals/queues/rabbitmq"
	_ "github.com/vinothyadav-777/chat-app/externals/queues/sqs"
	"github.com/vinothyadav-777/chat-app/services/consumer"
	"github.com/vinothyadav-777/chat-app/services/consumer/queue"
	config "github.com/vinothyadav-777/chat-app/utils"

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

func loadLocalLocation() *time.Location {
	loc, err := time.LoadLocation(constants.TimeZone)
	if err != nil {
		log.Fatalln("Error in Load Local Location")
	}
	return loc
}