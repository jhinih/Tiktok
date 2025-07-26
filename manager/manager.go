package manager

import (
	"Tiktok/log/zlog"
	"Tiktok/middleware"

	"github.com/gin-gonic/gin"
)

// 主要管理路由组和中间件的注册
// 该包采用管理器模式，集中管理所有路由组和中间件的注册逻辑
// 通过将路由组和中间件的管理抽象成独立的结构体，提高了代码的可维护性和可扩展性

// PathHandler 是一个用于注册路由组的函数类型
// 接收一个路由组指针作为参数，用于在该路由组下注册具体的路由处理函数
type PathHandler func(rg *gin.RouterGroup)

// Middleware 是一个用于生成中间件的函数类型
// 返回一个gin.HandlerFunc类型的中间件函数，用于处理HTTP请求
type Middleware func() gin.HandlerFunc

// RouteManager 管理不同的路由组，按业务功能分组
// 采用组合模式，将不同业务领域的路由组组合在一起
// 每个字段代表一个独立的业务功能模块的路由组
type RouteManager struct {
	LoginRoutes         *gin.RouterGroup // 登录相关的路由组，处理用户认证相关请求
	CommonRoutes        *gin.RouterGroup // 通用功能相关的路由组，处理跨业务模块的通用功能
	UserRoutes          *gin.RouterGroup // 用户相关的路由组，处理用户信息管理相关请求
	MessageRoutes       *gin.RouterGroup // 消息相关的路由组，处理消息通知相关功能
	PostRoutes          *gin.RouterGroup // 帖子相关的路由组，处理帖子内容管理
	FileRoutes          *gin.RouterGroup // 文件相关的路由组，处理文件上传下载
	ContestRoutes       *gin.RouterGroup // 竞赛相关的路由组，处理竞赛功能
	VideosRoutes        *gin.RouterGroup // 视频相关的路由组，处理视频功能
	CommunicationRoutes *gin.RouterGroup // 交流相关的路由组，处理交流功能
}

// NewRouteManager 创建一个新的 RouteManager 实例，包含各业务功能的路由组
// 工厂方法模式，集中创建和初始化所有路由组
// @param router *gin.Engine - Gin框架的路由引擎实例
// @return *RouteManager - 返回初始化好的路由管理器实例
func NewRouteManager(router *gin.Engine) *RouteManager {
	return &RouteManager{
		LoginRoutes:         router.Group("/api/login"),  // 初始化登录路由组
		CommonRoutes:        router.Group("/api/common"), // 通用功能相关的路由组
		UserRoutes:          router.Group("/api/user"),   // 用户相关的路由组
		MessageRoutes:       router.Group("/api/message"),
		PostRoutes:          router.Group("/api/post"),
		FileRoutes:          router.Group("/api/file"),
		ContestRoutes:       router.Group("/api/contest"),
		VideosRoutes:        router.Group("/api/videos"),
		CommunicationRoutes: router.Group("/api/communication"),
	}
}

// RegisterCommonRoutes 通用功能相关的路由组
// @param handler PathHandler - 路由处理函数
func (rm *RouteManager) RegisterCommonRoutes(handler PathHandler) {
	handler(rm.CommonRoutes)
}

// RegisterFileRoutes 注册文件相关的路由处理函数
// @param handler PathHandler - 路由处理函数
func (rm *RouteManager) RegisterFileRoutes(handler PathHandler) {
	handler(rm.FileRoutes)
}

// RegisterLoginRoutes 注册登录相关的路由处理函数
// @param handler PathHandler - 路由处理函数
func (rm *RouteManager) RegisterLoginRoutes(handler PathHandler) {
	handler(rm.LoginRoutes)
}

// RegisterUserRoutes 注册用户相关的路由处理函数
// @param handler PathHandler - 路由处理函数
func (rm *RouteManager) RegisterUserRoutes(handler PathHandler) {
	handler(rm.UserRoutes)
}

// RegisterMessageRoutes 注册消息相关的路由处理函数
// @param handler PathHandler - 路由处理函数
func (rm *RouteManager) RegisterMessageRoutes(handler PathHandler) {
	handler(rm.MessageRoutes)
}

// RegisterPostRoutes 注册帖子相关的路由处理函数
// @param handler PathHandler - 路由处理函数
func (rm *RouteManager) RegisterPostRoutes(handler PathHandler) {
	handler(rm.PostRoutes)
}

// RegisterContestRoutes 注册竞赛相关的路由处理函数
// @param handler PathHandler - 路由处理函数
func (rm *RouteManager) RegisterContestRoutes(handler PathHandler) {
	handler(rm.ContestRoutes)
}

// RegisterVideosRoutes 注册视频相关的路由处理函数
// @param handler PathHandler - 路由处理函数
func (rm *RouteManager) RegisterVideosRoutes(handler PathHandler) {
	handler(rm.VideosRoutes)
}

// RegisterCommunicationRoutes 注册视频相关的路由处理函数
// @param handler PathHandler - 路由处理函数
func (rm *RouteManager) RegisterCommunicationRoutes(handler PathHandler) {
	handler(rm.CommunicationRoutes)
}

// RegisterMiddleware 根据组名为对应的路由组注册中间件
// 策略模式的应用，根据不同的组名选择不同的路由组注册中间件
// @param group string - 路由组名称，可选值: "login"、"common"、"user"等
// @param middleware Middleware - 中间件生成函数
func (rm *RouteManager) RegisterMiddleware(group string, middleware Middleware) {
	switch group {
	case "login":
		rm.LoginRoutes.Use(middleware())
	case "common":
		rm.CommonRoutes.Use(middleware())
	case "user":
		rm.UserRoutes.Use(middleware())
	case "message":
		rm.MessageRoutes.Use(middleware())
	case "diary":
		rm.PostRoutes.Use(middleware())
	case "file":
		rm.FileRoutes.Use(middleware())
	case "contest":
		rm.ContestRoutes.Use(middleware())
	case "videos":
		rm.VideosRoutes.Use(middleware())
	case "communication":
		rm.CommunicationRoutes.Use(middleware())

	}
}

// RequestGlobalMiddleware 注册全局中间件，应用于所有路由
// 装饰器模式的应用，为所有路由添加通用的处理逻辑
// @param r *gin.Engine - Gin框架的路由引擎实例
func RequestGlobalMiddleware(r *gin.Engine) {
	// 添加带调试日志的CORS中间件
	r.Use(func(c *gin.Context) {
		zlog.Infof("Executing CORS middleware - Before handler")
		middleware.Cors()(c)
		c.Next()
		zlog.Infof("Executing CORS middleware - After handler")
	})
}
