package router

import (
	"github.com/gin-gonic/gin"
	"xxshop-api/goods-web/api/goods"
	"xxshop-api/goods-web/middlewares"
)

// 由调用方传递全局router
func InitGoodsRouter(Router *gin.RouterGroup) {
	GoodsRouter := Router.Group("goods")
	{
		GoodsRouter.GET("", goods.List)                                                   //商品列表
		GoodsRouter.POST("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.New) //该接口需要管理员权限
		GoodsRouter.GET("/:id", goods.Detail)                                             //获取商品详情
		GoodsRouter.DELETE("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Delete)
		GoodsRouter.GET("/:id/stocks", goods.Stocks)
		GoodsRouter.PATCH("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.UpdateStatus)
		GoodsRouter.PUT("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Update)
	}
}
