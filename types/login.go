package types

type SendCodeRequest struct {
	Email string `json:"email"`
}

type SendCodeResponse struct {
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Code     string `json:"code"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type RegisterResponse struct {
	Atoken string `json:"atoken"`
}

type LoginRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	IsRemember bool   `json:"is_remember"`
}

type LoginResponse struct {
	Atoken string `json:"atoken"`
	Rtoken string `json:"rtoken"`
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
