package types

type GetUserInfoRequest struct {
	ID string `form:"id"`
}

type GetUserInfoResponse struct {
	ID       int64  `json:"id,string"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Role     int    `json:"role"`
}

type GetUserProfileRequest struct {
	ID string `form:"id"`
}

type GetUserProfileResponse struct {
	ID       int64  `json:"id,string"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Role     int    `json:"role"`
}

type SetUserProfileRequest struct {
	OperatorID   string `json:"-"`
	OperatorRole int    `json:"-"`

	ID string `json:"id"`

	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Bio      string `json:"bio"`
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
