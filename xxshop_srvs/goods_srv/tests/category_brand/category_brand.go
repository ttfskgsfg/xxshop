package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"xxshop_srvs/goods_srv/proto"
)

var brandClient proto.GoodsClient
var conn *grpc.ClientConn

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	brandClient = proto.NewGoodsClient(conn)
}

func TestGetCategoryBrandList() {
	rsp, err := brandClient.CategoryBrandList(context.Background(), &proto.CategoryBrandFilterRequest{})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	fmt.Println(rsp.Data)
	//for _, category := range rsp.Data {
	//	fmt.Println(category.Name)
	//}
}

func TestGetSubCategoryList() {
	rsp, err := brandClient.GetSubCategory(context.Background(), &proto.CategoryListRequest{
		Id: 135487,
	})
	if err != nil {
		panic(err)
	}
	//fmt.Println(rsp.Total)
	//for _, category := range rsp.Data {
	//	fmt.Println(category.Name)
	//}
	fmt.Println(rsp.SubCategorys)
}

func main() {
	Init()
	TestGetCategoryBrandList()
	conn.Close()
}
