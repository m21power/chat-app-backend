package infrastructure

import (
	iuserrepository "chat-app/domain/contracts/repository"
	entities "chat-app/domain/entities"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}
func NewUserRepository(db *gorm.DB) iuserrepository.IUserRepository {
	return &UserRepository{db: db}
}
func (r *UserRepository) CreateUser(user *entities.User) error {
	user.Role = "USER"
	return r.db.Create(user).Error
}
func (r *UserRepository) Login(user *entities.User) (*entities.User,error) {
	existingUser, err := r.GetUserByPhoneNumber(user.PhoneNumber)
	if err != nil {
		return nil, err
	}
	if existingUser == nil {
		return nil, fmt.Errorf("user not found")
	}
	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))
	if err != nil {
		return nil, fmt.Errorf("invalid password")
	}
	return existingUser, nil
}

func (r *UserRepository) DeleteUser(id uint) error {
	response := r.db.Delete(&entities.User{}, id)
	if response.Error != nil  {
		return response.Error
	}
	if response.RowsAffected == 0{
		return fmt.Errorf("user not found")
	}
	return nil
}

func (r *UserRepository) GetUserByID(id uint) (*entities.User, error) {
	user := new(entities.User)
	err :=r.db.Where("id = ?",id).First(&user).Error
	if err != nil{
		return nil,err
	}
	return user,nil
}

func (r *UserRepository) GetUserByPhoneNumber(phoneNumber string) (*entities.User, error) {
	user := new(entities.User)
	err :=r.db.Where("phone_number = ?",phoneNumber).First(&user).Error
	if err != nil{
		return nil,err
	}
	return user,nil
}

func (r *UserRepository) UpdateUser(user *entities.User) error {
	return r.db.Save(user).Error
}



