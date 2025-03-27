package queues

import (
	"chat-app/models"
)

type Queues interface {
	Receive() ([]models.Message, error)
	Delete(string) error
	//DeleteBatch(map[string]string) error
	Send(string, int64) (string, error)
	SendBatch([]string, int64) ([]string, error)
}
