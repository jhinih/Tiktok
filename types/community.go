package types

type CreateCommunityRequest struct {
	OwnerId   string `json:"ownerId"`
	OwnerName string `json:"ownerName"`
	Name      string `json:"name"`
	Icon      string `json:"icon"`
	Desc      string `json:"desc"`
}

type CreateCommunityResponse struct {
}
type LoadCommunityRequest struct {
	OwnerId string `json:"ownerId"`
}

type LoadCommunityResponse struct {
}

type JoinGroupsRequest struct {
	UserId string `json:"userId"`
	ComId  string `json:"comId"`
}

type JoinGroupsResponse struct {
}
