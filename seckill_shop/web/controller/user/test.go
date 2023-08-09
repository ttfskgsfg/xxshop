package user

import (
	"github.com/gin-gonic/gin"
	"time"
	"context"
	"net/http"

	example "zhiliao_user_srv/proto/example"   // user_srv
	//example "zhiliao_product_srv/proto/example"   // product_srv
	//example "zhiliao_seckill_srv/proto/example"   // seckill_srv
	"github.com/micro/go-grpc"
	"fmt"
	"zhiliao_web/utils"
	"strings"
)

func Test(ctx *gin.Context)  {

	ctx.JSON(http.StatusOK,gin.H{
		"code":200,
		"msg":"success",
	})

}


func IndexPost(ctx *gin.Context)  {

	name := ctx.PostForm("name")

	fmt.Println("============")
	fmt.Println(name)
	service := grpc.NewService()

	// user_srv
	exampleClient := example.NewExampleService("zhiliao.user.srv.zhiliao_user_srv", service.Client())
	// product_srv
	//exampleClient := example.NewExampleService("zhiliao.product.srv.zhiliao_product_srv", service.Client())
	// seckill_srv
	//exampleClient := example.NewExampleService("zhiliao.seckill.srv.zhiliao_seckill_srv", service.Client())


	rsp, err := exampleClient.Call(context.TODO(), &example.Request{
		Name: name,
	})

	if err != nil {
		ctx.JSON(http.StatusHTTPVersionNotSupported,gin.H{
			"msg": "发生错误",
			"ref": time.Now().UnixNano(),
		})
	}

	ctx.JSON(http.StatusOK,gin.H{
		"msg": rsp.Msg,
		"ref": time.Now().UnixNano(),
	})


}

func TestToken(ctx *gin.Context)  {

	user_name := "admin"

	token, err := utils.GenToken(user_name,utils.FrontUserExpireDuration,utils.FrontUserSecretKey)

	fmt.Println(err)
	ctx.JSON(http.StatusOK,gin.H{
		"token":token,
	})

}

func TestValidToken(ctx *gin.Context)  {

	header_auth := ctx.Request.Header.Get("Authorization")
	fmt.Println("=============")
	//fmt.Println(header_auth)

	tokenString := strings.Split(header_auth," ")[1]

	claims,err := utils.AuthToken(tokenString,utils.FrontUserSecretKey)
	fmt.Println(err)
	fmt.Println(claims)

	ctx.String(200,"success")

}

