package applications

import (
	imessagerepository "chat-app/domain/contracts/repository"
	imessageusecase "chat-app/domain/contracts/usecases"
	entities "chat-app/domain/entities"

	"github.com/google/uuid"
)

type MessageUsecase struct {
	repository imessagerepository.IMessageRepository
}

func NewMessageUsecase(imrepo imessagerepository.IMessageRepository) imessageusecase.IMessageUsecase {
	return &MessageUsecase{repository: imrepo}
}

func (m *MessageUsecase) CreateMessage(message *entities.SendMessageRequest) (*entities.Message, error) {
	return m.repository.CreateMessage(message)
}
func (m *MessageUsecase) GetMessageByID(id uuid.UUID) (*entities.Message, error) {
	return m.repository.GetMessageByID(id)
}
func (m *MessageUsecase) GetMessagesByConversationID(conversationID uuid.UUID) ([]entities.Message, error) {
	return m.repository.GetMessagesByConversationID(conversationID)
}
func (m *MessageUsecase) GetMessagesBySenderAndReceiver(senderID, receiverID uint) ([]entities.Message, error) {
	return m.repository.GetMessagesBySenderAndReceiver(senderID, receiverID)
}
func (m *MessageUsecase) UpdateMessage(message *entities.Message) error {
	return m.repository.UpdateMessage(message)
}
func (m *MessageUsecase) DeleteMessage(id uuid.UUID) error {
	return m.repository.DeleteMessage(id)
}
