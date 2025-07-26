package logic

import (
	"Tiktok/global"
	"Tiktok/log/zlog"
	"Tiktok/model"
	"Tiktok/request"
	"Tiktok/response"
	"Tiktok/types"
	"Tiktok/utils"
	"Tiktok/utils/videoUtils"
	"context"
	"errors"
	"strings"
	"time"
)

type VideosLogic struct {
}

func NewVideosLogic() *VideosLogic {
	return &VideosLogic{}
}
func (l *VideosLogic) GetVideos(ctx context.Context, req types.GetVideosRequest) (resp types.GetVideosResponse, err error) {
	defer utils.RecordTime(time.Now())()

	// 记录请求参数
	zlog.CtxInfof(ctx, "获取视频列表请求参数: page=%d, pageSize=%d", req.Page, req.PageSize)

	// 设置默认分页值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 设置默认排序方式
	if req.OrderBy == "" {
		req.OrderBy = "random"
	}
	zlog.CtxInfof(ctx, "使用排序方式: %s", req.OrderBy)

	// 从数据库获取视频列表
	videos, err := request.NewVideosRequest(global.DB).GetVideos(req.Page, req.PageSize, req.OrderBy)
	if err != nil {
		zlog.CtxErrorf(ctx, "获取视频列表失败: %v", err)
		return resp, response.ErrResponse(err, response.PARAM_NOT_VALID)
	}

	// 获取总记录数
	total, err := request.NewVideosRequest(global.DB).GetTotalVideosCount()
	if err != nil {
		zlog.CtxErrorf(ctx, "获取视频总数失败: %v", err)
		return resp, response.ErrResponse(err, response.PARAM_NOT_VALID)
	}

	// 计算是否有更多数据
	hasMore := int64(req.Page*req.PageSize) < total

	// 构建响应数据
	resp.Data = make([]model.Video, 0, len(videos))
	for _, video := range videos {
		resp.Data = append(resp.Data, model.Video{
			ID:          video.ID,
			Title:       video.Title,
			Description: video.Description,
			URL:         video.URL,
			Cover:       video.Cover,
			Likes:       video.Likes,
			Comments:    video.Comments,
			Shares:      video.Shares,
			UserID:      video.UserID,
			CreatedAt:   video.CreatedAt,
			PublishTime: video.PublishTime,
			Type:        video.Type,
			IsPrivate:   video.IsPrivate,
		})
	}

	// 设置分页元数据
	resp.Page = req.Page
	resp.PageSize = req.PageSize
	resp.Total = total
	resp.HasMore = hasMore

	// 设置成功响应
	resp.Code = 200
	resp.Message = "获取视频列表成功"

	// 记录响应数据
	zlog.CtxInfof(ctx, "返回视频列表响应: 当前页%d, 每页%d条, 共%d条, 是否有更多:%v",
		resp.Page, resp.PageSize, resp.Total, resp.HasMore)

	return resp, nil
}
func (l *VideosLogic) UploadVideo(ctx context.Context, req types.UploadVideoRequest) (types.UploadVideoResponse, error) {
	defer utils.RecordTime(time.Now())()

	// 记录上传开始
	zlog.CtxInfof(ctx, "开始处理视频上传, 标题: %s, 用户ID: %d", req.Title, req.UserID)

	// 保存视频文件（不记录文件内容）
	filename, err := videoUtils.SaveUploadedFile(req.UploadFile, "./uploads/videos")
	if err != nil {
		zlog.CtxErrorf(ctx, "保存视频文件失败: %v, 文件名: %s",
			err, req.UploadFile.Filename)
		return types.UploadVideoResponse{}, errors.New("保存视频文件失败")
	}
	// 构建完整视频URL
	videoPath := strings.ReplaceAll(filename, "\\", "/")
	zlog.CtxInfof(ctx, "视频文件保存成功, 路径: %s", videoPath)

	// 保存封面图片
	coverPath := ""
	if req.CoverFile != nil {
		coverPath, err = videoUtils.SaveUploadedFile(req.CoverFile, "./uploads/covers")
		if err != nil {
			zlog.CtxErrorf(ctx, "保存封面图片失败: %v, 文件大小: %d, 文件名: %s",
				err, req.CoverFile.Size, req.CoverFile.Filename)
			return types.UploadVideoResponse{}, errors.New("保存封面图片失败")
		}
		// 将路径中的反斜杠替换为正斜杠，确保URL兼容性
		coverPath = strings.ReplaceAll(coverPath, "\\", "/")
		zlog.CtxInfof(ctx, "封面图片保存成功, 路径: %s", coverPath)
	}

	// 保存到数据库
	videoID, err := request.NewVideosRequest(global.DB).SaveVideo(
		videoPath,
		coverPath,
		req.Title,
		req.Description,
		req.IsPrivate,
		req.UserID,
	)
	if err != nil {
		zlog.CtxErrorf(ctx, "保存视频信息到数据库失败: %v, 视频路径: %s", err, videoPath)
		return types.UploadVideoResponse{}, errors.New("保存视频信息失败")
	}
	zlog.CtxInfof(ctx, "视频信息保存到数据库成功, ID: %d", videoID)

	// 构建响应
	resp := types.UploadVideoResponse{
		Code:    200,
		Message: "视频上传成功",
		Data: model.Video{
			ID:          videoID,
			Title:       req.Title,
			Description: req.Description,
			URL:         videoPath,
			Cover:       coverPath,
			IsPrivate:   req.IsPrivate,
		},
	}

	zlog.CtxInfof(ctx, "视频上传处理完成, 响应: %+v", resp)
	return resp, nil
}
