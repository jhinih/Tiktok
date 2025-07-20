package database

import (
	"Tiktok/configs"
	"Tiktok/global"
	"Tiktok/log/zlog"
	"gorm.io/gorm"
)

type DataBase interface {
	GetDsn(config configs.Config) string
	InitDataBase(config configs.Config) (*gorm.DB, error)
}

func InitDataBases(base DataBase, config configs.Config) {
	var err error
	global.DB, err = base.InitDataBase(config)
	if err != nil {
		zlog.Fatalf("无法初始化数据库 %v", err)
		return
	}
	zlog.Infof("初始化数据库成功！")
	////对该数据库注册 hook
	//logic.RegisterHook(global.DB)
	return
}
