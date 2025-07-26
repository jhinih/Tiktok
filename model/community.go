package model

type Community struct {
	ID int64 `json:"id" gorm:"column:id;primaryKey;type:bigint;type:bigint"`
	TimeModel
	Name    string
	OwnerId uint
	Img     string
	Desc    string
}
