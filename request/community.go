package request

import (
	"Tiktok/model"
	"gorm.io/gorm"
)

type CommunityRequest struct {
	DB *gorm.DB
}

func NewCommunityRequest(db *gorm.DB) *CommunityRequest {
	return &CommunityRequest{
		DB: db,
	}
}

// 创建群聊
func (r *CommunityRequest) CreateCommunity(community model.Community) (err error) {
	tx := r.DB.Begin()
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback()
		}
	}()

	if err = tx.Create(&community).Error; err != nil {
		return err
	}

	contact := model.Contact{
		OwnerId:    community.OwnerId,
		OwnerName:  community.OwnerName,
		TargetId:   community.ID,
		TargetName: community.Name,
		Type:       2, //群关系
		Desc:       community.Desc,
	}

	if err = tx.Create(&contact).Error; err != nil {
		return err
	}

	return tx.Commit().Error
}

// 查找某个群
func (r *CommunityRequest) FindCommunityByName(name string) model.Community {
	community := model.Community{}
	if err := r.DB.Where("name = ?", name).First(&community).Error; err != nil {
		return model.Community{}
	}
	return community
}

// 加载群列表
func (r *CommunityRequest) LoadUserCommunity(OwnerId int64) []model.Contact {
	contacts := make([]model.Contact, 0)
	r.DB.Where("owner_id = ? and type=2", OwnerId).Find(&contacts)
	return contacts
}
func (r *CommunityRequest) LoadCommunityUser(objIds []int64) []*model.Community {
	data := make([]*model.Community, 10)
	r.DB.Where("id in ?", objIds).Find(&data)
	return data
}

// 加入群聊
func (r *CommunityRequest) FindCommunityByNameOrId(ComId int64) model.Community {
	community := model.Community{}
	r.DB.Where("id=? or name=?", ComId, ComId).Find(&community)
	return community
}
func (r *CommunityRequest) IsInCommunity(OwnerId int64, community model.Community) model.Contact {
	contact := model.Contact{}
	contact.OwnerId = OwnerId
	contact.Type = 2
	r.DB.Where("owner_id=? and target_id=? and type =2 ", OwnerId, community.ID).Find(&contact)
	return contact
}
