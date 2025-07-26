package request

import (
	"Tiktok/model"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

type VideosRequest struct {
	DB *gorm.DB
}

func NewVideosRequest(db *gorm.DB) *VideosRequest {
	return &VideosRequest{
		DB: db,
	}
}

//func (r *VideosRequest) GetVideos(source string, before int64, count int) ([]model.Video, error) {
//	var posts []model.Video
//	err := r.DB.Model(&model.Video{}).Where("type = 'Video' AND source = ? AND id < ? ", source, before).Order("id DESC").Limit(count).Find(&model.Video{}).Error
//	return posts, err
//}

// GetVideosByUserId
// 根据作者的id来查询对应数据库数据，并TableVideo返回切片
func (r *VideosRequest) GetVideosByUserID(UserID int64) ([]model.Video, error) {
	//建立结果集接收
	var data []model.Video
	result := r.DB.Where(&model.Video{UserID: UserID}).Find(&data)
	//如果出现问题，返回对应到空，并且返回error
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

// GetVideoByVideoID
// 依据VideoId来获得视频信息
func (r *VideosRequest) GetVideoByVideoID(videoId int64) (model.Video, error) {
	var tableVideo model.Video
	tableVideo.ID = videoId
	//Init()
	result := r.DB.First(&tableVideo)
	if result.Error != nil {
		return tableVideo, result.Error
	}
	return tableVideo, nil

}

// GetVideosByLastTime
// 依据一个时间，来获取这个时间之前的一些视频
func (r *VideosRequest) GetVideosByLastTime(lastTime time.Time) ([]model.Video, error) {
	videos := make([]model.Video, 5)
	result := r.DB.Where("publish_time<?", lastTime).Order("publish_time desc").Limit(5).Find(&videos)
	if result.Error != nil {
		return videos, result.Error
	}
	return videos, nil
}

// GetVideos 获取分页视频列表
func (r *VideosRequest) GetVideos(page, pageSize int, orderBy string) ([]model.Video, error) {
	// 设置默认分页值
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	// 设置默认排序方式
	if orderBy == "" {
		orderBy = "random"
	}

	offset := (page - 1) * pageSize
	var videos []model.Video
	var subQuery *gorm.DB

	// 根据排序方式构建查询
	switch orderBy {
	case "random":
		// 获取所有符合条件的ID
		var allIDs []int64
		err := r.DB.Model(&model.Video{}).
			Select("id").
			Where("is_private = ?", false).
			Find(&allIDs).Error
		if err != nil {
			return nil, err
		}

		// 随机排序ID
		rand.Shuffle(len(allIDs), func(i, j int) {
			allIDs[i], allIDs[j] = allIDs[j], allIDs[i]
		})

		// 分页处理
		start := offset
		if start > len(allIDs) {
			start = len(allIDs)
		}
		end := start + pageSize
		if end > len(allIDs) {
			end = len(allIDs)
		}
		pageIDs := allIDs[start:end]

		// 获取分页视频数据
		err = r.DB.Model(&model.Video{}).
			Where("id IN (?)", pageIDs).
			Find(&videos).Error
		return videos, err
	case "latest":
		// 按创建时间降序
		subQuery = r.DB.Model(&model.Video{}).
			Select("id").
			Where("is_private = ?", false).
			Order("created_at DESC").
			Offset(offset).
			Limit(pageSize)
	case "popular":
		// 按热度(点赞数)降序
		subQuery = r.DB.Model(&model.Video{}).
			Select("id").
			Where("is_private = ?", false).
			Order("likes DESC, created_at DESC").
			Offset(offset).
			Limit(pageSize)
	default:
		// 默认随机排序
		subQuery = r.DB.Model(&model.Video{}).
			Select("id").
			Where("is_private = ?", false).
			Order("RAND()").
			Offset(offset).
			Limit(pageSize)
	}

	// 获取分页视频数据
	err := r.DB.Model(&model.Video{}).
		Where("id IN (?)", subQuery).
		Find(&videos).Error

	return videos, err
}

// GetTotalVideosCount 获取公开视频总数
func (r *VideosRequest) GetTotalVideosCount() (int64, error) {
	var count int64
	err := r.DB.Model(&model.Video{}).
		Where("is_private = ?", false).
		Count(&count).Error
	return count, err
}

// SaveVideo
func (r *VideosRequest) SaveVideo(videoPath, coverPath, title, description string, isPrivate bool, userID int64) (int64, error) {
	video := model.Video{
		Title:       title,
		Description: description,
		URL:         videoPath,
		Cover:       coverPath,
		UserID:      userID,
		IsPrivate:   isPrivate,
		Likes:       0,
		Comments:    0,
		Shares:      0,
		PublishTime: time.Now(),
	}

	if err := r.DB.Create(&video).Error; err != nil {
		return -1, err
	}

	return video.ID, nil
}

//func (r *VideosRequest) LikeVideo(videoID, userID int64) error {
//	// 检查是否已点赞
//	var like model.VideoLike
//	if err := r.DB.Where("video_id = ? AND user_id = ?", videoID, userID).First(&like).Error; err == nil {
//		// 已点赞，执行取消点赞
//		if err := r.DB.Delete(&like).Error; err != nil {
//			return err
//		}
//		// 减少视频点赞数
//		return r.DB.Model(&model.Video{}).Where("id = ?", videoID).
//			Update("likes", gorm.Expr("likes - 1")).Error
//	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
//		return err
//	}
//
//	// 未点赞，执行点赞
//	newLike := model.VideoLike{
//		VideoID: videoID,
//		UserID:  userID,
//	}
//	if err := r.DB.Create(&newLike).Error; err != nil {
//		return err
//	}
//
//	// 增加视频点赞数
//	return r.DB.Model(&model.Video{}).Where("id = ?", videoID).
//		Update("likes", gorm.Expr("likes + 1")).Error
//}
//
//func (r *VideosRequest) GetVideoLikes(videoID string) (int, error) {
//	var video model.Video
//	if err := r.DB.Select("likes").Where("id = ?", videoID).First(&video).Error; err != nil {
//		return 0, err
//	}
//	return video.Likes, nil
//}
