package router

import (
	"Tiktok/docs"
	"Tiktok/service"
	"github.com/gin-gonic/gin"
	"github.com/sony/gobreaker"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

func Router() *gin.Engine {

	r := gin.Default()

	// 创建限流器，每秒最多 1 个请求
	limiter := rate.NewLimiter(rate.Limit(10), 100)

	r.Use(func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "牢大你被限流了"})
			c.Abort()
			return
		}
		c.Next()
	})

	// 创建熔断器
	breaker := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		MaxRequests: 100,                             // 半开状态下允许的最大请求数
		Interval:    time.Second,                     // 熔断器关闭状态下重置计数器的时间间隔
		Timeout:     time.Duration(10) * time.Second, // 熔断器从打开状态到半开状态的超时时间
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 3 // 连续失败3次触发熔断
		},
	})

	//swagger
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	//静态资源
	r.Static("/asset", "asset/")
	r.StaticFile("/favicon.ico", "asset/images/favicon.ico")
	//	r.StaticFS()
	r.LoadHTMLGlob("views/**/*")

	//首页
	r.GET("/", service.GetIndex, func(c *gin.Context) {
		result, err := breaker.Execute(func() (interface{}, error) {
			// 模拟请求逻辑
			return "Product Service", nil
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": result,
		})
	})
	r.GET("/index", service.GetIndex)
	r.GET("/toRegister", service.ToRegister)
	r.GET("/toChat", service.ToChat)
	r.GET("/chat", service.Chat)
	r.POST("/searchFriends", service.SearchFriends)
	r.POST("/sendEmail", service.SendEmail)

	//用户模块
	r.POST("/user/getUserList", service.GetUserList)
	r.POST("/user/createUser", service.CreateUser)
	r.POST("/user/deleteUser", service.DeleteUser)
	r.POST("/user/updateUser", service.UpdateUser)
	r.POST("/user/findUserByNameAndPwd", service.FindUserByNameAndPwd)
	//r.POST("/user/find", service.FindByID)
	//发送消息
	r.GET("/user/sendMsg", service.SendMsg)
	//发送消息
	r.GET("/user/sendUserMsg", service.SendUserMsg)
	//添加好友
	r.POST("/contact/addfriend", service.AddFriend)
	//上传文件
	r.POST("/attach/upload", service.Upload)
	//创建群
	r.POST("/contact/createCommunity", service.CreateCommunity)
	//群列表
	r.POST("/contact/loadcommunity", service.LoadCommunity)
	//加入群
	r.POST("/contact/joinGroup", service.JoinGroups)
	//心跳续命 不合适  因为Node  所以前端发过来的消息再receProc里面处理
	// r.POST("/user/heartbeat", service.Heartbeat)
	r.POST("/user/redisMsg", service.RedisMsg)
	return r
}
