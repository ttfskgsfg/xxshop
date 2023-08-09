package main

import (
	"github.com/micro/go-log"
	"github.com/micro/go-micro"

	_ "zhiliao_user_srv/data_source"
	//"zhiliao_user_srv/subscriber"
	//"zhiliao_user_srv/handler"
	//example "zhiliao_user_srv/proto/example"\

	front_user "zhiliao_user_srv/proto/front_user"
	admin_user "zhiliao_user_srv/proto/admin_user"
	"github.com/micro/go-grpc"
	"zhiliao_user_srv/controller"
)

func main() {
	// New Service
	service := grpc.NewService(
		micro.Name("zhiliao.user.srv.zhiliao_user_srv"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	//example.RegisterExampleHandler(service.Server(), new(handler.Example))
	front_user.RegisterFrontUserHandler(service.Server(),new(controller.FrontUser))
	admin_user.RegisterAdminUserHandler(service.Server(),new(controller.AdminUser))




	// Register Struct as Subscriber
	//micro.RegisterSubscriber("zhiliao.user.srv.zhiliao_user_srv", service.Server(), new(subscriber.Example))

	// Register Function as Subscriber
	//micro.RegisterSubscriber("zhiliao.user.srv.zhiliao_user_srv", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
