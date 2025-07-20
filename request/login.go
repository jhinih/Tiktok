package request

import (
	"Tiktok/model"
	"gorm.io/gorm"
)

type LoginRequest struct {
	DB *gorm.DB
}

func NewLoginRequest(db *gorm.DB) *LoginRequest {
	return &LoginRequest{
		DB: db,
	}
}

// AddUser 新增用户
func (r *LoginRequest) AddUser(user model.User) error {
	return r.DB.Create(&user).Error
}

// GetUserByEmail 根据邮箱获取用户信息
func (r *LoginRequest) GetUserByEmail(email string) (model.User, error) {
	var user model.User
	err := r.DB.Where("email = ?", email).First(&user).Error
	return user, err
}
