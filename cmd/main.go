package main

import (
	"Tiktok/initialize"
	"Tiktok/log/zlog"
	routerg "Tiktok/router"
)

func main() {
	// 初始化
	initialize.Init()

	// 工程进入前夕，释放资源
	defer initialize.Eve()

	// 运行服务
	routerg.RunServer()
	zlog.Infof("程序运行完成！")
}
