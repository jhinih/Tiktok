package global

import (
	"Tiktok/configs"
	"Tiktok/utils/snowflakeUtils"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	Path     string
	DB       *gorm.DB
	Rdb      *redis.Client
	Config   *configs.Config
	ESClient *elasticsearch.Client

	SnowflakeNode *snowflakeUtils.Node // 默认雪花ID生成节点
)

var (
	OssClient *oss.Client
	OssBucket *oss.Bucket
)
