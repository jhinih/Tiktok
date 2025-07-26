package logic

import (
	"Tiktok/global"
	"Tiktok/log/zlog"
	"Tiktok/model"
	"Tiktok/models"
	"Tiktok/request"
	"Tiktok/response"
	"Tiktok/types"
	"Tiktok/utils"
	"context"
	"github.com/gin-gonic/gin"
	"strconv"
)

type CommunityLogic struct {
}

func NewCommunityLogic() *CommunityLogic {
	return &CommunityLogic{}
}

// 新建群
func (l *CommunityLogic) CreateCommunity(ctx context.Context, req types.CreateCommunityRequest) (resp types.CreateCommunityResponse, err error) {
	ownerId := req.OwnerId
	name := req.Name
	icon := req.Icon
	desc := req.Desc

	community := model.Community{}
	community.OwnerId = uint(ownerId)
	community.Name = name
	community.Img = icon
	community.Desc = desc
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
	} else {
		//"群聊已存在"
		response.ErrResponse(err, response.EMAIL_NOT_VALID)
	}
	err = request.NewCommunityRequest(global.DB).CreateCommunity(community)
	if err != nil {
		zlog.CtxErrorf(ctx, "创建群聊失败: %v", err)
		return resp, response.ErrResponse(err, response.DATABASE_ERROR)
	}
	resp.Ok = true
	return resp, nil
}

// 加载群列表
func (l *CommunityLogic) LoadCommunity(ctx context.Context, req types.LoadCommunityRequest) (resp types.LoadCommunityResponse, err error) {
	ownerId := req.OwnerId
	//	name := c.Request.FormValue("name")
	err := request.NewCommunityRequest(global.DB).LoadCommunity(ownerId)
	return resp, err
}

// 加入群 userId uint, comId uint
func (l *CommunityLogic) JoinGroups(ctx context.Context, req types.JoinGroupsRequest) (resp types.JoinGroupsResponse, err error) {
	userId := req.UserId
	comId := req.ComId
	//	name := c.Request.FormValue("name")
	err := request.NewCommunityRequest(global.DB).JoinGroup(userId, comId)
	return resp, err
}
