package rabbitmq

import (
	"github.com/vinothyadav-777/chat-app/models"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

func GetRabbitMQClient(amqpURI, queueName string) (*RabbitMQ, error) {
	// Connect to RabbitMQ server
	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	// Create a channel
	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	// Declare a queue
	queue, err := channel.QueueDeclare(
		queueName, // queue name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	// Return the RabbitMQ struct with the connection, channel, and queue
	return &RabbitMQ{
		conn:    conn,
		channel: channel,
		queue:   queue,
	}, nil
}

func (r *RabbitMQ) Receive() ([]models.Message, error) {
	// Consume messages from the queue
	messages, err := r.channel.Consume(
		r.queue.Name, // queue name
		"",           // consumer tag
		true,         // auto-acknowledge
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register a consumer: %w", err)
	}

	// Initialize the result slice
	var result []models.Message

	// Read the messages from the channel
	for msg := range messages {
		// Map the message properties (headers) to attributes
		attrs := make(map[string]string)
		for key, value := range msg.Headers {
			// You can adjust the logic here depending on what data you want to store
			attrs[key] = fmt.Sprintf("%v", value) // Converting to string
		}

		// Append the message to the result slice
		result = append(result, models.Message{
			ID:            msg.MessageId,
			ReceiptHandle: strconv.FormatUint(uint64(msg.DeliveryTag), 10), // DeliveryTag is the equivalent of the ReceiptHandle
			Body:          string(msg.Body),
			Attributes:    attrs,
		})
	}

	// If no messages were received
	if len(result) == 0 {
		return nil, nil
	}

	return result, nil
}

func (r *RabbitMQ) DeleteMessage(deliveryTag uint64) error {
	// Acknowledge the message to remove it from the queue
	err := r.channel.Ack(deliveryTag, false)
	if err != nil {
		return fmt.Errorf("failed to acknowledge message: %w", err)
	}
	return nil
}

func (r *RabbitMQ) Send(payload string, delay int64) error {
	// Send the message to the queue
	err := r.channel.Publish(
		"",           // exchange
		r.queue.Name, // routing key (queue name)
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         []byte(payload),
			DeliveryMode: amqp.Persistent, // Make the message persistent
			Expiration:   fmt.Sprintf("%d", delay/int64(time.Millisecond)),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}
	return nil
}

func (r *RabbitMQ) SendBatch(payloads []string, delay int64) error {
	for _, payload := range payloads {
		err := r.Send(payload, delay)
		if err != nil {
			log.Printf("Failed to send message: %s", err)
		}
	}
	return nil
}

func (r *RabbitMQ) Delete(messageID string) error {
	// Convert the messageID (string) into uint64 for DeliveryTag
	deliveryTag, err := strconv.ParseUint(messageID, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid message ID: %w", err)
	}

	// Acknowledge the message to remove it from the queue
	err = r.channel.Ack(deliveryTag, false)
	if err != nil {
		return fmt.Errorf("failed to acknowledge message: %w", err)
	}
	return nil
}

func (r *RabbitMQ) Close() error {
	// Close the channel and connection
	if err := r.channel.Close(); err != nil {
		return fmt.Errorf("failed to close channel: %w", err)
	}
	if err := r.conn.Close(); err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}
	return nil
}
