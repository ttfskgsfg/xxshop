package handler

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"xxshop_srvs/goods_srv/global"
	"xxshop_srvs/goods_srv/model"
	"xxshop_srvs/goods_srv/proto"
)

// 品牌方面
func (s *GoodsServer) BrandList(ctx context.Context, req *proto.BrandFilterRequest) (*proto.BrandListResponse, error) {
	//先定义返回值
	brandListResponse := proto.BrandListResponse{}
	var brands []model.Brands
	//brands := make([]model.Brands, 0) //也可以这么定义
	//fmt.Println(list)
	//在srv层分页
	//result := global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&brands)
	result := global.DB.Find(&brands)
	if result.Error != nil {
		return nil, result.Error
	}
	//查询表里面的所有数量
	var total int64
	global.DB.Model(&model.Brands{}).Count(&total)
	//fmt.Println(result.RowsAffected) //打印结果数量以及列表
	brandListResponse.Total = int32(total)
	fmt.Println(brandListResponse.Total)
	//返回类型是指针切片 先定义好
	var brandResponses []*proto.BrandInfoResponse
	for _, brand := range brands {
		//放入切片
		//brandResponse := proto.BrandInfoResponse{
		//	Id:brand.ID,
		//	Name: brand.Name,
		//	Logo: brand.Logo,
		//}
		//brandResponses = append(brandResponses, &brandResponse)
		//简洁写法
		brandResponses = append(brandResponses, &proto.BrandInfoResponse{
			Id:   brand.ID,
			Name: brand.Name,
			Logo: brand.Logo,
		})
	}
	brandListResponse.Data = brandResponses
	return &brandListResponse, nil
}

func (s *GoodsServer) CreateBrand(ctx context.Context, req *proto.BrandRequest) (*proto.BrandInfoResponse, error) {
	//新建品牌
	//先查询有没有重复的
	if result := global.DB.Where("name=?", req.Name).First(&model.Brands{}); result.RowsAffected == 1 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌已存在")
	}
	//接收前端传来的值
	brand := model.Brands{
		Name: req.Name,
		Logo: req.Logo}
	//保存
	global.DB.Save(brand)
	return &proto.BrandInfoResponse{Id: brand.ID}, nil
}

func (s *GoodsServer) DeleteBrand(ctx context.Context, req *proto.BrandRequest) (*emptypb.Empty, error) {
	//通过主键id来删除 先实例化brand对象
	if result := global.DB.Delete(&model.Brands{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "品牌不存在")
	}
	return &emptypb.Empty{}, nil
}

func (s *GoodsServer) UpdateBrand(ctx context.Context, req *proto.BrandRequest) (*emptypb.Empty, error) {
	brands := model.Brands{}
	if result := global.DB.First(&brands, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}
	if req.Name != "" {
		brands.Name = req.Name
	}
	if req.Name != "" {
		brands.Logo = req.Logo
	}
	//save会自动判断是update还是create
	global.DB.Save(&brands)
	return &emptypb.Empty{}, nil
}
