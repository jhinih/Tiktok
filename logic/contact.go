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

	"github.com/gin-gonic/gin"
)

type ContactLogic struct {
}

func NewContactLogic() *ContactLogic {
	return &ContactLogic{}
}

func (l *ContactLogic) AddFriend(ctx context.Context, req types.AddFriendRequest) (resp types.AddFriendResponse, err error) {
	if req.TargetName != "" {
		targetUser, err := request.NewUserRequest(global.DB).FindUserByName(req.TargetName)
		if uint(targetUser.ID) == req.UserId {
			//"不能加自己"
			return resp, response.ErrResponse(err, response.ME_AND_ME)
		}
		contact0 := model.Contact{}
		utils.DB.Where("owner_id =?  and target_id =? and type=1", req.UserId, targetUser.ID).Find(&contact0)
		if contact0.ID != 0 {
			//"不能重复添加"
			return resp, response.ErrResponse(err, response.FRIEND_YES_FRIEN)
		}
		err = request.NewContactRequest(global.DB).AddFriend(req.UserId, req.TargetName)
		if err != nil {
			return resp, response.ErrResponse(err, response.COMMON_FAIL)
		}
		//"添加好友成功"
		zlog.CtxDebugf(ctx, "添加好友成功: %v", req)
		return resp, nil

	}
	//"好友ID不能为空"
	return resp, response.ErrResponse(err, response.PARAM_IS_BLANK)
}
func (l *ContactLogic) SearchFriend(c *gin.Context, req types.SearchFriendRequest) (resp types.SearchFriendResponse, err error) {
	id := req.UserId
	users := request.NewContactRequest(global.DB).SearchFriend(uint(id))
	utils.RespOKList(c.Writer, users, len(users))
	resp.Ok = true
	return resp, nil
}
