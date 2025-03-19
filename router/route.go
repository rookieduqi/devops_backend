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
		serverNodeGroup.GET("/:node_id/view", controller.GetNodeViews)               // 获取节点视图 (Mock 数据)
		serverNodeGroup.POST("/:node_id/view", controller.AddNodeView)               // 添加节点视图 (Mock 数据)
		serverNodeGroup.PUT("/:node_id/view/:view_id", controller.UpdateNodeView)    // 更新节点视图 (Mock 数据)
		serverNodeGroup.DELETE("/:node_id/view/:view_id", controller.DeleteNodeView) // 删除节点视图 (Mock 数据)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}
