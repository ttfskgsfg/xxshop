package initialize

import (
	"github.com/gin-gonic/gin"
	"xxshop-api/goods-web/middlewares"
	"xxshop-api/goods-web/router"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	//配置所有router中的中间件
	Router.Use(middlewares.Cors()) //配置跨域
	//实体化路由组
	ApiGroup := Router.Group("/g/v1")
	router.InitGoodsRouter(ApiGroup)
	router.InitCategoryRouter(ApiGroup)
	router.InitBannerRouter(ApiGroup)
	router.InitBrandRouter(ApiGroup)
	return Router
}
