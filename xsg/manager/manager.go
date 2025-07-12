package manager

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"tgwp/middleware"
)

//主要管理路由组和中间件的注册

// PathHandler 是一个用于注册路由组的函数类型
type PathHandler func(rg *gin.RouterGroup)

// Middleware 是一个用于生成中间件的函数类型
type Middleware func() gin.HandlerFunc

// RouteManager 管理不同的路由组，按业务功能分组
type RouteManager struct {
	LoginRoutes   *gin.RouterGroup // 登录相关的路由组
	CommonRoutes  *gin.RouterGroup //通用功能相关的路由组
	UserRoutes    *gin.RouterGroup // 用户相关的路由组
	MessageRoutes *gin.RouterGroup
	PostRoutes    *gin.RouterGroup
	FileRoutes    *gin.RouterGroup
	ContestRoutes *gin.RouterGroup
}

// NewRouteManager 创建一个新的 RouteManager 实例，包含各业务功能的路由组
func NewRouteManager(router *gin.Engine) *RouteManager {
	return &RouteManager{
		LoginRoutes:   router.Group("/api/login"),  // 初始化登录路由组
		CommonRoutes:  router.Group("/api/common"), //通用功能相关的路由组
		UserRoutes:    router.Group("/api/user"),   // 用户相关的路由组
		MessageRoutes: router.Group("/api/message"),
		PostRoutes:    router.Group("/api/post"),
		FileRoutes:    router.Group("/api/file"),
		ContestRoutes: router.Group("/api/contest"),
	}
}

// RegisterCommonRoutes 通用功能相关的路由组
func (rm *RouteManager) RegisterCommonRoutes(handler PathHandler) {
	handler(rm.CommonRoutes)
}

func (rm *RouteManager) RegisterFileRoutes(handler PathHandler) {
	handler(rm.FileRoutes)
}

// RegisterLoginRoutes 注册登录相关的路由处理函数
func (rm *RouteManager) RegisterLoginRoutes(handler PathHandler) {
	handler(rm.LoginRoutes)
}

// RegisterUserRoutes 注册用户相关的路由处理函数
func (rm *RouteManager) RegisterUserRoutes(handler PathHandler) {
	handler(rm.UserRoutes)
}

func (rm *RouteManager) RegisterMessageRoutes(handler PathHandler) {
	handler(rm.MessageRoutes)
}

func (rm *RouteManager) RegisterPostRoutes(handler PathHandler) {
	handler(rm.PostRoutes)
}

func (rm *RouteManager) RegisterContestRoutes(handler PathHandler) {
	handler(rm.ContestRoutes)
}

// RegisterMiddleware 根据组名为对应的路由组注册中间件
// group 参数为 "login"、"profile"、"team"或"Common"，分别对应不同的路由组
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
	}
}

// RequestGlobalMiddleware 注册全局中间件，应用于所有路由
func RequestGlobalMiddleware(r *gin.Engine) {
	r.Use(requestid.New())
	r.Use(middleware.AddTraceId())
	r.Use(middleware.Cors())
}
