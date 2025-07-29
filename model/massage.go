package model

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	UserId     string //发送者
	UserName   string
	TargetId   string //接受者
	TargetName string
	Type       int    //发送类型  1私聊  2群聊  3心跳
	Media      int    //消息类型  1文字 2表情包 3语音 4图片 /表情包
	Content    string //消息内容
	CreateTime uint64 //创建时间
	ReadTime   uint64 //读取时间
	Pic        string
	Url        string
	Desc       string
	Amount     int //其他数字统计
}

func (table *Message) TableName() string {
	return "message"
}

// const (
// 	HeartbeatMaxTime = 1 * 60
// )

type Node struct {
	Conn          *websocket.Conn //连接
	Addr          string          //客户端地址
	FirstTime     uint64          //首次连接时间
	HeartbeatTime uint64          //心跳时间
	LoginTime     uint64          //登录时间
	DataQueue     chan []byte     //消息
	GroupSets     set.Interface   //好友 / 群
}

func (table *Node) TableName() string {
	return "node"
}

// 更新用户心跳
func (node Node) Heartbeat(currentTime uint64) {
	node.HeartbeatTime = currentTime
	return
}

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

// 用户心跳是否超时
func (node *Node) IsHeartbeatTimeOut(currentTime uint64) (timeout bool) {
	if node.HeartbeatTime+viper.GetUint64("timeout.HeartbeatMaxTime") <= currentTime {
		fmt.Println("心跳超时。。。自动下线", node)
		timeout = true
	}
	return
}
