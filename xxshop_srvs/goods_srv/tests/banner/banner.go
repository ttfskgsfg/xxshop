package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"xxshop_srvs/goods_srv/proto"
)

var bannerClient proto.GoodsClient
var conn *grpc.ClientConn

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	bannerClient = proto.NewGoodsClient(conn)
}

func TestGetBannerList() {
	rsp, err := bannerClient.BannerList(context.Background(), &empty.Empty{})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	for _, banner := range rsp.Data {
		fmt.Println(banner.Index)
	}
}

func main() {
	Init()
	TestGetBannerList()

	conn.Close()
}
