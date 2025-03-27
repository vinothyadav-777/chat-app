package sqs

import (
	"chat-app/models"
	"chat-app/utils"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	log "github.com/sirupsen/logrus"
)

type Sqs struct {
	client              *sqs.SQS
	url                 string
	MaxNumberOfMessages int64
	WaitTimeSeconds     int64
	VisibilityTimeout   int64
}

func GetSQSConsumer(client *sqs.SQS, queueURL string, MaxNumberOfMessages, WaitTimeSeconds, VisibilityTimeout int64) *Sqs {
	queue := Sqs{client, queueURL, MaxNumberOfMessages, WaitTimeSeconds, VisibilityTimeout}
	return &queue
}

func (s *Sqs) Receive() ([]models.Message, error) {
	res, err := s.client.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:              aws.String(s.url),
		MaxNumberOfMessages:   aws.Int64(s.MaxNumberOfMessages),
		WaitTimeSeconds:       aws.Int64(s.WaitTimeSeconds),
		MessageAttributeNames: aws.StringSlice([]string{"All"}),
		VisibilityTimeout:     aws.Int64(s.VisibilityTimeout),
	})
	if err != nil {
		return nil, fmt.Errorf("receive: %w", err)
	}

	if len(res.Messages) == 0 {
		return nil, nil
	}
	result := make([]models.Message, 0, 5)
	for _, message := range res.Messages {
		attrs := make(map[string]string)
		for key, attr := range message.MessageAttributes {
			attrs[key] = *attr.StringValue
		}

		result = append(result, models.Message{
			ID:            *message.MessageId,
			ReceiptHandle: *message.ReceiptHandle,
			Body:          *message.Body,
			Attributes:    attrs,
		})
	}
	return result, nil
}

func (s *Sqs) Delete(rcvHandle string) error {

	if _, err := s.client.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(s.url),
		ReceiptHandle: aws.String(rcvHandle),
	}); err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}

func (s *Sqs) DeleteBatch(rcvHandles map[string]string) error {
	deleteMsgInputBatch := make([]*sqs.DeleteMessageBatchRequestEntry, 0)
	for id, rcvHandle := range rcvHandles {
		deleteMsgInput := sqs.DeleteMessageBatchRequestEntry{
			Id:            aws.String(id),
			ReceiptHandle: aws.String(rcvHandle),
		}
		deleteMsgInputBatch = append(deleteMsgInputBatch, &deleteMsgInput)
	}
	req, err := s.client.DeleteMessageBatchRequest(&sqs.DeleteMessageBatchInput{QueueUrl: aws.String(s.url),
		Entries: deleteMsgInputBatch,
	})
	if err.Failed != nil {
		log.Errorln("Error creating Delete Batch Request ", err, rcvHandles)
		return fmt.Errorf("delete: %+v", err)
	}
	sendErrr := req.Send()
	if sendErrr != nil {
		log.Errorln("Error Deleting Batch Request ", err, rcvHandles)
		return fmt.Errorf("delete: %+v", err)
	}
	return nil
}

func (s *Sqs) Send(payload string, delay int64) (string, error) {
	sentTimeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	DataType := "String"
	messageAttr := sqs.MessageAttributeValue{StringValue: &sentTimeStamp, DataType: &DataType}
	messageAttrMap := map[string]*sqs.MessageAttributeValue{"sent_timestamp": &messageAttr}
	res, err := s.client.SendMessage(&sqs.SendMessageInput{

		DelaySeconds:      aws.Int64(delay),
		MessageBody:       aws.String(payload),
		QueueUrl:          aws.String(s.url),
		MessageAttributes: messageAttrMap,
	})
	if res == nil || res.MessageId == nil {
		return "", err
	}
	return *res.MessageId, err
}

func (s *Sqs) SendBatch(payloads []string, delay int64) ([]string, error) {
	var inputs []*sqs.SendMessageBatchRequestEntry
	for _, payload := range payloads {
		id := utils.GetNuid()
		sentTimeStamp := strconv.FormatInt(time.Now().Unix(), 10)
		str := "String"
		messageAttr := sqs.MessageAttributeValue{StringValue: &sentTimeStamp, DataType: &str}
		messageAttrMap := map[string]*sqs.MessageAttributeValue{"sent_timestamp": &messageAttr}
		res := sqs.SendMessageBatchRequestEntry{
			Id:                &id,
			MessageBody:       aws.String(payload),
			DelaySeconds:      aws.Int64(delay),
			MessageAttributes: messageAttrMap,
		}
		inputs = append(inputs, &res)

	}

	req := sqs.SendMessageBatchInput{
		Entries:  inputs,
		QueueUrl: aws.String(s.url),
	}

	res, err := s.client.SendMessageBatch(&req)

	if err != nil {
		return nil, err
	}
	var ackIds []string
	for _, obj := range res.Successful {
		ackIds = append(ackIds, *obj.MessageId)
	}
	for _, obj := range res.Failed {
		if obj.Message != nil {
			if *obj.Message != "" {
				return ackIds, errors.New(*obj.Message)
			}
		}
	}
	return ackIds, nil

}
