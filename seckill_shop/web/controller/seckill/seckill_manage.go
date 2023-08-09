package seckill

import (
	"github.com/gin-gonic/gin"
	"github.com/micro/go-grpc"
	"zhiliao_product_srv/proto/seckill"
	"context"
	"zhiliao_web/utils"
	"net/http"
	"fmt"
)

func GetSeckillList(ctx *gin.Context)  {
	currentPage := ctx.DefaultQuery("currentPage","1")
	pageSize := ctx.DefaultQuery("pageSize","10")

	// 和srv通信获取products数据
	service := grpc.NewService()

	seckill_service := zhiliao_product_srv.NewSecKillsService("zhiliao.product.srv.zhiliao_product_srv",service.Client())

	rep,err := seckill_service.SecKillList(context.TODO(),&zhiliao_product_srv.SecKillsRequest{
		CurrentPage:utils.StrToInt(currentPage),
		Pagesize:utils.StrToInt(pageSize),

	})

	if err != nil {
		ctx.JSON(http.StatusOK,gin.H{
			"code":500,
			"msg":"没有获取到数据",
		})
	}

	ctx.JSON(http.StatusOK,gin.H{
		"code":rep.Code,
		"msg":rep.Msg,
		"seckills":rep.Seckills,
		"total": rep.Total,
		"current_page":rep.Current,
		"page_size":rep.PageSize,
	})

}

func GetProducts(ctx *gin.Context)  {
	// 返回商品name 和 id

	// 和srv通信获取products数据
	service := grpc.NewService()

	seckill_service := zhiliao_product_srv.NewSecKillsService("zhiliao.product.srv.zhiliao_product_srv",service.Client())

	rep,err := seckill_service.GetProducts(context.TODO(),&zhiliao_product_srv.ProductRequest{})

	if err != nil {
		ctx.JSON(http.StatusOK,gin.H{
			"code":500,
			"msg":"没有可选的商品",
		})
	}

	ctx.JSON(http.StatusOK,gin.H{
		"code":200,
		"msg":"成功",
		"products":rep.Products,
	})


}

func SecKillAdd(ctx *gin.Context)  {
	name := ctx.PostForm("name")
	price := ctx.PostForm("price")
	num := ctx.PostForm("num")
	pid := ctx.PostForm("pid")
	start_time := ctx.PostForm("start_time")
	end_time := ctx.PostForm("end_time")

	fmt.Println("===================")
	fmt.Println(name)
	fmt.Println(price)
	fmt.Println(num)
	fmt.Println(pid)
	fmt.Println(start_time)
	fmt.Println(end_time)

	service := grpc.NewService()

	seckill_service := zhiliao_product_srv.NewSecKillsService("zhiliao.product.srv.zhiliao_product_srv",service.Client())

	rep,_ := seckill_service.SecKillAdd(context.TODO(),&zhiliao_product_srv.SecKill{
		Name:name,
		Price:utils.StrToFloat32(price),
		Num:utils.StrToInt(num),
		Pid:utils.StrToInt(pid),
		StartTime:start_time,
		EndTime:end_time,
	})

	ctx.JSON(http.StatusOK,gin.H{
		"code":rep.Code,
		"msg":rep.Msg,
	})



}

func SecKillDel(ctx *gin.Context)  {
	id := ctx.PostForm("id")
	fmt.Println(id)

	// 和srv通信获取products数据
	service := grpc.NewService()

	seckill_service := zhiliao_product_srv.NewSecKillsService("zhiliao.product.srv.zhiliao_product_srv",service.Client())

	rep,_ := seckill_service.SecKillDel(context.TODO(),&zhiliao_product_srv.SecKillDelRequest{
		Id:utils.StrToInt(id),
	})

	ctx.JSON(http.StatusOK,gin.H{
		"code":rep.Code,
		"msg":rep.Msg,
	})


}


func SeckillToEdit(ctx *gin.Context)  {

	id := ctx.Query("id")
	fmt.Println(id)

	// 和srv通信获取products数据
	service := grpc.NewService()

	seckill_service := zhiliao_product_srv.NewSecKillsService("zhiliao.product.srv.zhiliao_product_srv",service.Client())
	rep,err := seckill_service.SecKillToEdit(context.TODO(),&zhiliao_product_srv.SecKillDelRequest{
		Id:utils.StrToInt(id),
	})

	if err != nil {
		ctx.JSON(http.StatusOK,gin.H{
			"code":500,
			"msg":"没有查询到数据",
		})
	}

	ctx.JSON(http.StatusOK,gin.H{
		"code":rep.Code,
		"msg":rep.Msg,
		"seckill":rep.Seckill,
		"products_no":rep.ProductsNo,
	})

}


func ProductDoEdit(ctx *gin.Context)  {

	id := ctx.PostForm("id")
	name := ctx.PostForm("name")
	price := ctx.PostForm("price")
	num := ctx.PostForm("num")
	pid := ctx.PostForm("pid")
	start_time := ctx.PostForm("start_time")
	end_time := ctx.PostForm("end_time")

	// 和srv通信获取front_users数据
	service := grpc.NewService()

	seckill_service := zhiliao_product_srv.NewSecKillsService("zhiliao.product.srv.zhiliao_product_srv",service.Client())


	rep,_ := seckill_service.SecKillDoEdit(context.TODO(),&zhiliao_product_srv.SecKill{
		Id:utils.StrToInt(id),
		Name:name,
		Num:utils.StrToInt(num),
		Price:utils.StrToFloat32(price),
		StartTime:start_time,
		EndTime:end_time,
		Pid:utils.StrToInt(pid),

	})

	ctx.JSON(http.StatusOK,gin.H{
		"code":rep.Code,
		"msg":rep.Msg,
	})


}

func GetFrontSeckillList(ctx *gin.Context)  {

	currentPage := ctx.DefaultQuery("currentPage","1")
	pageSize := ctx.DefaultQuery("pageSize","8")
	service := grpc.NewService()

	seckill_service := zhiliao_product_srv.NewSecKillsService("zhiliao.product.srv.zhiliao_product_srv",service.Client())
	rep,err := seckill_service.FrontSecKillList(context.TODO(),&zhiliao_product_srv.FrontSecKillRequest{
		CurrentPage:utils.StrToInt(currentPage),
		Pagesize:utils.StrToInt(pageSize),
	})

	if err !=nil {
		ctx.JSON(http.StatusOK,gin.H{
			"code":500,
			"msg":"没有获取到数据",
		})
	}

	for _,seckill := range rep.SeckillList{
		seckill.Pic = utils.Img2Base64(seckill.Pic)
	}

	ctx.JSON(http.StatusOK,gin.H{
		"code":rep.Code,
		"msg":rep.Msg,
		"current":rep.Current,
		"page_size":rep.PageSize,
		"total_page":rep.TotalPage,
		"seckill_list":rep.SeckillList,
	})


}

func SecKillDetail(ctx *gin.Context)  {
	id := ctx.Query("id")

	fmt.Println("+++++++++++++++++")
	fmt.Println(id)
	// grpc 通信
	service := grpc.NewService()

	seckill_service := zhiliao_product_srv.NewSecKillsService("zhiliao.product.srv.zhiliao_product_srv",service.Client())
	rep,err := seckill_service.FrontSecKillDetail(context.TODO(),&zhiliao_product_srv.SecKillDelRequest{
		Id:utils.StrToInt(id),
	})


	if err != nil {
		ctx.JSON(http.StatusOK,gin.H{
			"code":500,
			"msg":"传递参数错误",
		})
	}

	rep.Seckill.Pic = utils.Img2Base64(rep.Seckill.Pic)
	ctx.JSON(http.StatusOK,gin.H{
		"code":rep.Code,
		"msg":rep.Msg,
		"seckill":rep.Seckill,
	})


}

