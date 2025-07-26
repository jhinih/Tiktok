package request

import (
	"Tiktok/model"
	"Tiktok/utils"
	"fmt"
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
func (r *CommunityRequest) CreateCommunity(community model.Community) (err error) {
	tx := utils.DB.Begin()
	//事务一旦开始，不论什么异常最终都会 Rollback
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := r.DB.Create(&community).Error; err != nil {
		fmt.Println(err)
		tx.Rollback()
		return err
	}
	contact := model.Contact{}
	contact.OwnerId = community.OwnerId
	contact.TargetId = uint(community.ID)
	contact.Type = 2 //群关系
	if err := r.DB.Create(&contact).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return err
}

// 查找某个群

func (r *CommunityRequest) FindCommunityByName(name string) model.Community {
	community := model.Community{}
	utils.DB.Where("name = ?", name).First(&community)
	return community
}

func (r *CommunityRequest) LoadCommunity(OwnerId int) (err error) {
	contacts := make([]model.Contact, 0)
	objIds := make([]uint64, 0)
	r.DB.Where("owner_id = ? and type=2", OwnerId).Find(&contacts)
	for _, v := range contacts {
		objIds = append(objIds, uint64(v.TargetId))
	}

	data := make([]*model.Community, 10)
	r.DB.Where("id in ?", objIds).Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}
	//utils.DB.Where()
	return nil
}

// 加入群聊
func (r *CommunityRequest) JoinGroup(OwnerId, ComId int) (err error) {
	contact := model.Contact{}
	contact.OwnerId = uint(OwnerId)
	//contact.TargetId = comId
	contact.Type = 2
	community := model.Community{}

	r.DB.Where("id=? or name=?", ComId, ComId).Find(&community)
	if community.Name == "" {
		//"没有找到群"
	}
	//utils.DB.Where("owner_id=? and target_id=? and type =2 ", userId, comId).Find(&contact)
	r.DB.Where("owner_id=? and target_id=? and type =2 ", OwnerId, community.ID).Find(&contact)
	if contact.TimeModel.CreatedTime != 0 {
		//"已加过此群"
	} else {
		contact.TargetId = uint(community.ID)
		utils.DB.Create(&contact)
		//"加群成功"
	}
	return nil
}
