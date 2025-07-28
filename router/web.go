package router

import (
	"Tiktok/configs"
	"Tiktok/global"
	"Tiktok/internal/api"
	"Tiktok/log/zlog"
	"Tiktok/manager"
	"Tiktok/middleware"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RunServer 启动服务器 路由层
func RunServer() {
	r, err := listen()
	if err != nil {
		zlog.Errorf("Listen error: %v", err)
		panic(err.Error())
	}
	r.Run(fmt.Sprintf("%s:%d", configs.Conf.App.Host, configs.Conf.App.Port)) // 启动 Gin 服务器
}

// 自定义ResponseWriter类型
type responseWriter struct {
	gin.ResponseWriter
	headers map[string]string
}

func (w *responseWriter) WriteHeader(code int) {
	for k, v := range w.headers {
		w.ResponseWriter.Header().Set(k, v)
	}
	w.ResponseWriter.WriteHeader(code)
}

func (w *responseWriter) Write(data []byte) (int, error) {
	for k, v := range w.headers {
		w.ResponseWriter.Header().Set(k, v)
	}
	return w.ResponseWriter.Write(data)
}

// listen 配置 Gin 服务器
func listen() (*gin.Engine, error) {
	r := gin.New()

	// 强制在所有响应中添加CORS头
	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}

		// 设置CORS头
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		// 处理OPTIONS请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// 添加Recovery中间件
	r.Use(gin.Recovery())

	// 注册其他全局中间件
	manager.RequestGlobalMiddleware(r)

	// 静态文件路由（支持视频播放）
	r.Static("/uploads", "uploads")

	// 添加视频响应头中间件
	r.Use(func(c *gin.Context) {
		if strings.HasSuffix(c.Request.URL.Path, ".mp4") {
			c.Header("Content-Type", "video/mp4")
			c.Header("Accept-Ranges", "bytes")
		}
	})
	// 创建 RouteManager 实例
	routeManager := manager.NewRouteManager(r)
	// 注册各业务路由组的具体路由
	registerRoutes(routeManager)
	return r, nil
}

// registerRoutes 注册各业务路由的具体处理函数
func registerRoutes(routeManager *manager.RouteManager) {

	// 注册登录相关路由组
	routeManager.RegisterLoginRoutes(func(rg *gin.RouterGroup) {
		rg.POST("/send-code", middleware.Limiter(rate.Every(time.Minute)*4, 4), api.SendCode)
		rg.POST("/register", middleware.Limiter(rate.Every(time.Minute)*4, 4), api.Register)
		rg.POST("/login", middleware.Limiter(rate.Every(time.Minute)*4, 4), api.Login)
		rg.GET("/loginpage", middleware.Limiter(rate.Every(time.Minute)*4, 4), api.Loginpage)
		rg.POST("/refresh-token", middleware.Limiter(rate.Every(time.Second)*4, 8), api.RefreshToken)

		//rg.GET("/test", middleware.Limiter(rate.Every(time.Second)*2, 5), middleware.Authentication(global.ROLE_SUPER_ADMIN), api.TokenTest)
	})

	//注册用户相关路由组
	routeManager.RegisterUserRoutes(func(rg *gin.RouterGroup) {
		rg.GET("/info", middleware.Limiter(rate.Every(time.Second)*20, 40), api.GetUserInfo)
		rg.GET("/my-info", middleware.Limiter(rate.Every(time.Second)*5, 10), middleware.Authentication(global.ROLE_GUEST), api.GetMyUserInfo)
		//获取和修改用户资料
		rg.GET("/profile", middleware.Limiter(rate.Every(time.Second)*10, 20), api.GetProfile)
		rg.POST("/profile", middleware.Limiter(rate.Every(time.Second)*4, 8), middleware.Authentication(global.ROLE_GUEST), api.SetProfile)
		rg.POST("/role", middleware.Limiter(rate.Every(time.Second)*4, 8), middleware.Authentication(global.ROLE_GUEST), api.SetRole)

	})

	// 注册视频相关路由组
	routeManager.RegisterVideosRoutes(func(rg *gin.RouterGroup) {
		// 添加响应头验证
		rg.GET("", func(c *gin.Context) { // 移除斜杠避免重定向
			middleware.Limiter(rate.Every(time.Minute)*20, 40)(c)
			api.GetVideos(c)
		})

		rg.POST("/upload", middleware.Limiter(rate.Every(time.Minute)*20, 40), api.UploadVideo)
		//rg.POST("/getvideosbylasttime", middleware.Limiter(rate.Every(time.Minute)*3, 3), middleware.Authentication(global.ROLE_ADMIN), api.GetVideosByLastTime)
	})
	//注册聊天相关路由组
	routeManager.RegisterCommunicationRoutes(func(rg *gin.RouterGroup) {
		rg.POST("/add-friend", middleware.Limiter(rate.Every(time.Minute)*20, 40), api.AddFriend)
		rg.POST("/search-friends", middleware.Limiter(rate.Every(time.Minute)*20, 40), api.SearchFriend)

		rg.GET("/Send-msg", middleware.Limiter(rate.Every(time.Minute)*20, 40), api.SendMsg)
		rg.GET("/send-user-msg", middleware.Limiter(rate.Every(time.Minute)*20, 40), api.SendUserMsg)
		rg.GET("/send-group-msg", middleware.Limiter(rate.Every(time.Minute)*20, 40), api.SendGroupMsg)

		rg.GET("/get-user-list", middleware.Limiter(rate.Every(time.Minute)*20, 40), api.GetUserList)
		rg.GET("/chat", middleware.Limiter(rate.Every(time.Minute)*20, 40), api.Chats)

		rg.POST("/create-community", middleware.Limiter(rate.Every(time.Minute)*20, 40), api.CreateCommunity)
		rg.POST("/load-community", middleware.Limiter(rate.Every(time.Minute)*20, 40), api.LoadCommunity)
		rg.POST("/join-groups", middleware.Limiter(rate.Every(time.Minute)*20, 40), api.JoinGroups)

		rg.POST("/redis-msg", middleware.Limiter(rate.Every(time.Minute)*20, 40), api.RedisMsg)

		rg.POST("/upload", middleware.Limiter(rate.Every(time.Minute)*20, 40), api.Upload)

	})
}
