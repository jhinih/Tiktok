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

	//// 注册通用路由组
	//routeManager.RegisterCommonRoutes(func(rg *gin.RouterGroup) {
	//	rg.POST("/test", middleware.Limiter(rate.Every(time.Hour)*1, 10), api.Template)
	//
	//	rg.GET("/signin-list", middleware.Limiter(rate.Every(time.Minute)*5, 8), middleware.Authentication(global.ROLE_USER), api.SigninList)
	//	rg.POST("/signin", middleware.Limiter(rate.Every(time.Minute)*5, 8), middleware.Authentication(global.ROLE_USER), api.Signin)
	//	rg.POST("/signin-teacher", middleware.Limiter(rate.Every(time.Minute)*5, 8), middleware.Authentication(global.ROLE_USER), api.SigninTeacher)
	//
	//	rg.GET("/auto-list", middleware.Limiter(rate.Every(time.Second)*2, 4), middleware.Authentication(global.ROLE_USER), api.GetAutoList)
	//	rg.POST("/auto-setting", middleware.Limiter(rate.Every(time.Minute)*3, 5), middleware.Authentication(global.ROLE_USER), api.AutoSetting)
	//})

	//// 注册文件上传相关路由组
	//routeManager.RegisterFileRoutes(func(rg *gin.RouterGroup) {
	//	rg.POST("/upload", middleware.Limiter(rate.Every(time.Minute)*3, 5), middleware.Authentication(global.ROLE_GUEST), api.UploadFile)
	//})

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
			middleware.Limiter(rate.Every(time.Minute)*3, 3)(c)
			api.GetVideos(c)
		})

		rg.POST("/upload", middleware.Limiter(rate.Every(time.Minute)*3, 3), api.UploadVideo)
		//rg.POST("/getvideosbylasttime", middleware.Limiter(rate.Every(time.Minute)*3, 3), middleware.Authentication(global.ROLE_ADMIN), api.GetVideosByLastTime)
	})

	routeManager.RegisterCommunicationRoutes(func(rg *gin.RouterGroup) {
		rg.POST("/addfriend", middleware.Limiter(rate.Every(time.Minute)*3, 3), api.AddFriend)
		rg.POST("/searchFriends", middleware.Limiter(rate.Every(time.Minute)*3, 3), api.SearchFriend)
		rg.GET("/sendUserMsg", middleware.Limiter(rate.Every(time.Minute)*3, 3), api.SendUserMsg)
		rg.GET("/SendMsg", middleware.Limiter(rate.Every(time.Minute)*3, 3), api.SendMsg)
		rg.GET("/getUserList", middleware.Limiter(rate.Every(time.Minute)*3, 3), api.GetUserList)
		rg.GET("/chat", middleware.Limiter(rate.Every(time.Minute)*3, 3), api.Chats)
		rg.GET("/createCommunity", middleware.Limiter(rate.Every(time.Minute)*3, 3), api.CreateCommunity)
		rg.POST("/loadCommunity", middleware.Limiter(rate.Every(time.Minute)*3, 3), api.LoadCommunity)
		rg.POST("/redisMsg", middleware.Limiter(rate.Every(time.Minute)*3, 3), api.RedisMsg)
		rg.POST("/joinGroups", middleware.Limiter(rate.Every(time.Minute)*3, 3), api.JoinGroups)

	})
}

////上传文件
//r.POST("/attach/upload", service.Upload)

////心跳续命 不合适  因为Node  所以前端发过来的消息再receProc里面处理
//// r.POST("/user/heartbeat", service.Heartbeat)
//r.POST("/user/redisMsg", service.RedisMsg)
