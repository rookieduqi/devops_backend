package router

import (
	"bluebell/controller"
	"bluebell/logger"
	"bluebell/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) // gin设置成发布模式
	}
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	// 注册
	r.POST("/signup", controller.SignUpHandler)
	// 登录
	r.POST("/login", controller.LoginHandler)

	r.GET("/ping", middlewares.JWTAuthMiddleware(), func(c *gin.Context) {
		// 如果是登录的用户,判断请求头中是否有 有效的JWT  ？
		c.JSON(http.StatusOK, gin.H{
			"msg": "ok",
		})
	})

	// 节点管理
	//r.GET("/server/nodes", controller.GetNameServerNodes)
	serverNodeGroup := r.Group("/server/node")
	{
		serverNodeGroup.POST("", controller.AddServerNode)          // 新增
		serverNodeGroup.GET("", controller.GetServerNodes)          // 获取
		serverNodeGroup.PUT("", controller.UpdateServerNode)        // 更新
		serverNodeGroup.DELETE("/:id", controller.DeleteServerNode) // 删除
	}

	serverNodeGroup = r.Group("/server/node_view")
	{
		serverNodeGroup.POST("/get/view", controller.GetNodeViews)
	}

	serverNodeGroup = r.Group("/server/view")
	{
		serverNodeGroup.POST("/get", controller.GetNodeViewsT)
	}

	serverNodeGroup = r.Group("/server/view_jobs")
	{
		serverNodeGroup.POST("/get/job", controller.GetNodeJobsT)
		serverNodeGroup.POST("/start/job", controller.StartNodeJobsT)
		serverNodeGroup.POST("/stop/job", controller.StopNodeJobsT)
	}

	serverNodeGroup = r.Group("/server/view_console")
	{
		serverNodeGroup.POST("/get", controller.GetNodeConsole)
		serverNodeGroup.POST("/pipeline/overview", controller.GetConsolePipeOverview)
		serverNodeGroup.POST("/pipeline/console", controller.GetConsolePipeConsole)
		serverNodeGroup.POST("/build/previous", controller.ConsoleBuildPrevious)
		serverNodeGroup.DELETE("/build/delete", controller.ConsoleBuildDelete)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "接口不存在",
		})
	})
	return r
}
