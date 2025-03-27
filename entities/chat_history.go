package entities

import (
	"time"
)

type ChatHistory struct {
	MessageID   string     `gorm:"primaryKey;type:varchar(255);not null"`                // Primary Key (MessageID)
	Sender      string     `gorm:"type:varchar(255);not null;index:idx_sender_receiver"` // Composite index on Sender and Receiver
	Receiver    string     `gorm:"type:varchar(255);not null;index:idx_sender_receiver"` // Composite index on Sender and Receiver
	Message     string     `gorm:"type:text;not null"`
	FileURL     *string    `gorm:"type:varchar(255);default:null" json:"file_url,omitempty"`
	MsgType     string     `gorm:"type:varchar(50);not null"`
	Timestamp   time.Time  `gorm:"autoCreateTime;index"` // Index on Timestamp
	DeliveredAt *time.Time `gorm:"type:timestamp;default:null" json:"delivered_at,omitempty"`
	SeenAt      *time.Time `gorm:"type:timestamp;default:null" json:"seen_at,omitempty"`
	Status      string     `gorm:"type:varchar(50);default:'sent'"` //  sent --> delivered --> viewed
	EditedAt    *time.Time `gorm:"type:timestamp;default:null" json:"edited_at,omitempty"`
	DeletedAt   *time.Time `gorm:"type:timestamp;default:null" json:"deleted_at,omitempty"`
	ExpiresAt   *time.Time `gorm:"type:timestamp;default:null" json:"expires_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at,omitempty" gorm:"autoCreateTime"` // Automatically set the time when the record is created
	UpdatedAt   time.Time  `json:"updated_at,omitempty" gorm:"autoUpdateTime"` // Automatically set the time when the record is updated
}
