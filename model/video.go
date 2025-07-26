package model

import "time"

// Video 视频数据结构
type Video struct {
	ID          int64  `json:"id" gorm:"column:id;primaryKey;type:bigint;type:bigint"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Cover       string `json:"cover"`
	Likes       int    `json:"likes"`
	Comments    int    `json:"comments"`
	Shares      int    `json:"shares"`
	UserID      int64  `json:"user_id" gorm:"column:user_id;type:bigint;comment:用户ID"`
	CreatedAt   string `json:"created_at"`
	PublishTime time.Time
	Type        string `json:"type" gorm:"column:type;type:varchar(63);comment:类型"`

	IsPrivate bool `json:"is_private" gorm:"not null;column:is_private;type:bool;comment:是否私密"`
}

func (Video) TableName() string {
	return "videos"
}

type VideoLike struct {
	ID int64 `json:"id" gorm:"column:id;primaryKey;type:bigint;type:bigint"`
	TimeModel

	VideoID int64 `json:"video_id" gorm:"column:video_id;type:bigint;comment:视频ID;uniqueIndex:idx_video_user_unique"`
	UserID  int64 `json:"user_id" gorm:"column:user_id;type:bigint;comment:用户ID;uniqueIndex:idx_video_user_unique"`
}

func (VideoLike) TableName() string {
	return "video_like"
}

type Comment struct {
	ID int64 `json:"id" gorm:"column:id;primaryKey;type:bigint;type:bigint"`
	TimeModel

	Content  string `json:"content" gorm:"column:content;type:text;comment:内容"`
	VideoID  int64  `json:"video_id" gorm:"column:video_id;type:bigint;comment:视频ID;index"`
	FatherID int64  `json:"father_id" gorm:"column:father_id;type:bigint;comment:父评论ID"`
	UserID   int64  `json:"user_id" gorm:"column:user_id;type:bigint;comment:用户ID"`
	Likes    int    `json:"likes" gorm:"not null;column:likes;type:int;comment:点赞数"`
	//IsAdminLike bool   `json:"is_admin_like" gorm:"not null;column:is_admin_like;type:bool;comment:是否有管理员点赞"`
}

func (Comment) TableName() string {
	return "comments"
}

type CommentLike struct {
	ID int64 `json:"id" gorm:"column:id;primaryKey;type:bigint;type:bigint"`
	TimeModel

	CommentID int64 `json:"comment_id" gorm:"column:comment_id;type:bigint;comment:评论ID;uniqueIndex:idx_comment_user_unique"`
	UserID    int64 `json:"user_id" gorm:"column:user_id;type:bigint;comment:用户ID;uniqueIndex:idx_comment_user_unique"`
}

func (CommentLike) TableName() string {
	return "comment_like"
}
