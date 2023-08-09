package middle_ware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"zhiliao_web/utils"
)

func JwtTokenValid(ctx *gin.Context)  {

	auth_header := ctx.Request.Header.Get("Authorization")

	if auth_header == "" {
		ctx.JSON(http.StatusOK,gin.H{
			"code":401,
			"msg":"请携带token",
		})
		ctx.Abort()
		return
	}

	auths := strings.Split(auth_header," ")

	bearer := auths[0]
	token := auths[1]

	if len(token) == 0 || len(bearer) == 0 {
		ctx.JSON(http.StatusOK,gin.H{
			"code":401,
			"msg":"请携带正确格式的token",
		})
		ctx.Abort()
		return
	}

	user, err := utils.AuthToken(token,utils.AdminUserSecretKey)

	if err != nil {
		ctx.JSON(http.StatusOK,gin.H{
			"code":401,
			"msg":"无效的token",
		})
		ctx.Abort()
		return
	}

	ctx.Set("admin_user_name",user.UserName)
	ctx.Next()

}


func JwtTokenFrontValid(ctx *gin.Context)  {

	auth_header := ctx.Request.Header.Get("Authorization")

	if auth_header == "" {
		ctx.JSON(http.StatusOK,gin.H{
			"code":401,
			"msg":"请携带token",
		})
		ctx.Abort()
		return
	}

	auths := strings.Split(auth_header," ")

	bearer := auths[0]
	token := auths[1]

	if len(token) == 0 || len(bearer) == 0 {
		ctx.JSON(http.StatusOK,gin.H{
			"code":401,
			"msg":"请携带正确格式的token",
		})
		ctx.Abort()
		return
	}

	user, err := utils.AuthToken(token,utils.FrontUserSecretKey)

	if err != nil {
		ctx.JSON(http.StatusOK,gin.H{
			"code":401,
			"msg":"无效的token",
		})
		ctx.Abort()
		return
	}

	ctx.Set("front_user_name",user.UserName)
	ctx.Next()

}
