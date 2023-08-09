package user

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"zhiliao_web/utils"
	"net/http"
	"github.com/micro/go-grpc"
	"zhiliao_user_srv/proto/front_user"
	"context"
)



func SendEmail(ctx *gin.Context)  {

	email := ctx.PostForm("email")
	is_ok := utils.VerifyEmail(email)

	if is_ok{ // 位true则进行grpc通信
		service := grpc.NewService()
		front_user_service := zhiliao_user_srv.NewFrontUserService("zhiliao.user.srv.zhiliao_user_srv",service.Client())

		rep ,_ := front_user_service.FrontUserSendEmail(context.TODO(),&zhiliao_user_srv.FrontUserMailRequest{Email:email})
		fmt.Println(rep)

		ctx.JSON(http.StatusOK,gin.H{
			"code":rep.Code,
			"msg":rep.Msg,
		})

	}else {
		ctx.JSON(http.StatusOK,gin.H{
			"code":500,
			"msg":"邮箱格式不正确",
		})
	}

}

func FrontUserRegister(ctx *gin.Context)  {

	email := ctx.PostForm("email")
	is_ok := utils.VerifyEmail(email)
	if is_ok{
		captche := ctx.PostForm("captche")
		password := ctx.PostForm("password")
		repassword := ctx.PostForm("repassword")

		if password != repassword {
			ctx.JSON(http.StatusOK,gin.H{
				"code":200,
				"msg":"两次密码不一致",
			})
		}else { // 所有都没问题
			// grpc通信
			service := grpc.NewService()
			front_user_service := zhiliao_user_srv.NewFrontUserService("zhiliao.user.srv.zhiliao_user_srv",service.Client())
			rep,err := front_user_service.FrontUserRegister(context.TODO(),&zhiliao_user_srv.FrontUserRequest{Email:email,Code:captche,Password:password,Reassword:repassword})

			if err != nil {
				ctx.JSON(http.StatusOK,gin.H{
					"code":500,
					"msg":"注册失败",
				})
			}else {
				ctx.JSON(http.StatusOK,gin.H{
					"code":rep.Code,
					"msg":rep.Msg,
				})
			}

		}

	}else {
		ctx.JSON(http.StatusOK,gin.H{
			"code":500,
			"msg":"邮箱格式不正确",
		})
	}

}

func FrontUserLogin(ctx *gin.Context)  {

	mail := ctx.PostForm("mail")
	password := ctx.PostForm("password")

	fmt.Println("===============")
	fmt.Println(mail)
	fmt.Println(password)

	is_ok := utils.VerifyEmail(mail)

	if is_ok {
		service := grpc.NewService()

		front_user_service := zhiliao_user_srv.NewFrontUserService("zhiliao.user.srv.zhiliao_user_srv",service.Client())
		rep,err := front_user_service.FrontUserLogin(context.TODO(),&zhiliao_user_srv.FrontUserRequest{
			Email:mail,
			Password:password,
		})


		if err != nil {
			ctx.JSON(http.StatusOK,gin.H{
				"code":500,
				"msg":"用户名或密码错误",
			})
		}else {

			// 生成token
			tokenString,err2 := utils.GenToken(rep.UserName,utils.FrontUserExpireDuration,utils.FrontUserSecretKey)

			if err2 != nil {
				ctx.JSON(http.StatusOK,gin.H{
					"code":rep.Code,
					"msg":rep.Msg,
				})
			}else {
				ctx.JSON(http.StatusOK,gin.H{
					"code":rep.Code,
					"msg":rep.Msg,
					"token":tokenString,
					"username":rep.UserName,
				})
			}

		}

	}else {
		ctx.JSON(http.StatusOK,gin.H{
			"code":500,
			"msg":"邮箱地址不正确",
		})
	}




	
}