package infrastructure

import (
	imessagerepository "chat-app/domain/contracts/repository"
	entities "chat-app/domain/entities"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) imessagerepository.IMessageRepository {
	return &MessageRepository{db: db}
}
func (r *MessageRepository) CreateMessage(message *entities.SendMessageRequest) (*entities.Message, error) {
	var convId uuid.UUID
	if message.ConversationID != uuid.Nil {
		convId = message.ConversationID
	}
	messageEntity := &entities.Message{
		ID:             uuid.New(),
		SenderID:       message.SenderID,
		ReceiverID:     message.ReceiverID,
		ConversationID: convId,
		Content:        message.Content,
		Deleted:        false,
	}

	if err := r.db.Create(messageEntity).Error; err != nil {
		return nil, err
	}
	convId, err := r.UpdateConversation(message, messageEntity.ID)
	if err != nil {
		return nil, err
	}
	messageEntity.ConversationID = convId
	if err := r.db.Save(messageEntity).Error; err != nil {
		return nil, err
	}
	return messageEntity, nil
}
func (r *MessageRepository) GetMessageByID(id uuid.UUID) (*entities.Message, error) {
	var message entities.Message
	if err := r.db.First(&message, id).Error; err != nil {
		return nil, err
	}
	return &message, nil
}
func (r *MessageRepository) GetMessagesByConversationID(conversationID uuid.UUID) ([]entities.Message, error) {
	var messages []entities.Message
	if err := r.db.Where("conversation_id = ? AND deleted = ?", conversationID, false).Order("sent_at DESC").Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}
func (r *MessageRepository) GetMessagesBySenderAndReceiver(senderID, receiverID uint) ([]entities.Message, error) {
	var messages []entities.Message
	if err := r.db.Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?) AND deleted = ?", senderID, receiverID, receiverID, senderID, false).Order("sent_at DESC").Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}
func (r *MessageRepository) UpdateMessage(message *entities.Message) error {
	if err := r.db.Model(&entities.Message{}).Where("id = ?", message.ID).Update("content", message.Content).Error; err != nil {
		return err
	}
	return nil
}
func (r *MessageRepository) DeleteMessage(id uuid.UUID) error {
	var message entities.Message
	if err := r.db.First(&message, id).Error; err != nil {
		return err
	}
	message.Deleted = true
	if err := r.db.Save(&message).Error; err != nil {
		return err
	}
	var conversation entities.Conversation
	if err := r.db.Where("last_message_id = ?", id).First(&conversation).Error; err == nil {
		var lastMessages []entities.Message
		if err := r.db.Where("conversation_id = ? AND deleted = ?", conversation.ID, false).
			Order("sent_at DESC").Limit(2).Find(&lastMessages).Error; err != nil {
			return err
		}
		fmt.Println(lastMessages)
		if len(lastMessages) > 1 {
			conversation.LastMessageID = lastMessages[0].ID
		} else {
			conversation.LastMessageID = uuid.Nil
		}
		if err := r.db.Save(&conversation).Error; err != nil {
			return err
		}
	}
	return nil
}
func (r *MessageRepository) UpdateConversation(request *entities.SendMessageRequest, lastMessageId uuid.UUID) (uuid.UUID, error) {
	conv := entities.Conversation{
		ID:            request.ConversationID,
		User1ID:       request.SenderID,
		User2ID:       request.ReceiverID,
		LastMessageID: lastMessageId,
	}
	if request.ConversationID == uuid.Nil {
		var existingConv entities.Conversation
		if err := r.db.Where("(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)", request.SenderID, request.ReceiverID, request.ReceiverID, request.SenderID).First(&existingConv).Error; err == nil {
			existingConv.LastMessageID = lastMessageId
			if err := r.db.Save(&existingConv).Error; err != nil {
				return existingConv.ID, err
			}
			return existingConv.ID, nil
		}
		conv.ID = uuid.New()
		if err := r.db.Create(&conv).Error; err != nil {
			return uuid.Nil, err
		}
		return conv.ID, nil
	} else {
		if err := r.db.Save(&conv).Error; err != nil {
			return request.ConversationID, err
		}
		return request.ConversationID, nil

	}

}
