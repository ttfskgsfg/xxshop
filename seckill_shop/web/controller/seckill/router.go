package seckill

import (
	"github.com/gin-gonic/gin"
	"zhiliao_web/middle_ware"
)

func Router(router *gin.RouterGroup)  {

	// 管理端
	router.GET("/get_seckill_list",middle_ware.JwtTokenValid,GetSeckillList)
	router.GET("/get_products",middle_ware.JwtTokenValid,GetProducts)
	router.POST("/seckill_add",middle_ware.JwtTokenValid,SecKillAdd)
	router.POST("/seckill_del",middle_ware.JwtTokenValid,SecKillDel)
	router.GET("/seckill_to_edit",middle_ware.JwtTokenValid,SeckillToEdit)
	router.POST("/seckill_do_edit",middle_ware.JwtTokenValid,ProductDoEdit)

	// 前端列表
	router.GET("/front/get_seckill_list",GetFrontSeckillList)
	// 前端详情
	router.GET("/front/seckill_detail",middle_ware.JwtTokenFrontValid,SecKillDetail)


	// 秒杀接口
	router.POST("/front/seckill",middle_ware.JwtTokenFrontValid,SecKill)

	// 获取下单结果
	router.GET("/front/get_seckill_result",middle_ware.JwtTokenFrontValid,GetSeckillResult)



}
