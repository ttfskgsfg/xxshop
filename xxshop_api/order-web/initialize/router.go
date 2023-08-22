package initialize

import (
	"github.com/gin-gonic/gin"
	"xxshop-api/order-web/middlewares"
	"xxshop-api/order-web/router"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	//配置所有router中的中间件
	Router.Use(middlewares.Cors()) //配置跨域
	//实体化路由组
	ApiGroup := Router.Group("/o/v1")
	router.InitOrderRouter(ApiGroup)
	router.InitShopCartRouter(ApiGroup)

	return Router
}
