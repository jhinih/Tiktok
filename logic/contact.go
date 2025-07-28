package logic

import (
	"Tiktok/global"
	"Tiktok/log/zlog"
	"Tiktok/model"
	"Tiktok/request"
	"Tiktok/response"
	"Tiktok/types"
	"context"
)

type ContactLogic struct {
}

func NewContactLogic() *ContactLogic {
	return &ContactLogic{}
}

func (l *ContactLogic) AddFriend(ctx context.Context, req types.AddFriendRequest) (resp types.AddFriendResponse, err error) {
	if req.TargetName == "" {
		zlog.CtxErrorf(ctx, "好友用户名不能为空")
		return resp, response.ErrResponse(err, response.PARAM_IS_BLANK)
	}

	// 查询发起用户
	ownerUser, err := request.NewUserRequest(global.DB).FindUserById(req.UserId)
	if err != nil {
		zlog.CtxErrorf(ctx, "查询发起用户失败: %v, userId: %d", err, req.UserId)
		return resp, response.ErrResponse(err, response.USER_NOT_EXIST)
	}

	// 查询目标用户
	targetUser, err := request.NewUserRequest(global.DB).FindUserByName(req.TargetName)
	if err != nil {
		zlog.CtxErrorf(ctx, "查询目标用户失败: %v, username: %s", err, req.TargetName)
		return resp, response.ErrResponse(err, response.USER_NOT_EXIST)
	}

	// 检查是否添加自己
	if targetUser.ID == req.UserId {
		zlog.CtxErrorf(ctx, "不能添加自己为好友, userId: %d", req.UserId)
		return resp, response.ErrResponse(err, response.ME_AND_ME)
	}

	// 检查是否已是好友
	contact, err := request.NewContactRequest(global.DB).IsFriend(req.UserId, targetUser.ID)
	if err != nil {
		zlog.CtxErrorf(ctx, "检查好友关系失败: %v", err)
		return resp, response.ErrResponse(err, response.COMMON_FAIL)
	}
	if contact.ID != 0 {
		zlog.CtxErrorf(ctx, "已是好友, 不能重复添加")
		return resp, response.ErrResponse(err, response.FRIEND_YES_FRIEN)
	}

	// 添加好友
	// 添加详细日志记录
	zlog.CtxDebugf(ctx, "准备添加好友: ownerId=%d, targetId=%d", ownerUser.ID, targetUser.ID)

	// 添加事务重试机制
	var retryCount = 0
	maxRetries := 3
	var lastErr error

	for retryCount < maxRetries {
		if err := request.NewContactRequest(global.DB).AddFriend(ownerUser, targetUser); err != nil {
			lastErr = err
			zlog.CtxErrorf(ctx, "添加好友失败(尝试 %d/%d): %v", retryCount+1, maxRetries, err)
			retryCount++
			continue
		}
		break
	}

	if retryCount == maxRetries {
		zlog.CtxErrorf(ctx, "添加好友最终失败: %v", lastErr)
		return resp, response.ErrResponse(lastErr, response.COMMON_FAIL)
	}

	zlog.CtxInfof(ctx, "添加好友成功: ownerId=%d, targetId=%d", ownerUser.ID, targetUser.ID)
	return resp, nil
}

func (l *ContactLogic) SearchFriend(ctx context.Context, req types.SearchFriendRequest) (resp types.SearchFriendResponse, err error) {
	users, err := request.NewContactRequest(global.DB).SearchFriend(req.UserId)
	if len(users) == 0 {
		zlog.CtxErrorf(ctx, "搜索朋友失败: %v", err)
	}
	resp.Users = users
	return resp, nil
}

func (l *ContactLogic) SearchUserByGroupId(ctx context.Context, req types.SearchUserByGroupIdRequest) (resp types.SearchUserByGroupIdResponse, err error) {
	contacts := make([]model.Contact, 0)
	objIds := make([]int64, 0)
	contacts = request.NewContactRequest(global.DB).SearchUserByGroupId(req.CommunityId)
	for _, v := range contacts {
		objIds = append(objIds, v.OwnerId)
	}
	resp.UserIds = objIds
	zlog.CtxDebugf(ctx, "查找群友成功: %v", req)
	return resp, err
}
