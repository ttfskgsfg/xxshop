package initialize

import (
	"github.com/gin-gonic/gin"
	"xxshop-api/user-web/middlewares"
	"xxshop-api/user-web/router"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	//配置所有router中的中间件
	Router.Use(middlewares.Cors()) //配置跨域
	//实体化路由组
	ApiGroup := Router.Group("/u/v1")
	router.InitUserRouter(ApiGroup)
	router.InitBaseRouter(ApiGroup)
	return Router
}
