package api

import (
	"Tiktok/log/zlog"
	"Tiktok/logic"
	"Tiktok/response"
	"Tiktok/types"
	"github.com/gin-gonic/gin"
)

func Chats(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.ChatRequest](c)
	if err != nil {
		response.Response(c, nil, err)
		return
	}
	zlog.CtxInfof(ctx, "Chat request: %v", req)

	// Implement chat logic here
	resp, err := logic.NewChatLogic().Chat(c, req)
	response.Response(c, resp, err)
}

func RedisMsg(c *gin.Context) {

	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.RedisMsgRequest](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "发送个人消息请求: %v", req)
	resp, err := logic.NewChatLogic().RedisMsg(ctx, req)
	response.Response(c, resp, err)
}

func GetUserList(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)

	users, err := logic.NewChatLogic().GetUserList()
	if err != nil {
		zlog.CtxErrorf(ctx, "GetUserList failed: %v", err)
		response.Response(c, nil, err)
		return
	}

	zlog.CtxDebugf(ctx, "GetUserList success, count: %d", len(users))
	response.Response(c, gin.H{
		"data": users,
		"code": 0,
		"msg":  "success",
	}, nil)
}
func SendMsg(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.SendMsgRequest](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "发送消息请求: %v", req)
	resp, err := logic.NewChatLogic().SendMsg(req)
	response.Response(c, resp, err)
}
func SendUserMsg(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.SendMsgRequest](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "发送个人消息请求: %v", req)
	resp, err := logic.NewChatLogic().SendMsg(req)
	response.Response(c, resp, err)
}
func SendGroupMsg(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.SendGroupMsgRequest](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "发送群消息请求: %v", req)
	resp, err := logic.NewChatLogic().SendGroupMsg(req)
	response.Response(c, resp, err)
}

func Upload(c *gin.Context) {
	resp, err := logic.NewChatLogic().Upload(c)
	response.Response(c, resp, err)
}
