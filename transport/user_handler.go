package transport

import (
	auth "chat-app/Auth"
	domain "chat-app/domain/contracts/usecases"
	entities "chat-app/domain/entities"
	"chat-app/util"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type UserHandler struct{
	iuserusecase domain.IUserUsecase
}

func NewUserHandler(iuserusecase domain.IUserUsecase) UserHandler{
	return UserHandler{iuserusecase: iuserusecase}
}
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request){
	user := new(entities.User)
	err := json.NewDecoder(r.Body).Decode(&user)
	if err !=nil{
		util.WriteError(w,http.StatusBadRequest,err)
		return
	}
	existingUser,err := h.iuserusecase.Login(user)
	if err != nil{
		util.WriteError(w,http.StatusBadRequest,err)
		return
	}

	token, err :=auth.CreateToken(existingUser.ID,string(existingUser.Role))
	if err != nil{
		util.WriteError(w,http.StatusBadRequest,err)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	})
	util.WriteJSON(w,http.StatusOK,map[string]string{"token":token})
}
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request){
	user:= new(entities.User)
	err :=json.NewDecoder(r.Body).Decode(&user)
	if err !=nil{
		util.WriteError(w,http.StatusBadRequest,err)
		return
	}
	err =h.iuserusecase.CreateUser(user)
	if err != nil{
		util.WriteError(w, http.StatusBadRequest,err)
		return
	}
	util.WriteJSON(w,http.StatusOK,map[string]string{"message":"User Created Successfully!"})
}
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request){
	value :=mux.Vars(r)
	id := value["id"]
	Id, err := strconv.Atoi(id)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}
	user,err :=h.iuserusecase.GetUserByID(uint(Id))
	user.Password = ""
	if err != nil{
		util.WriteError(w,http.StatusBadRequest,err)
		return
	}
	util.WriteJSON(w,http.StatusOK,user)
}
func (h *UserHandler) GetUserByPhoneNumber(w http.ResponseWriter, r *http.Request){
	value :=mux.Vars(r)
	phoneNumber := value["phoneNumber"]
	user,err :=h.iuserusecase.GetUserByPhoneNumber(phoneNumber)
	if err != nil{
		util.WriteError(w,http.StatusBadRequest,err)
		return
	}
	user.Password = ""
	util.WriteJSON(w,http.StatusOK,user)
}
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request){
	value :=mux.Vars(r)
	id := value["id"]
	userID := r.Context().Value(auth.ContextUserID).(string)
	userRole := r.Context().Value(auth.ContextUserRole).(string)
	if userID != id && userRole != "ADMIN"{
		util.WriteError(w,http.StatusUnauthorized,fmt.Errorf("Unauthorized"))
		return
	}
	Id,err := strconv.Atoi(id)
	if err !=nil{
		util.WriteError(w,http.StatusBadRequest,err)
		return
	}
	oldUser,err :=h.iuserusecase.GetUserByID(uint(Id))
	if err !=nil{
		util.WriteError(w,http.StatusBadRequest,err)
		return
	}
	newUser:= new(entities.User)
	err =json.NewDecoder(r.Body).Decode(&newUser)
	if err !=nil{
		util.WriteError(w,http.StatusBadRequest,err)
		return
	}
	updatedUser := updateUser(oldUser,newUser)
	err =h.iuserusecase.UpdateUser(updatedUser)
	if err != nil{
		util.WriteError(w, http.StatusBadRequest,err)
		return
	}
	util.WriteJSON(w,http.StatusOK,map[string]string{"message":"User Updated Successfully!"})
}
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request){
	value := mux.Vars(r)
	id := value["id"]
	userID := r.Context().Value(auth.ContextUserID).(string)
	userRole := r.Context().Value(auth.ContextUserRole).(string)
	if userID != id && userRole != "ADMIN"{
		util.WriteError(w,http.StatusUnauthorized,fmt.Errorf("Unauthorized"))
		return
	}
	Id,err := strconv.Atoi(id)
	if err != nil{
		util.WriteError(w,http.StatusBadRequest,err)
		return
	}
	err =h.iuserusecase.DeleteUser(uint(Id))
	if err != nil{
		util.WriteError(w,http.StatusBadRequest,err)
		return
	}
	util.WriteJSON(w,http.StatusOK,map[string]string{"message":"User Deleted Successfully!"})

}


func updateUser(oldUser *entities.User,newUser *entities.User) *entities.User{
	if newUser.BackgroundPicture != "" {
		oldUser.BackgroundPicture = newUser.BackgroundPicture
	}
	if newUser.ProfilePicture != "" {
		oldUser.ProfilePicture = newUser.ProfilePicture
	}
	if newUser.Username != "" {
		oldUser.Username = newUser.Username
	}
	if newUser.PhoneNumber != "" {
		oldUser.PhoneNumber = newUser.PhoneNumber
	}
	if newUser.Password != "" {
		oldUser.Password = newUser.Password
	}
	if newUser.Bio != "" {
		oldUser.Bio = newUser.Bio
	}
	if newUser.Hobbies != ""{
		oldUser.Hobbies = newUser.Hobbies
	}
	return oldUser
}