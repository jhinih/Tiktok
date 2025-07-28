package types

import "Tiktok/model"

type AddFriendRequest struct {
	UserId     int64  `json:"userId"`
	TargetName string `json:"username"`
}

type AddFriendResponse struct {
}
type SearchFriendRequest struct {
	UserId int64 `json:"userId"`
}

type SearchFriendResponse struct {
	Users []model.User `json:"users"`
}
type SearchUserByGroupIdRequest struct {
	CommunityId int64 `json:"communityId"`
}

type SearchUserByGroupIdResponse struct {
	UserIds []int64 `json:"userIds"`
}
