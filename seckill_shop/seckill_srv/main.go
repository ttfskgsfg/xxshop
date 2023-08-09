package main

import (
	"github.com/micro/go-log"
	"github.com/micro/go-micro"

	seckill "zhiliao_seckill_srv/proto/seckill"
	"github.com/micro/go-grpc"
	"zhiliao_seckill_srv/controller"
	_ "zhiliao_seckill_srv/data_source"
)

func main() {
	// New Service
	service := grpc.NewService(
		micro.Name("zhiliao.seckill.srv.zhiliao_seckill_srv"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler

	seckill.RegisterSecKillHandler(service.Server(),new(controller.SecKill))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
