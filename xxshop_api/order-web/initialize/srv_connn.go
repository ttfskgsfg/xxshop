package initialize

import (
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"xxshop-api/order-web/global"
	"xxshop-api/order-web/proto"
)

func InitSrvConn() {
	consulInfo := global.Serverconfig.ConsulInfo
	goodsConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s&", consulInfo.Host, consulInfo.Port, global.Serverconfig.GoodsSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【商品服务失败】")
	}
	//goodsSrvClient := proto.NewGoodsClient(userConn)
	//global.GoodsSrvClient = goodsSrvClient
	//这样写更简洁
	global.GoodsSrvClient = proto.NewGoodsClient(goodsConn)

	orderConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s&", consulInfo.Host, consulInfo.Port, global.Serverconfig.OrderSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【商品服务失败】")
	}
	//goodsSrvClient := proto.NewGoodsClient(userConn)
	//global.GoodsSrvClient = goodsSrvClient
	//这样写更简洁
	global.OrderSrvClient = proto.NewOrderClient(orderConn)

	invConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.Serverconfig.InventorySrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【库存服务失败】")
	}
	global.InventorySrvClient = proto.NewInventoryClient(invConn)
}
