package api

import (
	"Tiktok/log/zlog"
	"Tiktok/logic"
	"Tiktok/models"
	"Tiktok/response"
	"Tiktok/types"
	"Tiktok/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

// 加好友
func AddFriend(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.AddFriendRequest](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "注册请求: %v", req)
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
	zlog.CtxInfof(ctx, "注册请求: %v", req)
	resp, err := logic.NewContactLogic().SearchFriend(c, req)
	response.Response(c, resp, err)
}
func SendUserMsg(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}
func SendMsg(c *gin.Context) {
	// 防止跨域站点伪造请求
	var upGrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(ws)
	MsgHandler(c, ws)
}
func MsgHandler(c *gin.Context, ws *websocket.Conn) {
	for {
		msg, err := utils.Subscribe(c, utils.PublishKey)
		if err != nil {
			fmt.Println(" MsgHandler 发送失败", err)
		}

		tm := time.Now().Format("2006-01-02 15:04:05")
		m := fmt.Sprintf("[ws][%s]:%s", tm, msg)
		err = ws.WriteMessage(1, []byte(m))
		if err != nil {
			log.Fatalln(err)
		}
	}
}
func GetUserList(c *gin.Context) {
	data := make([]*models.UserBasic, 10)
	data = models.GetUserList()
	c.JSON(200, gin.H{
		"code":    0, //  0成功   -1失败
		"message": "用户名已注册！",
		"data":    data,
	})
}
