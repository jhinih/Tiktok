package api

import (
	"Tiktok/log/zlog"
	"Tiktok/logic"
	"Tiktok/response"
	"Tiktok/types"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"strings"
	"time"
	"unicode/utf8"
)

//	func GetVideosByLastTime(c *gin.Context) {
//		ctx := zlog.GetCtxFromGin(c)
//		req, err := types.BindRequest[types.GetVideosByLastTimeRequest](c)
//		if err != nil {
//			return
//		}
//		zlog.CtxInfof(ctx, "获取视频请求: %v", req)
//		resp, err := logic.NewVideosLogic().GetVideosByLastTime(ctx, req)
//		response.Response(c, resp, err)
//	}
func GetVideos(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindRequest[types.GetVideosRequest](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "获取视频请求: %v", req)
	resp, err := logic.NewVideosLogic().GetVideos(ctx, req)
	response.Response(c, resp, err)
}

// UploadVideo 上传视频
func UploadVideo(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)

	// 记录请求头信息
	zlog.CtxInfof(ctx, "Request headers: %v", c.Request.Header)

	// 检查Content-Type是否为multipart/form-data
	contentType := c.Request.Header.Get("Content-Type")
	if !strings.Contains(contentType, "multipart/form-data") {
		err := errors.New("invalid content type, expected multipart/form-data")
		zlog.CtxErrorf(ctx, "Content-Type error: %v, got: %s", err, contentType)
		response.Response(c, nil, response.ErrResponse(err, response.PARAM_NOT_VALID))
		return
	}

	// 设置100MB最大上传大小和5分钟超时
	const maxUploadSize = 100 << 20 // 100MB
	const timeout = 5 * time.Minute

	// 设置超时上下文
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// 检查文件大小
	if c.Request.ContentLength > maxUploadSize {
		err := errors.New("文件大小超过100MB限制")
		zlog.CtxErrorf(ctx, "文件大小检查失败: %v, 大小: %d", err, c.Request.ContentLength)
		response.Response(c, nil, response.ErrResponse(err, response.PARAM_NOT_VALID))
		return
	}

	// 手动解析multipart表单，限制最大内存100MB
	if err := c.Request.ParseMultipartForm(maxUploadSize); err != nil {
		if err == io.EOF {
			zlog.CtxErrorf(ctx, "ParseMultipartForm EOF error: client may have disconnected, size: %d", c.Request.ContentLength)
		} else {
			zlog.CtxErrorf(ctx, "ParseMultipartForm error: %v, size: %d", err, c.Request.ContentLength)
		}
		response.Response(c, nil, response.ErrResponse(err, response.PARAM_NOT_VALID))
		return
	}

	// 记录表单数据
	zlog.CtxInfof(ctx, "Form data: %v", c.Request.MultipartForm.Value)

	// 获取上传文件
	videoFile, err := c.FormFile("upload")
	if err != nil {
		zlog.CtxErrorf(ctx, "Get video file error: %v", err)
		response.Response(c, nil, response.ErrResponse(err, response.PARAM_NOT_VALID))
		return
	}

	// 记录文件信息
	zlog.CtxInfof(ctx, "Upload file info: %+v", videoFile)

	// 验证文件类型
	if !strings.HasPrefix(videoFile.Header.Get("Content-Type"), "video/") {
		err := errors.New("仅支持视频文件上传")
		zlog.CtxErrorf(ctx, err.Error())
		response.Response(c, nil, response.ErrResponse(err, response.PARAM_NOT_VALID))
		return
	}

	// 获取表单字段
	title := c.PostForm("title")
	description := c.PostForm("description")
	isPrivate := c.PostForm("is_private") == "true"

	// 验证标题长度
	if utf8.RuneCountInString(title) > 100 {
		err := errors.New("标题不能超过100个字符")
		zlog.CtxErrorf(ctx, err.Error())
		response.Response(c, nil, response.ErrResponse(err, response.PARAM_NOT_VALID))
		return
	}

	// 验证描述长度
	if utf8.RuneCountInString(description) > 500 {
		err := errors.New("描述不能超过500个字符")
		zlog.CtxErrorf(ctx, err.Error())
		response.Response(c, nil, response.ErrResponse(err, response.PARAM_NOT_VALID))
		return
	}

	// 构建请求对象
	req := types.UploadVideoRequest{
		Title:       title,
		Description: description,
		IsPrivate:   isPrivate,
		UploadFile:  videoFile,
	}
	if err != nil {
		response.Response(c, nil, response.ErrResponse(err, response.PARAM_NOT_VALID))
		return
	}
	if !strings.HasPrefix(videoFile.Header.Get("Content-Type"), "video/") {
		response.Response(c, nil, response.ErrResponse(errors.New("仅支持视频文件上传"), response.PARAM_NOT_VALID))
		return
	}
	if utf8.RuneCountInString(req.Title) > 100 {
		response.Response(c, nil, response.ErrResponse(errors.New("标题不能超过100个字符"), response.PARAM_NOT_VALID))
		return
	}
	if utf8.RuneCountInString(req.Description) > 500 {
		response.Response(c, nil, response.ErrResponse(errors.New("描述不能超过500个字符"), response.PARAM_NOT_VALID))
		return
	}

	resp, err := logic.NewVideosLogic().UploadVideo(ctx, req)
	response.Response(c, resp, err)
}
