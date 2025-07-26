package request

import (
	"Tiktok/global"
	"Tiktok/model"
	"Tiktok/utils"
	"gorm.io/gorm"
)

type ContactRequest struct {
	DB *gorm.DB
}

func NewContactRequest(db *gorm.DB) *ContactRequest {
	return &ContactRequest{
		DB: db,
	}
}

// 添加好友   自己的ID  ， 好友的ID
func (r *ContactRequest) AddFriend(userId uint, targetName string) (err error) {
	targetUser, err := NewUserRequest(global.DB).FindUserByName(targetName)

	tx := utils.DB.Begin()
	//事务一旦开始，不论什么异常最终都会 Rollback
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	contact := model.Contact{}
	contact.OwnerId = userId
	contact.TargetId = uint(targetUser.ID)
	contact.Type = 1
	if err := r.DB.Create(&contact).Error; err != nil {
		tx.Rollback()
		return err
	}
	contact1 := model.Contact{}
	contact1.OwnerId = uint(targetUser.ID)
	contact1.TargetId = userId
	contact1.Type = 1
	if err := r.DB.Create(&contact1).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return err
}
func (r *ContactRequest) SearchFriend(userId uint) []model.User {
	contacts := make([]model.Contact, 0)
	objIds := make([]uint64, 0)
	r.DB.Where("owner_id = ? and type=1", userId).Find(&contacts)
	for _, v := range contacts {
		objIds = append(objIds, uint64(v.TargetId))
	}
	users := make([]model.User, 0)
	r.DB.Where("id in ?", objIds).Find(&users)
	return users
}
