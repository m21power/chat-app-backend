package domain

import "time"
type Role string
const(
	Admin Role = "ADMIN"
	user Role = "USER"
)
type User struct{
	ID uint `gorm:"primaryKey;autoincrement" json:"id"`
	Name string `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Role Role `gorm:"column:role;type:varchar(255);not null" json:"role"`
	Username string `gorm:"column:username;type:varchar(255);unique" json:"username"`
	Password string `gorm:"column:password;type:varchar(255);not null" json:"password"`
	PhoneNumber string `gorm:"column:phone_number;type:varchar(255);unique;not null" json:"phone_number"`
	ProfilePicture string `gorm:"column:profile_picture;type:varchar(255)" json:"profile_picture"`
	BackgroundPicture string `gorm:"column:background_picture;type:varchar(255)" json:"background_picture"`
	Bio string `gorm:"column:bio;type:varchar(255)" json:"bio"`
	Hobbies string `gorm:"column:hobbies;type:varchar(255)" json:"hobbies"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}