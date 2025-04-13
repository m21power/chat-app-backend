package applications

import (
	iuserrepository "chat-app/domain/contracts/repository"
	iuserusecase "chat-app/domain/contracts/usecases"
	entities "chat-app/domain/entities"

	"golang.org/x/crypto/bcrypt"
)


type UserUsecase struct{
	 repository iuserrepository.IUserRepository
}
func NewUserUsecase(repository iuserrepository.IUserRepository) iuserusecase.IUserUsecase{
	return &UserUsecase{repository: repository}
}

func (u *UserUsecase) Login(user *entities.User) (*entities.User,error){
	return u.repository.Login(user)
}
func (u *UserUsecase) CreateUser(user *entities.User) error{
	hashedPassword,err := hashPassword(user.Password)
	if err != nil{
		return err
	}
	user.Password = hashedPassword
	return u.repository.CreateUser(user)
}
func (u *UserUsecase) GetUserByID(id uint) (*entities.User,error){
	return u.repository.GetUserByID(id)
}
func (u *UserUsecase) GetUserByPhoneNumber(phoneNumber string) (*entities.User,error){
	return u.repository.GetUserByPhoneNumber(phoneNumber)
}
func (u *UserUsecase) UpdateUser(user *entities.User) error{
	hashedPassword,err := hashPassword(user.Password)
	if err != nil{
		return err
	}
	user.Password = hashedPassword
	return u.repository.UpdateUser(user)
}
func (u *UserUsecase)DeleteUser(id uint) error{
	return u.repository.DeleteUser(id)
}
func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
