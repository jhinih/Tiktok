package api

import (
	"Tiktok/log/zlog"
	"Tiktok/logic"
	"Tiktok/response"
	"Tiktok/types"
	"github.com/gin-gonic/gin"
)

func CreateCommunity(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.CreateCommunityRequest](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "登录请求: %v", req)
	resp, err := logic.NewCommunityLogic().CreateCommunity(ctx, req)
	response.Response(c, resp, err)
}

func LoadCommunity(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.LoadCommunityRequest](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "登录请求: %v", req)
	resp, err := logic.NewCommunityLogic().LoadCommunity(ctx, req)
	response.Response(c, resp, err)
}
func JoinGroups(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.JoinGroupsRequest](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "登录请求: %v", req)
	resp, err := logic.NewCommunityLogic().JoinGroups(ctx, req)
	response.Response(c, resp, err)
}
