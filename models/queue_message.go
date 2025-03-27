package models

type Message struct {
	ID            string
	ReceiptHandle string
	Body          string
	Attributes    map[string]string
}
