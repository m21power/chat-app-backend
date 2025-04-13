package domain

import (
	domain "chat-app/domain/entities"

	"github.com/google/uuid"
)

type IMessageUsecase interface {
	CreateMessage(message *domain.SendMessageRequest) (*domain.Message, error)
	GetMessageByID(id uuid.UUID) (*domain.Message, error)
	GetMessagesByConversationID(conversationID uuid.UUID) ([]domain.Message, error)
	GetMessagesBySenderAndReceiver(senderID, receiverID uint) ([]domain.Message, error)
	UpdateMessage(message *domain.Message) error
	DeleteMessage(id uuid.UUID) error
}
