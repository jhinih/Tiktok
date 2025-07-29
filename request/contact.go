package request

import (
	"Tiktok/log/zlog"
	"Tiktok/model"
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

// 添加好友
func (r *ContactRequest) AddFriend(ownerUser, targetUser model.User) (err error) {
	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			return
		}
		if err != nil {
			tx.Rollback()
		}
	}()

	// 创建正向联系人记录
	contact := model.Contact{
		OwnerId:    ownerUser.ID,
		OwnerName:  ownerUser.Username,
		TargetId:   targetUser.ID,
		TargetName: targetUser.Username,
		Type:       1,
	}

	// 添加详细日志
	zlog.Debugf("Creating forward contact: %+v", contact)
	if err = tx.Create(&contact).Error; err != nil {
		zlog.Errorf("Failed to create forward contact: %v", err)
		return err
	}
	zlog.Debugf("Forward contact created successfully, ID: %d", contact.ID)

	// 创建反向联系人记录
	reverseContact := model.Contact{
		OwnerId:    targetUser.ID,
		OwnerName:  targetUser.Username,
		TargetId:   ownerUser.ID,
		TargetName: ownerUser.Username,
		Type:       1,
	}

	zlog.Debugf("Creating reverse contact: %+v", reverseContact)
	if err = tx.Create(&reverseContact).Error; err != nil {
		zlog.Errorf("Failed to create reverse contact: %v", err)
		return err
	}
	zlog.Debugf("Reverse contact created successfully, ID: %d", reverseContact.ID)

	if err = tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// 搜索所有好友
func (r *ContactRequest) SearchFriend(userId int64) ([]model.User, error) {
	contacts := make([]model.Contact, 0)
	objIds := make([]int64, 0)
	err := r.DB.Where("owner_id = ? and type=1", userId).Find(&contacts).Error
	if err != nil {
		return nil, err
	}
	for _, v := range contacts {
		objIds = append(objIds, v.TargetId)
	}
	users := make([]model.User, 0)
	err = r.DB.Where("id in ?", objIds).Find(&users).Error
	return users, err
}

// 是否是好友
func (r *ContactRequest) IsFriend(UserId, TargetID int64) (model.Contact, error) {
	var contact model.Contact
	err := r.DB.Where("owner_id =?  and target_id =? and type=1", UserId, TargetID).Find(&contact).Error
	return contact, err
}

// 创建群聊
func (r *ContactRequest) CreatCommunity(contact model.Contact) error {
	err := r.DB.Create(&contact).Error
	return err
}

//// 群是否被拥有
//func (r *ContactRequest) SearchUserByGroupId(communityId int64) []model.Contact {
//	contacts := make([]model.Contact, 0)
//	r.DB.Where("target_id = ? and type=2", communityId).Find(&contacts)
//	return contacts
//}
