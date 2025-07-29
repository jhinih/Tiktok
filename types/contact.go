package types

import "Tiktok/model"

type AddFriendRequest struct {
	UserId     string `json:"userId"`
	TargetName string `json:"username"`
}

type AddFriendResponse struct {
}
type SearchFriendRequest struct {
	UserId string `json:"userId"`
}

type SearchFriendResponse struct {
	Users []model.User `json:"users"`
}
type SearchUserByGroupIdRequest struct {
	CommunityId string `json:"communityId"`
}

type SearchUserByGroupIdResponse struct {
	UserIds []int64 `json:"userIds"`
}
