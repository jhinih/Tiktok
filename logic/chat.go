package logic

import (
	"Tiktok/global"
	"Tiktok/log/zlog"
	"Tiktok/model"
	"Tiktok/request"
	"Tiktok/types"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
	"io"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ChatLogic struct {
}

func NewChatLogic() *ChatLogic {
	return &ChatLogic{}
}

// 映射关系
var clientMap map[int64]*model.Node = make(map[int64]*model.Node, 0)

// 读写锁
var rwLocker sync.RWMutex

// 需要 ：发送者ID ，接受者ID ，消息类型，发送的内容，发送类型
func (l *ChatLogic) Chat(c *gin.Context, req types.ChatRequest) (resp types.ChatResponse, err error) {
	userId, _ := strconv.ParseInt(strconv.Itoa(int(req.UserId)), 10, 64)
	isvalida := true //checkToke()  待.........
	conn, err := (&websocket.Upgrader{
		//token 校验
		CheckOrigin: func(r *http.Request) bool {
			return isvalida
		},
	}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	//2.获取conn
	currentTime := uint64(time.Now().Unix())
	node := &model.Node{
		Conn:          conn,
		Addr:          conn.RemoteAddr().String(), //客户端地址
		HeartbeatTime: currentTime,                //心跳时间
		LoginTime:     currentTime,                //登录时间
		DataQueue:     make(chan []byte, 50),
		GroupSets:     set.New(set.ThreadSafe),
	}
	//3. 用户关系
	//4. userid 跟 node绑定 并加锁
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()
	//5.完成发送逻辑
	go NewChatLogic().sendProc(node)
	//6.完成接受逻辑
	go NewChatLogic().recvProc(node)
	//7.加入在线用户到缓存
	model.SetUserOnlineInfo("online_"+strconv.Itoa(int(req.UserId)), []byte(node.Addr), time.Duration(viper.GetInt("timeout.RedisOnlineTime"))*time.Hour)

	//sendMsg(userId, []byte("欢迎进入聊天系统"))
	return
}

func (l *ChatLogic) sendProc(node *model.Node) {
	for {
		select {
		case data := <-node.DataQueue:
			fmt.Println("[ws]sendProc >>>> msg :", string(data))
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func (l *ChatLogic) recvProc(node *model.Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		msg := model.Message{}
		err = json.Unmarshal(data, &msg)
		if err != nil {
			fmt.Println(err)
		}
		//心跳检测 msg.Media == -1 || msg.Type == 3
		//if msg.Type == 3 {
		//	currentTime := uint64(time.Now().Unix())
		//	node.Heartbeat(currentTime)
		//} else {
		NewChatLogic().dispatch(data)
		NewChatLogic().broadMsg(data) //todo 将消息广播到局域网
		fmt.Println("[ws] recvProc <<<<< ", string(data))
		//}

	}
}

var udpsendChan chan []byte = make(chan []byte, 1024)

func (l *ChatLogic) broadMsg(data []byte) {
	udpsendChan <- data
}

func (l *ChatLogic) init() {
	go NewChatLogic().udpSendProc()
	go NewChatLogic().udpRecvProc()
	fmt.Println("init goroutine ")
}

// 完成udp数据发送协程
func (l *ChatLogic) udpSendProc() {
	con, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(192, 168, 0, 255),
		Port: viper.GetInt("port.udp"),
	})
	defer con.Close()
	if err != nil {
		fmt.Println(err)
	}
	for {
		select {
		case data := <-udpsendChan:
			fmt.Println("udpSendProc  data :", string(data))
			_, err := con.Write(data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}

}

// 完成udp数据接收协程
func (l *ChatLogic) udpRecvProc() {
	con, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: viper.GetInt("port.udp"),
	})
	if err != nil {
		fmt.Println(err)
	}
	defer con.Close()
	for {
		var buf [512]byte
		n, err := con.Read(buf[0:])
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("udpRecvProc  data :", string(buf[0:n]))
		NewChatLogic().dispatch(buf[0:n])
	}
}

// 后端调度逻辑处理
func (l *ChatLogic) dispatch(data []byte) {
	msg := model.Message{}
	msg.CreateTime = uint64(time.Now().Unix())
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch msg.Type {
	case 1: //私信
		fmt.Println("dispatch  data :", string(data))
		//sendMsg(msg.TargetId, data)
		req := types.SendMsgRequest{
			UserId: msg.TargetId,
			Msg:    data,
		}
		NewChatLogic().SendMsg(req)
	case 2: //群发
		//sendGroupMsg(msg.TargetId, data) //发送的群ID ，消息内容
		req := types.SendGroupMsgRequest{
			TargetId: msg.TargetId,
			Msg:      data,
		}
		NewChatLogic().SendGroupMsg(req)
		// case 4: // 心跳
		// 	node.Heartbeat()
		//case 4:
		//
	}
}
func (l *ChatLogic) SendGroupMsg(req types.SendGroupMsgRequest) (resp types.SendGroupMsgResponse, err error) {
	fmt.Println("开始群发消息")
	//userIds := SearchUserByGroupId(uint(req.TargetId))
	//var contact []model.Contact
	contacts, _ := request.NewChatRequest(global.DB).SearchUserByGroupId(uint(req.TargetId))
	userIds := make([]uint, 0)
	for _, v := range contacts {
		userIds = append(userIds, uint(v.OwnerId))
	}
	for i := 0; i < len(userIds); i++ {
		//排除给自己的
		if req.TargetId != int64(userIds[i]) {
			//sendMsg(int64(userIds[i]), msg)
			req := types.SendMsgRequest{
				UserId: req.TargetId,
				Msg:    req.Msg,
			}
			_, err := NewChatLogic().SendMsg(req)
			return resp, err

		}

	}
	return resp, nil
}
func (l *ChatLogic) SendMsg(req types.SendMsgRequest) (resp types.SendMsgResponse, err error) {
	rwLocker.RLock()
	node, ok := clientMap[req.UserId]
	rwLocker.RUnlock()
	jsonMsg := model.Message{}
	json.Unmarshal(req.Msg, &jsonMsg)
	ctx := context.Background()
	targetIdStr := strconv.Itoa(int(req.UserId))
	userIdStr := strconv.Itoa(int(jsonMsg.UserId))
	jsonMsg.CreateTime = uint64(time.Now().Unix())
	r, err := global.Rdb.Get(ctx, "online_"+userIdStr).Result()
	if err != nil {
		fmt.Println(err)
	}
	if r != "" {
		if ok {
			fmt.Println("sendMsg >>> userID: ", req.UserId, "  msg:", string(req.Msg))
			node.DataQueue <- req.Msg
		}
	}
	var key string
	if req.UserId > jsonMsg.UserId {
		key = "msg_" + userIdStr + "_" + targetIdStr
	} else {
		key = "msg_" + targetIdStr + "_" + userIdStr
	}

	res, err := global.Rdb.ZRevRange(ctx, key, 0, -1).Result()
	if err != nil {
		fmt.Println(err)
	}
	score := float64(cap(res)) + 1
	ress, e := global.Rdb.ZAdd(ctx, key, redis.Z{Score: score, Member: string(req.Msg)}).Result()
	//res, e := utils.Red.Do(ctx, "zadd", key, 1, jsonMsg).Result() //备用 后续拓展 记录完整msg
	if e != nil {
		fmt.Println(e)
	}
	fmt.Println(ress)
	return resp, nil
}

//// 需要重写此方法才能完整的msg转byte[]
//func (msg model.Message) MarshalBinary() ([]byte, error) {
//	return json.Marshal(msg)
//}

// 获取缓存里面的消息
func (l *ChatLogic) RedisMsg(ctx context.Context, req types.RedisMsgRequest) (resp types.RedisMsgResponse, err error) {
	//func RedisMsg(userIdA int64, userIdB int64, start int64, end int64, isRev bool) []string {
	rwLocker.RLock()
	rwLocker.RUnlock()
	ctx = context.Background()
	userIdStr := strconv.Itoa(int(req.UserIdA))
	targetIdStr := strconv.Itoa(int(req.UserIdB))
	var key string
	if req.UserIdA > req.UserIdB {
		key = "msg_" + targetIdStr + "_" + userIdStr
	} else {
		key = "msg_" + userIdStr + "_" + targetIdStr
	}
	//key = "msg_" + userIdStr + "_" + targetIdStr
	//rels, err := utils.Red.ZRevRange(ctx, key, 0, 10).Result()  //根据score倒叙

	var rels []string
	if req.IsRev {
		rels, err = global.Rdb.ZRange(ctx, key, req.Start, req.End).Result()
	} else {
		rels, err = global.Rdb.ZRevRange(ctx, key, req.Start, req.End).Result()
	}
	if err != nil {
		fmt.Println(err) //没有找到
	}
	// 发送推送消息
	/**
	// 后台通过websoket 推送消息
	for _, val := range rels {
		fmt.Println("sendMsg >>> userID: ", userIdA, "  msg:", val)
		node.DataQueue <- []byte(val)
	}**/
	resp.Rels = rels
	return resp, nil
}

//// 更新用户心跳
//func (node *model.Node) Heartbeat(currentTime uint64) {
//	node.HeartbeatTime = currentTime
//	return
//}

//// 清理超时连接
//func CleanConnection(param interface{}) (result bool) {
//	result = true
//	defer func() {
//		if r := recover(); r != nil {
//			fmt.Println("cleanConnection err", r)
//		}
//	}()
//	//fmt.Println("定时任务,清理超时连接 ", param)
//	//node.IsHeartbeatTimeOut()
//	currentTime := uint64(time.Now().Unix())
//	for i := range clientMap {
//		node := clientMap[i]
//		if node.IsHeartbeatTimeOut(currentTime) {
//			fmt.Println("心跳超时..... 关闭连接：", node)
//			node.Conn.Close()
//		}
//	}
//	return result
//}
//
//// 用户心跳是否超时
//func (node *model.Node) IsHeartbeatTimeOut(currentTime uint64) (timeout bool) {
//	if node.HeartbeatTime+viper.GetUint64("timeout.HeartbeatMaxTime") <= currentTime {
//		fmt.Println("心跳超时。。。自动下线", node)
//		timeout = true
//	}
//	return
//}

func (l *ChatLogic) GetUserList() ([]map[string]interface{}, error) {
	var users []model.User
	if err := request.NewChatRequest(global.DB).SearchUsers(&users); err != nil {
		if err == gorm.ErrRecordNotFound {
			return []map[string]interface{}{}, nil
		}
		return nil, fmt.Errorf("failed to get user list: %w", err)
	}

	result := make([]map[string]interface{}, 0, len(users))
	for _, user := range users {
		result = append(result, map[string]interface{}{
			"id":     user.ID,
			"name":   user.Username,
			"avatar": user.Avatar,
		})
	}

	if len(result) == 0 {
		return []map[string]interface{}{}, nil
	}

	return result, nil
}

// 上传文件到本地
func (l *ChatLogic) Upload(c *gin.Context) (int int, err error) {

	//func Upload(c *gin.Context) {
	w := c.Writer
	r := c.Request
	srcFile, head, err := r.FormFile("file")
	suffix := ".png"
	ofilName := head.Filename
	tem := strings.Split(ofilName, ".")
	if len(tem) > 1 {
		suffix = "." + tem[len(tem)-1]
	}
	fileName := fmt.Sprintf("%d%04d%s", time.Now().Unix(), rand.Int31(), suffix)
	dstFile, err := os.Create("./upload/asset" + fileName)

	_, err = io.Copy(dstFile, srcFile)

	url := "./upload/asset/" + fileName
	zlog.CtxDebugf(c, url, w)
	return 666, err
}
