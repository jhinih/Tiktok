package api

import (
	"Tiktok/log/zlog"
	"Tiktok/logic"
	"Tiktok/response"
	"Tiktok/types"
	"Tiktok/utils/jwtUtils"
	"github.com/gin-gonic/gin"
)

// GetUserInfo 获取用户基础信息
func GetUserInfo(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.GetUserInfoRequest](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "获取用户基础信息请求: %v", req)
	resp, err := logic.NewUserLogic().GetUserInfo(ctx, req)
	response.Response(c, resp, err)
}

// GetMyUserInfo 获取自己的用户基础信息
func GetMyUserInfo(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.GetUserInfoRequest](c)
	if err != nil {
		return
	}
	// 直接从token中获取用户ID，然后调用UserInfo接口
	req.ID = jwtUtils.GetUserId(c)
	zlog.CtxInfof(ctx, "获取自己的用户基础信息请求: %v", req)
	resp, err := logic.NewUserLogic().GetUserInfo(ctx, req)
	response.Response(c, resp, err)
}

// GetProfile 获取用户资料
func GetProfile(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.GetUserProfileRequest](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "获取用户资料请求: %v", req)
	resp, err := logic.NewUserLogic().GetUserProfile(ctx, req)
	response.Response(c, resp, err)
}

// SetProfile 设置用户资料
func SetProfile(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.SetUserProfileRequest](c)
	if err != nil {
		return
	}
	// 设置者
	req.OperatorID = jwtUtils.GetUserId(c)
	req.OperatorRole = jwtUtils.GetRole(c)
	zlog.CtxInfof(ctx, "解析token成功，role: %v", req.OperatorRole)
	zlog.CtxInfof(ctx, "修改用户资料请求: %v", req)
	resp, err := logic.NewUserLogic().SetUserProfile(ctx, req)
	response.Response(c, resp, err)
}

// SetRole 设置用户角色
func SetRole(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.SetUserRoleRequest](c)
	if err != nil {
		return
	}
	// 设置者
	req.OperatorRole = jwtUtils.GetRole(c)
	zlog.CtxInfof(ctx, "解析token成功，role: %v", req.OperatorRole)
	zlog.CtxInfof(ctx, "修改用户权限请求: %v", req)
	resp, err := logic.NewUserLogic().SetUserRole(ctx, req)
	response.Response(c, resp, err)
}
