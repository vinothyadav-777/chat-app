package queues

import (
	"github.com/vinothyadav-777/chat-app/models"
)

type Queues interface {
	Receive() ([]models.Message, error)
	Delete(string) error
	//DeleteBatch(map[string]string) error
	Send(string, int64) error
	SendBatch([]string, int64) error
}
