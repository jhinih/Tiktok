package request

import (
	"Tiktok/model"
	"gorm.io/gorm"
)

type ChatRequest struct {
	DB *gorm.DB
}

func NewChatRequest(db *gorm.DB) *ChatRequest {
	return &ChatRequest{
		DB: db,
	}
}

func (r *ChatRequest) SearchUserByGroupId(communityId int64) ([]model.Contact, error) {
	contacts := make([]model.Contact, 0)
	//objIds := make([]uint, 0)
	err := r.DB.Where("target_id = ? and type=2", communityId).Find(&contacts).Error
	//for _, v := range contacts {
	//	objIds = append(objIds, uint(v.OwnerId))
	//}
	return contacts, err
}
func (r *ChatRequest) SearchUsers(users *[]model.User) error {
	if err := r.DB.Where("is_logout = ?", false).Find(users).Error; err != nil {
		return err
	}
	if len(*users) == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
