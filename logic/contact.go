package logic

import (
	"Tiktok/global"
	"Tiktok/log/zlog"
	"Tiktok/request"
	"Tiktok/response"
	"Tiktok/types"
	"context"
	"strconv"
	"strings"
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
	clean := strings.Trim(req.UserId, `"`)
	Id, err := strconv.ParseInt(clean, 10, 64)
	if err != nil {
		zlog.CtxErrorf(ctx, "用户ID转换失败: %v, userId: %s", err, req.UserId)
		return resp, response.ErrResponse(err, response.PARAM_IS_INVALID)
	}
	ownerUser, err := request.NewUserRequest(global.DB).FindUserById(int64(Id))
	if err != nil {
		zlog.CtxErrorf(ctx, "查询发起用户失败: %v, userId: %s", err, req.UserId)
		return resp, response.ErrResponse(err, response.USER_NOT_EXIST)
	}

	// 查询目标用户
	targetUser, err := request.NewUserRequest(global.DB).FindUserByName(req.TargetName)
	if err != nil {
		zlog.CtxErrorf(ctx, "查询目标用户失败: %v, username: %s", err, req.TargetName)
		return resp, response.ErrResponse(err, response.USER_NOT_EXIST)
	}

	// 检查是否添加自己
	if targetUser.ID == int64(Id) {
		zlog.CtxErrorf(ctx, "不能添加自己为好友, userId: %d", req.UserId)
		return resp, response.ErrResponse(err, response.ME_AND_ME)
	}

	// 检查是否已是好友
	contact, err := request.NewContactRequest(global.DB).IsFriend(int64(Id), int64(targetUser.ID))
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
	clean := strings.Trim(req.UserId, `"`)
	Id, _ := strconv.ParseInt(clean, 10, 64)
	users, err := request.NewContactRequest(global.DB).SearchFriend(int64(Id))
	if len(users) == 0 {
		zlog.CtxErrorf(ctx, "搜索朋友失败: %v", err)
	}
	resp.Users = users
	return resp, nil
}

//群是否被拥有
//func (l *ContactLogic) SearchUserByGroupId(ctx context.Context, req types.SearchUserByGroupIdRequest) (resp types.SearchUserByGroupIdResponse, err error) {
//	contacts := make([]model.Contact, 0)
//	objIds := make([]int64, 0)
//	clean := strings.Trim(req.CommunityId, `"`)
//	Id, _ := strconv.ParseInt(clean, 10, 64)
//	contacts = request.NewContactRequest(global.DB).SearchUserByGroupId(int64(Id))
//	for _, v := range contacts {
//		objIds = append(objIds, v.OwnerId)
//	}
//	resp.UserIds = objIds
//	zlog.CtxDebugf(ctx, "查找群友成功: %v", req)
//	return resp, err
//}
