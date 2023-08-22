package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"xxshop-api/user-web/api"
	"xxshop-api/user-web/middlewares"
)

// 由调用方传递全局router
func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user")
	zap.S().Infof("配置用户相关url")
	{
		//UserRouter.GET("list", api.GetUserList)
		UserRouter.GET("list", middlewares.JWTAuth(), middlewares.IsAdminAuth(), api.GetUserList)
		UserRouter.POST("pwd_login", api.PassWordLogin)
		//UserRouter.POST("register", api.Register)
	}
}
