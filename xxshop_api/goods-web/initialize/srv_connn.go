package initialize

import (
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"xxshop-api/goods-web/global"
	"xxshop-api/goods-web/proto"
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
	//goodsSrvClient := proto.NewGoodsClient(userConn)
	//global.GoodsSrvClient = goodsSrvClient
	//这样写更简洁
	global.GoodsSrvClient = proto.NewGoodsClient(userConn)
}
