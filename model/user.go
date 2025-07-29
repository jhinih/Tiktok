package model

import "time"

type User struct {
	ID int64 `json:"id" gorm:"column:id;primaryKey;type:bigint"`
	TimeModel
	Username string `json:"username" gorm:"column:username;type:varchar(255);comment:用户名"`
	Password string `json:"password" gorm:"column:password;type:varchar(255);comment:密码"`
	Email    string `json:"email" gorm:"column:email;type:varchar(255);comment:邮箱"`
	Avatar   string `json:"avatar" gorm:"column:avatar;type:varchar(255);comment:头像URL"`
	CreateAt int64  `json:"create_at" gorm:"column:create_at;type:bigint"`
	Role     int    `json:"role" gorm:"column:role;type:int;comment:权限等级"`

	Phone         string    `gorm:"column:phone;type:varchar(20);comment:手机号" valid:"matches(^1[3-9]{1}\\d{9}$)"`
	ClientIp      string    `gorm:"column:client_ip;type:varchar(50);comment:客户端IP"`
	ClientPort    string    `gorm:"column:client_port;type:varchar(20);comment:客户端端口"`
	LoginTime     time.Time `gorm:"column:login_time;comment:登录时间"`
	HeartbeatTime time.Time `gorm:"column:heartbeat_time;comment:心跳时间"`
	LoginOutTime  time.Time `gorm:"column:login_out_time;comment:登出时间" json:"login_out_time"`
	IsLogout      bool      `gorm:"column:is_logout;comment:是否登出"`
	DeviceInfo    string    `gorm:"column:device_info;type:varchar(255);comment:设备信息"`
	Bio           string    `gorm:"column:bio;type:varchar(255);comment:个人简介"`
	// 0: 游客(未实名) 1:普通用户 2.正式成员 3:管理员 4:超级管理员
}
