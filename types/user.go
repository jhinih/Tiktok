package types

type GetUserInfoRequest struct {
	ID string `form:"id"`
}

type GetUserInfoResponse struct {
	ID       string `json:"id,string"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Role     int    `json:"role"`
	Bio      string `json:"bio"`
	CreateAt int64  `json:"create_at,string"`
	Email    string `json:"email"`
}

type GetUserProfileRequest struct {
	ID string `form:"id"`
}

type GetUserProfileResponse struct {
	ID       string `json:"id,string"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Role     int    `json:"role"`
	Bio      string `json:"bio"`
	Email    string `json:"email"`
	CreateAt int64  `json:"create_at,string"`
}

type SetUserProfileRequest struct {
	OperatorID   string `json:"-"`
	OperatorRole int    `json:"-"`

	ID string `json:"id"`

	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Bio      string `json:"bio"`
	Email    string `json:"email"`
}

type SetUserProfileResponse struct {
}

type SetUserRoleRequest struct {
	OperatorRole int    `json:"-"`
	ID           string `json:"id"`
	Role         int    `json:"role"`
}

type SetUserRoleResp struct {
}
