package models

import (
	"errors"
	"regexp"
)

type MessageEvent struct {
	MessageRequest *MessageRequest `json:"request"`
	QReceiveTime   int64           `json:"queue_receive_time"`
	QPublishTime   int64           `json:"publish_time"`
	Delay          int64           `json:"delay"`
}

type MessageRequest struct {
	UserId         string `json:"user_id"`
	MessageContent string `json:"message_content"`
	MessageType    string `json:"message_type"`
	FileUrl        string `json:"file_url"`
	Status         string `json:"status"`
	ReceiverID     string `json:"receiver_id"`
}

// Validate checks the validity of MessageEvent and its nested MessageRequest
func (m *MessageEvent) Validate() error {
	if m.MessageRequest == nil {
		return errors.New("message request cannot be nil")
	}
	// Validate MessageRequest fields
	if err := m.MessageRequest.Validate(); err != nil {
		return err
	}

	// Validate MessageEvent fields
	if m.QReceiveTime <= 0 {
		return errors.New("queue_receive_time must be a positive integer")
	}

	if m.QPublishTime <= 0 {
		return errors.New("publish_time must be a positive integer")
	}

	if m.Delay < 0 {
		return errors.New("delay cannot be negative")
	}

	return nil
}

// Validate checks the validity of MessageRequest fields
func (m *MessageRequest) Validate() error {
	if m.UserId == "" {
		return errors.New("user_id cannot be empty")
	}

	if m.MessageContent == "" {
		return errors.New("message_content cannot be empty")
	}

	if m.MessageType == "" {
		return errors.New("message_type cannot be empty")
	}

	if m.FileUrl != "" {
		if !isValidURL(m.FileUrl) {
			return errors.New("file_url is not a valid URL")
		}
	}

	if m.Status == "" {
		return errors.New("status cannot be empty")
	}

	return nil
}

func isValidURL(url string) bool {
	re := regexp.MustCompile(`^(http|https)://[a-zA-Z0-9-_.]+(\.[a-zA-Z]{2,})+.*$`)
	return re.MatchString(url)
}
