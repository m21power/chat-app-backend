package transport

import (
	auth "chat-app/Auth"
	imessageusecase "chat-app/domain/contracts/usecases"
	entities "chat-app/domain/entities"
	"chat-app/util"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type MessageHandler struct {
	usecase imessageusecase.IMessageUsecase
}

func NewMessageHandler(usecase imessageusecase.IMessageUsecase) *MessageHandler {
	return &MessageHandler{
		usecase: usecase,
	}
}
func (h *MessageHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	message := new(entities.SendMessageRequest)
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}
	userID, err := strconv.Atoi(r.Context().Value(auth.ContextUserID).(string))
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if uint(userID) != message.SenderID {
		util.WriteError(w, http.StatusUnauthorized, fmt.Errorf("Unauthorized"))
		return
	}
	createdMessage, err := h.usecase.CreateMessage(message)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}
	util.WriteJSON(w, http.StatusOK, createdMessage)
}
func (h *MessageHandler) GetMessageByID(w http.ResponseWriter, r *http.Request) {
	// get message by id
}
func (h *MessageHandler) GetMessagesByConversationID(w http.ResponseWriter, r *http.Request) {
	value := mux.Vars(r)
	convId := value["convID"]
	uuidConvID, err := uuid.Parse(convId)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid conversation id"))
		return
	}
	messages, err := h.usecase.GetMessagesByConversationID(uuidConvID)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}
	util.WriteJSON(w, http.StatusOK, messages)

}
func (h *MessageHandler) GetMessagesBySenderAndReceiver(w http.ResponseWriter, r *http.Request) {
	type payload struct {
		SenderID   uint `json:"senderId"`
		ReceiverID uint `json:"receiverId"`
	}
	var requestPayload payload
	err := json.NewDecoder(r.Body).Decode(&requestPayload)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}
	userID := r.Context().Value(auth.ContextUserID).(string)
	userid, err := strconv.Atoi(userID)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if uint(userid) != requestPayload.SenderID && uint(userid) != requestPayload.ReceiverID {
		util.WriteError(w, http.StatusUnauthorized, fmt.Errorf("Unauthorized"))
		return
	}
	messages, err := h.usecase.GetMessagesBySenderAndReceiver(requestPayload.SenderID, requestPayload.ReceiverID)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}
	util.WriteJSON(w, http.StatusOK, messages)
}
func (h *MessageHandler) UpdateMessage(w http.ResponseWriter, r *http.Request) {
	value := mux.Vars(r)
	id := value["id"]
	messageID, err := uuid.Parse(id)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid message id"))
		return
	}
	message := new(entities.Message)
	err = json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}
	userID := r.Context().Value(auth.ContextUserID).(string)
	userid, err := strconv.Atoi(userID)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}
	oldMessage, err := h.usecase.GetMessageByID(messageID)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if uint(userid) != oldMessage.SenderID {
		util.WriteError(w, http.StatusUnauthorized, fmt.Errorf("Unauthorized"))
		return
	}
	message.ID = messageID
	err = h.usecase.UpdateMessage(message)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}
	oldMessage.Content = message.Content
	util.WriteJSON(w, http.StatusOK, oldMessage)
}
func (h *MessageHandler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	value := mux.Vars(r)
	id := value["id"]
	messageID, err := uuid.Parse(id)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid message id"))
		return
	}
	userID := r.Context().Value(auth.ContextUserID).(string)
	userid, err := strconv.Atoi(userID)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}
	oldMessage, err := h.usecase.GetMessageByID(messageID)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if uint(userid) != oldMessage.SenderID {
		util.WriteError(w, http.StatusUnauthorized, fmt.Errorf("Unauthorized"))
		return
	}
	err = h.usecase.DeleteMessage(messageID)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}
	util.WriteJSON(w, http.StatusOK, map[string]string{"message": "Message Deleted Successfully!"})
}
