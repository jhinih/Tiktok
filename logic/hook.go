package logic

import (
	"Tiktok/log/zlog"
	//"Tiktok/model"
	"gorm.io/gorm"
)

// RegisterHook 注册 GORM 钩子
func RegisterHook(db *gorm.DB) {
	zlog.Infof("注册 GORM hooks...")
	//db.Callback().Create().Before("gorm:Create").Register("before_create_BaseModel", BeforeCreateBaseModel)
}

//func BeforeCreateBaseModel(db *gorm.DB) {
//	if db.Statement.Schema != nil {
//		if baseModel, ok := db.Statement.Model.(*model.CommonModel); ok {
//			baseModel.BeforeCreate(db)
//		}
//		if baseModel, ok := db.Statement.Model.(*model.TimeModel); ok {
//			baseModel.BeforeCreate(db)
//		}
//	}
//}
