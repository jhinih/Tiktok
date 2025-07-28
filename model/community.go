package model

type Community struct {
	ID int64 `json:"id,string" gorm:"column:id;primaryKey;type:bigint;type:bigint"`
	TimeModel
	Name      string `json:"name" gorm:"column:name;type:varchar(255);size:255;unique;not null;"`
	OwnerId   int64  `json:"owner_id" gorm:"column:owner_id;type:bigint;type:bigint"`
	OwnerName string
	Img       string
	Desc      string
}
