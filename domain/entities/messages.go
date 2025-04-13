package domain

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	SenderID       uint      `gorm:"not null" json:"sender_id"`
	ReceiverID     uint      `gorm:"not null" json:"receiver_id"`
	Content        string    `gorm:"type:text;not null" json:"content"`
	SentAt         time.Time `gorm:"autoCreateTime" json:"sent_at"`
	Seen           bool      `gorm:"default:false" json:"seen"`
	Deleted        bool      `gorm:"default:false" json:"deleted"`
	ConversationID uuid.UUID `gorm:"type:uuid;not null" json:"conversation_id"`
}
type Conversation struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	User1ID       uint      `gorm:"not null" json:"user1_id"`
	User2ID       uint      `gorm:"not null" json:"user2_id"`
	LastMessageID uuid.UUID `gorm:"type:uuid" json:"last_message_id"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Optional: Associations
	User1       User    `gorm:"foreignKey:User1ID"`
	User2       User    `gorm:"foreignKey:User2ID"`
	LastMessage Message `gorm:"foreignKey:LastMessageID"`
}

type SendMessageRequest struct {
	SenderID       uint      `json:"sender_id"`
	ReceiverID     uint      `json:"receiver_id"`
	Content        string    `json:"content"`
	ConversationID uuid.UUID `json:"conversation_id"`
}
