package types

import (
	"Tiktok/model"
	"mime/multipart"
	"time"
)

// GetVideosRequest 获取视频列表请求
type GetVideosRequest struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	OrderBy  string `json:"order_by"` // random:随机, latest:最新, popular:最热
}

// GetVideosResponse 获取视频列表响应
type GetVideosResponse struct {
	Code     int           `json:"code"`
	Message  string        `json:"message"`
	Data     []model.Video `json:"data"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
	Total    int64         `json:"total"`
	HasMore  bool          `json:"has_more"`
}

// UploadVideoRequest 上传视频请求参数
type UploadVideoRequest struct {
	Title       string                `json:"title" binding:"required,max=100"`
	Description string                `json:"description" binding:"max=500"`
	IsPrivate   bool                  `json:"is_private"`
	CoverFile   *multipart.FileHeader `form:"cover" binding:"-"`
	UploadFile  *multipart.FileHeader `form:"upload" binding:"-"`
	UserID      int64                 `json:"user_id" binding:"-"`
}

// UploadVideoResponse 上传视频响应
type UploadVideoResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    model.Video `json:"data"`
}

// LikeVideoRequest 点赞视频请求参数
type LikeVideoRequest struct {
	VideoID int64 `json:"video_id" binding:"required"`
}

// LikeVideoResponse 点赞视频响应
type LikeVideoResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Likes int `json:"likes"`
	} `json:"data"`
}

//// CommentVideoRequest 评论视频请求
//type CommentVideoRequest struct {
//	VideoID int64 `json:"video_id"`
//	Content string `json:"content"`
//}
//
//// CommentVideoResponse 评论视频响应
//type CommentVideoResponse struct {
//	Code    int    `json:"code"`
//	Message string `json:"message"`
//}

// GetVideosByLastTimeRequest 获取远期视频列表请求
type GetVideosByLastTimeRequest struct {
	LastTime time.Time
}

// GetVideosByLastTimeResponse 获取远期视频列表响应
type GetVideosByLastTimeResponse struct {
	Data []model.Video `json:"data"`
	Time time.Time
}
