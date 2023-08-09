package main

import (
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	_ "zhiliao_product_srv/data_source"
	product "zhiliao_product_srv/proto/product"
	seckill "zhiliao_product_srv/proto/seckill"
	"github.com/micro/go-grpc"
	"zhiliao_product_srv/controller"
)

func main() {
	// New Service
	service := grpc.NewService(
		micro.Name("zhiliao.product.srv.zhiliao_product_srv"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()
	product.RegisterProductsHandler(service.Server(),new(controller.Products))
	seckill.RegisterSecKillsHandler(service.Server(),new(controller.SecKills))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
