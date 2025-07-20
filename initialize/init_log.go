package initialize

import (
	"Tiktok/configs"
	"Tiktok/log"
	"Tiktok/log/zlog"
)

func InitLog(config *configs.Config) {
	logger := log.GetZap(config)
	zlog.InitLogger(logger)
}
