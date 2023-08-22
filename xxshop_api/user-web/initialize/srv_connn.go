package initialize

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"xxshop-api/user-web/global"
	"xxshop-api/user-web/proto"
)

func InitSrvConn() {
	consulInfo := global.Serverconfig.ConsulInfo
	userConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s&", consulInfo.Host, consulInfo.Port, global.Serverconfig.UserSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】")
	}
	userSrvClient := proto.NewUserClient(userConn)
	global.UserSrvClient = userSrvClient
}

func InitSrvConn2() {
	//从注册中心获取到用户服务的信息
	cfg := api.DefaultConfig()
	consulInfo := global.Serverconfig.ConsulInfo
	cfg.Address = fmt.Sprintf("%s:%d", consulInfo.Host, consulInfo.Port)
	//cfg.Address = "127.0.0.1:8500"
	userSrvHost := ""
	userSrvPort := 0
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	//data, err := client.Agent().ServicesWithFilter(fmt.Sprintf("Service == \"%s\"", global.Serverconfig.UserSrvInfo.Name))
	//通过查找注册服务名 动态获取服务中心地址
	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf(`Service == "%s"`, global.Serverconfig.UserSrvInfo.Name))
	if err != nil {
		panic(err)
	}
	for _, value := range data {
		userSrvHost = value.Address
		userSrvPort = value.Port
		break
	}
	if userSrvHost == "" {
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】")
		return
	}
	//拨号连接用户grpc服务器 会有跨域问题
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", userSrvHost, userSrvPort), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 【用户服务失败】",
			"msg", err.Error())
	}
	//1、后续的用户服务下线了 2、改端口 3、该ip

	//2、已经事先创立好了连接，后续不用再进行tcp的三次握手
	//3、一个连接多个groutine共用，性能 - 连接池
	userSrvClient := proto.NewUserClient(userConn)
	global.UserSrvClient = userSrvClient
}
