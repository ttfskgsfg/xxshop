package user

import (
	"github.com/gin-gonic/gin"

	"github.com/micro/go-grpc"
	"zhiliao_user_srv/proto/admin_user"
	"context"
	"net/http"
	"zhiliao_web/utils"
)

func AdminLogin(ctx *gin.Context)  {


	username := ctx.PostForm("username")
	password := ctx.PostForm("password")

	// grpc通信
	service := grpc.NewService()

	admin_user_service := zhiliao_user_srv.NewAdminUserService("zhiliao.user.srv.zhiliao_user_srv",service.Client())
	rep,err := admin_user_service.AdminUserLogin(context.TODO(),&zhiliao_user_srv.AdminUserRequest{
		Username:username,
		Password:password,
	})


	if err != nil {
		ctx.JSON(http.StatusOK,gin.H{
			"code":500,
			"msg":"用户名或密码错误",
		})
	}else {
		admin_token,err1 := utils.GenToken(username,utils.AdminUserExpireDuration,utils.AdminUserSecretKey)

		if err1 != nil {
			ctx.JSON(http.StatusOK,gin.H{
				"code":500,
				"msg":"token错误",
			})
		}
		ctx.JSON(http.StatusOK,gin.H{
			"code":rep.Code,
			"msg":rep.Msg,
			"admin_token":admin_token,
			"user_name":rep.UserName,
		})
	}


}

func FrontUserList(ctx *gin.Context)  {

	currentPage := ctx.DefaultQuery("currentPage","1")
	pageSize := ctx.DefaultQuery("pageSize","10")

	// 和srv通信获取front_users数据
	service := grpc.NewService()
	admin_user_service := zhiliao_user_srv.NewAdminUserService("zhiliao.user.srv.zhiliao_user_srv",service.Client())
	rep,err := admin_user_service.FrontUserList(context.TODO(),&zhiliao_user_srv.FrontUsersRequest{
		CurrentPage:utils.StrToInt(currentPage),
		Pagesize:utils.StrToInt(pageSize),

	})

	if err != nil {
		ctx.JSON(http.StatusOK,gin.H{
			"code":500,
			"msg":"没有查询到数据",
		})
	}


	ctx.JSON(http.StatusOK,gin.H{
		"code":200,
		"msg":"成功",
		"front_users":rep.FrontUsers,
		"total": rep.Total,
		"current_page":rep.Current,
		"page_size":rep.PageSize,
	})


}




