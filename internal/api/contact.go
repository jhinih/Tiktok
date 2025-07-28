package api

import (
	"Tiktok/log/zlog"
	"Tiktok/logic"
	"Tiktok/response"
	"Tiktok/types"
	"github.com/gin-gonic/gin"
)

// 加好友
func AddFriend(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.AddFriendRequest](c)
	if err != nil {
		return
	}
	// 从token中获取用户ID
	zlog.CtxInfof(ctx, "加好友请求: %v", req)
	resp, err := logic.NewContactLogic().AddFriend(ctx, req)
	response.Response(c, resp, err)
}

// 搜索好友
func SearchFriend(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.SearchFriendRequest](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "搜索好友请求: %v", req)
	resp, err := logic.NewContactLogic().SearchFriend(c, req)
	response.Response(c, resp, err)
}
