package types

type SendCodeRequest struct {
	Email string `json:"email"`
}

type SendCodeResponse struct {
	Ok bool `json:"ok"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Code     string `json:"code"`
	Password string `json:"password"`
	Username string `json:"username"`

	Avatar string `json:"avatar" gorm:"column:avatar;type:varchar(255);comment:头像URL"`

	//Phone         string `gorm:"column:phone;type:varchar(20);comment:手机号" valid:"matches(^1[3-9]{1}\\d{9}$)"`
	//ClientIp      string `gorm:"column:client_ip;type:varchar(50);comment:客户端IP"`
	//ClientPort    string `gorm:"column:client_port;type:varchar(20);comment:客户端端口"`
	//DeviceInfo    string `gorm:"column:device_info;type:varchar(255);comment:设备信息"`
	//Bio			string `gorm:"column:bio;type:varchar(255);comment:个人简介"`
}

type RegisterResponse struct {
	Atoken string `json:"atoken"`
	Ok     bool   `json:"ok"`
}

type LoginRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	IsRemember bool   `json:"is_remember"`
}

type LoginResponse struct {
	Atoken string `json:"atoken"`
	Rtoken string `json:"rtoken"`
	Ok     bool   `json:"ok"`
}

type RefreshTokenRequest struct {
	Rtoken string `json:"rtoken"`
}

type RefreshTokenResponse struct {
	Atoken string `json:"atoken"`
}

//type TokenTestRequest struct {
//}
//
//type TokenTestResponse struct {
//}

//type RefreshTokenRequest struct {
//	Rtoken string `json:"rtoken"`
//}
//
//type RefreshTokenResponse struct {
//	Atoken string `json:"atoken"`
//}

type TokenTestRequest struct {
}

type TokenTestResponse struct {
}

//type LoginResponse struct {
//	Token    string `json:"atoken"`
//	UserInfo interface{}
//}
