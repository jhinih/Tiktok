package logic

import (
	"Tiktok/global"
	"Tiktok/log/zlog"
	"Tiktok/model"
	"Tiktok/request"
	"Tiktok/response"
	"Tiktok/types"
	"Tiktok/utils"
	"context"
	"strconv"
	"strings"
	"time"
)

type CommunityLogic struct {
}

func NewCommunityLogic() *CommunityLogic {
	return &CommunityLogic{}
}

// 新建群
func (l *CommunityLogic) CreateCommunity(ctx context.Context, req types.CreateCommunityRequest) (resp types.CreateCommunityResponse, err error) {
	clean := strings.Trim(req.OwnerId, `"`)
	Id, _ := strconv.ParseInt(clean, 10, 64)
	community := model.Community{
		OwnerId:   int64(Id),
		OwnerName: req.OwnerName,
		Name:      req.Name,
		Img:       req.Icon,
		Desc:      req.Desc,
		TimeModel: model.TimeModel{
			CreatedTime: time.Now().Unix(),
			UpdatedTime: time.Now().Unix(),
		},
	}
	if len(community.Name) == 0 {
		//"群名称不能为空"
		response.ErrResponse(err, response.COMMUNITY_IS_BLANK)
	}
	comIS := request.NewCommunityRequest(global.DB).FindCommunityByName(community.Name)
	if comIS.Name == "" {
		if utils.IsNumeric(community.Name) {
			//"开发者不允许你拿数字建群"
			response.ErrResponse(err, response.FACK_FACK_FACK)
		}
	} else if comIS.Name == community.Name {
		//"群聊已存在"
		response.ErrResponse(err, response.EMAIL_NOT_VALID)
	}
	err = request.NewCommunityRequest(global.DB).CreateCommunity(community)
	if err != nil {
		zlog.CtxErrorf(ctx, "创建群聊失败: %v", err)
		return resp, response.ErrResponse(err, response.DATABASE_ERROR)
	}
	return resp, nil
}

// 加载群列表
func (l *CommunityLogic) LoadCommunity(ctx context.Context, req types.LoadCommunityRequest) (resp types.LoadCommunityResponse, err error) {
	clean := strings.Trim(req.OwnerId, `"`)
	ownerId, _ := strconv.ParseInt(clean, 10, 64)
	//	name := c.Request.FormValue("name")
	contacts := make([]model.Contact, 0)
	objIds := make([]int64, 0)
	contacts = request.NewCommunityRequest(global.DB).LoadUserCommunity(int64(ownerId))
	for _, v := range contacts {
		objIds = append(objIds, v.TargetId)
	}
	data := request.NewCommunityRequest(global.DB).LoadCommunityUser([]int64(objIds))
	if len(data) != 0 {
		//response.ErrResponse(err, response.COMMUNITY_IS_BLANK)
	} else {
		//zlog.Infof()
	}
	return resp, err
}

// 加入群 userId uint, comId uint
func (l *CommunityLogic) JoinGroups(ctx context.Context, req types.JoinGroupsRequest) (resp types.JoinGroupsResponse, err error) {
	clean := strings.Trim(req.UserId, `"`)
	userId, _ := strconv.ParseInt(clean, 10, 64)
	clean = strings.Trim(req.ComId, `"`)
	comId, _ := strconv.ParseInt(clean, 10, 64)

	community := model.Community{}
	community = request.NewCommunityRequest(global.DB).FindCommunityByNameOrId(int64(comId))
	contact := model.Contact{}
	contact = request.NewCommunityRequest(global.DB).IsInCommunity(userId, community)
	if contact.TimeModel.CreatedTime != 0 {
		//"已加过此群"
	} else {
		contact.OwnerId = int64(userId)
		contact.TargetId = int64(community.ID)
		contact.TargetName = community.Name
		contact.TimeModel.CreatedTime = time.Now().Unix()
		contact.TimeModel.UpdatedTime = time.Now().Unix()
		contact.Type = 2
		contact.Desc = "" //后续添加描述
		err = request.NewContactRequest(global.DB).CreatCommunity(contact)
		//"加群成功"
	}
	return resp, err
}
