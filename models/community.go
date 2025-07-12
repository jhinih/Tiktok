package models

import (
	"Tiktok/utils"
	"fmt"
	"gorm.io/gorm"
)

type Community struct {
	gorm.Model
	Name    string
	OwnerId uint
	Img     string
	Desc    string
}

func CreateCommunity(community Community) (int, string) {
	tx := utils.DB.Begin()
	//事务一旦开始，不论什么异常最终都会 Rollback
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	//if utils.DB.Where("name = ?", community.Name).First(&Community{})

	if len(community.Name) == 0 {
		return -1, "群名称不能为空"
	}
	if community.OwnerId == 0 {
		return -1, "请先登录"
	}
	comIS := FindCommunityByName(community.Name)
	if comIS.Name == "" {

		if utils.IsNumeric(community.Name) {
			return -1, "开发者不允许你拿数字建群"
		}
		if err := utils.DB.Create(&community).Error; err != nil {
			fmt.Println(err)
			tx.Rollback()
			return -1, "建群失败"
		}
		contact := Contact{}
		contact.OwnerId = community.OwnerId
		contact.TargetId = community.ID
		contact.Type = 2 //群关系
		if err := utils.DB.Create(&contact).Error; err != nil {
			tx.Rollback()
			return -1, "添加群关系失败"
		}
	} else {
		return -1, "群聊已存在"
	}

	tx.Commit()
	return 0, "建群成功"

}

// 加入群聊
func JoinGroup(userId uint, comId string) (int, string) {
	contact := Contact{}
	contact.OwnerId = userId
	//contact.TargetId = comId
	contact.Type = 2
	community := Community{}

	utils.DB.Where("id=? or name=?", comId, comId).Find(&community)
	if community.Name == "" {
		return -1, "没有找到群"
	}
	//utils.DB.Where("owner_id=? and target_id=? and type =2 ", userId, comId).Find(&contact)
	utils.DB.Where("owner_id=? and target_id=? and type =2 ", userId, community.ID).Find(&contact)
	if !contact.CreatedAt.IsZero() {
		return -1, "已加过此群"
	} else {
		contact.TargetId = community.ID
		utils.DB.Create(&contact)
		return 0, "加群成功"
	}
}

func LoadCommunity(ownerId uint) ([]*Community, string) {
	contacts := make([]Contact, 0)
	objIds := make([]uint64, 0)
	utils.DB.Where("owner_id = ? and type=2", ownerId).Find(&contacts)
	for _, v := range contacts {
		objIds = append(objIds, uint64(v.TargetId))
	}

	data := make([]*Community, 10)
	utils.DB.Where("id in ?", objIds).Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}
	//utils.DB.Where()
	return data, "查询成功"
}

// 查找某个群
func FindCommunityByName(name string) Community {
	community := Community{}
	utils.DB.Where("name = ?", name).First(&community)
	return community
}
