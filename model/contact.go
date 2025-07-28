package model

type Contact struct {
	ID int64 `json:"id,string" gorm:"column:id;primaryKey;type:bigint;type:bigint"`
	TimeModel
	OwnerId    int64 //谁的关系信息
	OwnerName  string
	TargetId   int64 //对应的谁 /群 ID
	TargetName string

	Type int //对应的类型  1好友  2群  3xx
	Desc string
}

func (table *Contact) TableName() string {
	return "contact"
}
