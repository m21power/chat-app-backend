package infrastructure

import (
	iuserrepository "chat-app/domain/contracts/repository"
	entities "chat-app/domain/entities"
	"fmt"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}
func NewUserRepository(db *gorm.DB) iuserrepository.IUserRepository {
	return &UserRepository{db: db}
}
func (r *UserRepository) CreateUser(user *entities.User) error {
	return r.db.Create(user).Error
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



