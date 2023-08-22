package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xxshop-api/order-web/models"
)

//验证用户是否是管理员

func IsAdminAuth() gin.HandlerFunc {
	//将一些公用代码抽出来然后公用 优势时减少代码量 但是不利于维护 且有连锁反应 --版本管理得做好
	//不单独抽出来  则颗粒性更细
	return func(ctx *gin.Context) {
		claims, _ := ctx.Get("claims")
		currentUser := claims.(*models.CustomClaims)
		if currentUser.AuthorityId != 2 {
			ctx.JSON(http.StatusForbidden, gin.H{
				"msg": "无权限",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
