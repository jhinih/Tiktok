package types

type CreateCommunityRequest struct {
	OwnerId int    `json:"ownerId"`
	Name    string `json:"name"`
	Icon    string `json:"icon"`
	Desc    string `json:"desc"`
}

type CreateCommunityResponse struct {
	Ok bool `json:"ok"`
}
type LoadCommunityRequest struct {
	OwnerId int `json:"ownerId"`
}

type LoadCommunityResponse struct {
	Ok bool `json:"ok"`
}

type JoinGroupsRequest struct {
	UserId int `json:"userId"`
	ComId  int `json:"comId"`
}

type JoinGroupsResponse struct {
	Ok bool `json:"ok"`
}
