package initialize

import (
	"Tiktok/global"
	"Tiktok/log/zlog"
	"Tiktok/utils/snowflakeUtils"
)

func InitSnowflake() {
	var err error
	global.SnowflakeNode, err = snowflakeUtils.NewNode(global.DEFAULT_NODE_ID)
	if err != nil {
		zlog.Errorf("初始化雪花ID生成节点失败: %v", err)
		panic(err)
	}
}
