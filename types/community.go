package types

type CreateCommunityRequest struct {
	OwnerId   int64  `json:"ownerId"`
	OwnerName string `json:"ownerName"`
	Name      string `json:"name"`
	Icon      string `json:"icon"`
	Desc      string `json:"desc"`
}

type CreateCommunityResponse struct {
}
type LoadCommunityRequest struct {
	OwnerId int64 `json:"ownerId"`
}

type LoadCommunityResponse struct {
}

type JoinGroupsRequest struct {
	UserId int64 `json:"userId"`
	ComId  int64 `json:"comId"`
}

type JoinGroupsResponse struct {
}
