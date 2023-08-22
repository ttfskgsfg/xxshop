package handler

import (
	"context"
	"encoding/json"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"xxshop_srvs/goods_srv/global"
	"xxshop_srvs/goods_srv/model"
	"xxshop_srvs/goods_srv/proto"
)

// 商品分类
func (s *GoodsServer) GetAllCategorysList(ctx context.Context, req *emptypb.Empty) (*proto.CategoryListResponse, error) {
	/*
		[
			{
				//一级类目
				"id":xxx,
				"name":"",
				"level":1,
				"is_tab":false,
				"parent":13xxx,
				"sub_category":[  //二级类目
					"id":xxx,
					"name":"",
					"level":1,
					"is_tab":false,
					"sub_category":[]  //三级类目
				]
			}
		]
	*/ //srv层直接返回即可 让web层来实现分类
	var categorys []model.Category
	//要指明预加载哪个字段  反向查询
	//存在问题 只能拿到一级类 一级类拿不到二级类
	//global.DB.Preload("SubCategory").Find(&categorys)
	//此时能拿到2级类目
	//global.DB.Where(&model.Category{Level: 1}).Preload("SubCategory").Find(&categorys)
	//sql实际上执行了三条语句 查询到第三级类目
	global.DB.Where(&model.Category{Level: 1}).Preload("SubCategory.SubCategory").Find(&categorys)
	//for _, category := range categorys {
	//	fmt.Println(category.Name)
	//}  //测试用  //减轻数据库压力可以加缓存
	b, _ := json.Marshal(&categorys)
	//要返回所有种类 可以不用写
	//var bb []*proto.CategoryInfoResponse
	//for _, bbb := range categorys {
	//		bb = append(bb, &proto.CategoryInfoResponse{
	//			Id:             bbb.ID,
	//			Name:           bbb.Name,
	//			ParentCategory: bbb.ParentCategoryID,
	//			Level:          bbb.Level,
	//			IsTab:          bbb.IsTab,
	//		})
	//}
	//业务上只要返回JsonData即可
	return &proto.CategoryListResponse{JsonData: string(b)}, nil
	//return &proto.CategoryListResponse{JsonData: string(b), Data: bb}, nil
}

// // 获取子分类
func (s *GoodsServer) GetSubCategory(ctx context.Context, req *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error) {
	//实例化要返回的对象
	categoryListResponse := proto.SubCategoryListResponse{}
	//实例化查询对象 接住从数据库中取出的值
	var category model.Category
	if result := global.DB.First(&category, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	//构造商品分类本身的信息
	categoryListResponse.Info = &proto.CategoryInfoResponse{
		//从数据库中取出的数
		Id:             category.ID,
		Name:           category.Name,
		Level:          category.Level,
		IsTab:          category.IsTab,
		ParentCategory: category.ParentCategoryID,
	}
	//构造商品子分类
	//实例化对象 接住从数据库中查询到的值，将其成员值赋给返回值里的参数
	var subCategorys []model.Category
	//实例化要返回的值里面的参数  分类的切片
	var subCategoryResponse []*proto.CategoryInfoResponse
	//preloads := "SubCategory"
	//if category.Level == 1 {
	//	preloads = "SubCategory.SubCategory"
	//}
	//global.DB.Where(&model.Category{ParentCategoryID: req.Id}).Preload(preloads).Find(&subCategorys)
	//因为只查询子分类 这样写也可以
	global.DB.Where(&model.Category{ParentCategoryID: req.Id}).Find(&subCategorys)
	for _, subCategory := range subCategorys {
		//类之间赋值
		subCategoryResponse = append(subCategoryResponse, &proto.CategoryInfoResponse{
			Id:             subCategory.ID,
			Name:           subCategory.Name,
			Level:          subCategory.Level,
			IsTab:          subCategory.IsTab,
			ParentCategory: subCategory.ParentCategoryID,
		})
	}
	categoryListResponse.SubCategorys = subCategoryResponse
	//需要分页的场景 就要返回total
	return &categoryListResponse, nil
}

func (s *GoodsServer) CreateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error) {
	category := &model.Category{}
	category.Name = req.Name
	category.Level = req.Level
	if req.Level != 1 {
		//去查询父类目是否存在  也可以让web端去查
		category.ParentCategoryID = req.ParentCategory
	}
	category.IsTab = req.IsTab
	global.DB.Save(&category)
	return &proto.CategoryInfoResponse{Id: int32(category.ID)}, nil
}

func (s *GoodsServer) DeleteCategory(ctx context.Context, req *proto.DeleteCategoryRequest) (*emptypb.Empty, error) {
	if result := global.DB.Delete(&model.Category{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "分类不存在")
	}
	return &emptypb.Empty{}, nil
}

func (s *GoodsServer) UpdateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*emptypb.Empty, error) {
	category := model.Category{}
	if result := global.DB.First(&category, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "分类不存在")
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	//有默认值 要判断下是否有传递值进来
	if req.ParentCategory != 0 {
		category.ParentCategoryID = req.ParentCategory
	}
	if req.Level != 0 {
		category.Level = req.Level
	}
	if req.IsTab {
		category.IsTab = req.IsTab
	}
	global.DB.Save(&category)

	return &emptypb.Empty{}, nil
}
