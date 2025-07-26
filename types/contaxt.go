package types

type AddFriendRequest struct {
	UserId     uint
	TargetName string
}

type AddFriendResponse struct {
	Ok bool `json:"ok"`
}
type SearchFriendRequest struct {
	UserId uint
}

type SearchFriendResponse struct {
	Ok bool `json:"ok"`
}
type ChatRequest struct {
	UserId uint
}

type ChatResponse struct {
	Ok bool `json:"ok"`
}
