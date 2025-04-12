package domain

import domain "chat-app/domain/entities"

type IUserUsecase interface{
	CreateUser(user *domain.User) error
	GetUserByID(id uint) (*domain.User,error)
	GetUserByPhoneNumber(phoneNumber string) (*domain.User,error)
	UpdateUser(user *domain.User) error
	DeleteUser(id uint) error
}