package model

import (
	"Tiktok/global"
	"context"
	"time"
)

func SetUserOnlineInfo(key string, val []byte, timeTTL time.Duration) {
	ctx := context.Background()
	global.Rdb.Set(ctx, key, val, timeTTL)
}
