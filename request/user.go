package request

import (
	"Tiktok/model"
	"fmt"
	"gorm.io/gorm"
)

type UserRequest struct {
	DB *gorm.DB
}

func NewUserRequest(db *gorm.DB) *UserRequest {
	return &UserRequest{
		DB: db,
	}
}

// GetUserProfileByID  获取用户信息
func (r *UserRequest) GetUserProfileByID(id int64) (model.User, error) {
	var user model.User
	err := r.DB.Where("id = ?", id).First(&user).Error
	return user, err
}

// UpdateUserProfile  更新用户信息
func (r *UserRequest) UpdateUserProfile(user model.User) error {
	err := r.DB.Save(&user).Error
	return err
}

// SetUserRole  设置用户角色
func (r *UserRequest) SetUserRole(id int64, role int) error {
	err := r.DB.Model(&model.User{}).Where("id = ?", id).Update("role", role).Error
	return err
}

func (r *UserRequest) GetUserList() []*model.User {
	data := make([]*model.User, 10)
	r.DB.Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}
	return data
}
func (r *UserRequest) FindUserByName(name string) (model.User, error) {
	user := model.User{}
	err := r.DB.Where("username = ?", name).First(&user).Error
	return user, err
}
func (r *UserRequest) FindUserById(id int64) (model.User, error) {
	user := model.User{}
	err := r.DB.Where("id = ?", id).First(&user).Error
	return user, err
}
