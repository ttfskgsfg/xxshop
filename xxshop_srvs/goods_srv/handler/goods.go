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

type GoodsServer struct {
	proto.UnimplementedGoodsServer
}

func ModelToResponse(goods model.Goods) proto.GoodsInfoResponse {
	return proto.GoodsInfoResponse{
		Id:              goods.ID,
		CategoryId:      goods.CategoryID,
		Name:            goods.Name,
		GoodsSn:         goods.GoodsSn,
		ClickNum:        goods.ClickNum,
		SoldNum:         goods.SoldNum,
		FavNum:          goods.FavNum,
		MarketPrice:     goods.MarketPrice,
		ShopPrice:       goods.ShopPrice,
		GoodsBrief:      goods.GoodsBrief,
		ShipFree:        goods.ShipFree,
		GoodsFrontImage: goods.GoodsFrontImage,
		IsNew:           goods.IsNew,
		IsHot:           goods.IsHot,
		OnSale:          goods.OnSale,
		DescImages:      goods.DescImages,
		Images:          goods.Images,
		//外键
		Category: &proto.CategoryBriefInfoResponse{
			Id:   goods.Category.ID,
			Name: goods.Category.Name,
		},
		Brand: &proto.BrandInfoResponse{
			Id:   goods.Brands.ID,
			Name: goods.Brands.Name,
			Logo: goods.Brands.Logo,
		},
	}
}

// 商品接口
func (s *GoodsServer) GoodsList(ctx context.Context, req *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	//关键词搜索、查询新品、查询热门商品、通过价格区间筛选、通过商品分类筛选
	goodsListResponse := &proto.GoodsListResponse{}
	//要做到兼容 能判断前端是否传来那些值 通过默认值
	var goods []model.Goods
	//这样写会有问题 即用户同时传递了keywords和ishot，都不会走查询
	//if req.KeyWords != "" {
	//	//搜索
	//	global.DB.Where("name LIKE ?","%"+req.KeyWords+"%").Find(&goods)
	//}
	//if req.IsHot {
	//	global.DB.Where("is_hot=true").Find(&goods)
	//}
	//严重错误 影响了全局变量 global.DB = global.DB.Where("name LIKE ?","%"+req.KeyWords+"%").Find(&goods)
	//正确写法一
	//var queryMap map[string]interface{}  //map这样写会报空的错误
	//queryMap := map[string]interface{}{}
	//if req.KeyWords != "" {
	//	//搜索
	//	queryMap["name"] = "%" + req.KeyWords + "%"
	//}
	//if req.IsHot {
	//	queryMap["is_hot"] = true
	//}
	//if req.IsNew {
	//	queryMap["is_new"] = true
	//}
	//if req.PriceMin > 0 {
	//	//当价格查询比较复杂时 这种方式不灵活
	//	queryMap["shop_price"] = true
	//}
	//正确写法二
	//定义一个局部查询变量  指明要查询的表
	localDB := global.DB.Model(model.Goods{})
	if req.KeyWords != "" {
		//搜索
		localDB = localDB.Where("name LIKE ?", "%"+req.KeyWords+"%")
	}
	if req.IsHot {
		//localDB.Where("is_hot=true") 写法一
		localDB = localDB.Where(model.Goods{IsHot: true}) //写法二
	}
	if req.IsNew {
		localDB = localDB.Where("is_new=true")
	}
	if req.PriceMin > 0 {
		localDB = localDB.Where("shop_price>=?", req.PriceMin)
	}
	if req.PriceMax > 0 {
		localDB = localDB.Where("shop_price<=?", req.PriceMax)
	}
	if req.Brand > 0 {
		localDB = localDB.Where("brand_id=?", req.Brand)
	}
	var subQuery string
	//通过category查询商品  子查询
	if req.TopCategory > 0 {
		var category model.Category
		if result := global.DB.First(&category, req.TopCategory); result.RowsAffected == 0 {
			return nil, status.Errorf(codes.NotFound, "商品分类不存在")
		}
		if category.Level == 1 {
			subQuery = fmt.Sprintf("select id from category WHERE parent_category_id in (select id from category WHERE parent_category_id=%d)", req.TopCategory)
		} else if category.Level == 2 {
			subQuery = fmt.Sprintf("select id from category WHERE parent_category_id=%d", req.TopCategory)
		} else if category.Level == 3 {
			subQuery = fmt.Sprintf("select id from category WHERE id=%d", req.TopCategory)
		}
		localDB = localDB.Where(fmt.Sprintf("category_id in (%s)", subQuery))
	}
	//这么拼凑会有问题 语句子查询要加括号
	//result := localDB.Where("category_id in ?",subQuery).Find(&goods)
	//这样有两个问题 一是没有分页，而是要topcategory有值才走这个语句
	//result := localDB.Where(fmt.Sprintf("category_id in (%s)", subQuery)).Find(&goods)
	//先拿total再分页
	var count int64
	localDB.Count(&count)
	goodsListResponse.Total = int32(count)
	//有外键 要预加载
	result := localDB.Preload("Category").Preload("Brands").Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&goods)

	if result.Error != nil {
		return nil, result.Error
	}
	for _, good := range goods {
		goodsInfoResponse := ModelToResponse(good)
		goodsListResponse.Data = append(goodsListResponse.Data, &goodsInfoResponse)
	}

	return goodsListResponse, nil
}

// 处理通过前端传来的批量获取商品信息请求 用户提交订单时查询
func (s *GoodsServer) BatchGetGoods(ctx context.Context, req *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error) {
	//声明返回对象
	goodsListResponse := &proto.GoodsListResponse{}
	var goods []model.Goods
	//result := global.DB.Where(&goods, req.Id).Find(&goods) //一定要加find语句 只有where是不会执行查询的，只是生成sql
	//调用find和first才会执行
	//Where(&goods, req.Id)中只有是主键才能这么传
	result := global.DB.Where(req.Id).Find(&goods)
	for _, good := range goods {
		goodsInfoResponse := ModelToResponse(good)
		goodsListResponse.Data = append(goodsListResponse.Data, &goodsInfoResponse)
	}
	goodsListResponse.Total = int32(result.RowsAffected)
	return goodsListResponse, nil
}

func (s *GoodsServer) GetGoodsDetail(ctx context.Context, req *proto.GoodsInfoRequest) (*proto.GoodsInfoResponse, error) {
	//定义变量 接住数据库查询来的值
	var goods model.Goods
	//要preload才能查出外键 即联查
	if result := global.DB.Preload("Category").First(&goods, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	goodsInfoResponse := ModelToResponse(goods)
	return &goodsInfoResponse, nil
}

func (s *GoodsServer) CreateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	var brand model.Brands
	if result := global.DB.Preload("Category").First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "品牌不存在")
	}
	//这里没有看到图片文件是如何上传 在微服务中 普通的文件上传方式已经不在使用
	goods := model.Goods{
		Brands:          brand,
		BrandsID:        brand.ID,
		Category:        category,
		CategoryID:      category.ID,
		Name:            req.Name,
		GoodsSn:         req.GoodsSn,
		MarketPrice:     req.MarketPrice,
		ShopPrice:       req.ShopPrice,
		GoodsBrief:      req.GoodsBrief,
		ShipFree:        req.ShipFree,
		Images:          req.Images, //前端传来的字符串
		DescImages:      req.DescImages,
		GoodsFrontImage: req.GoodsFrontImage,
		IsNew:           req.IsNew,
		IsHot:           req.IsHot,
		OnSale:          req.OnSale,
	}
	global.DB.Save(&goods)
	//返回给web层的数据
	return &proto.GoodsInfoResponse{
		Id:         goods.ID,
		CategoryId: goods.CategoryID,
	}, nil
}

func (s *GoodsServer) DeleteGoods(ctx context.Context, req *proto.DeleteGoodsInfo) (*emptypb.Empty, error) {
	if result := global.DB.First(&model.Goods{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	return &emptypb.Empty{}, nil
}

func (s *GoodsServer) UpdateGoods(cxt context.Context, req *proto.CreateGoodsInfo) (*emptypb.Empty, error) {
	var goods model.Goods
	if result := global.DB.First(&goods, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}
	var brand model.Brands
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}
	goods.Brands = brand
	goods.BrandsID = brand.ID
	goods.Category = category
	goods.CategoryID = category.ID
	goods.Name = req.Name
	goods.GoodsSn = req.GoodsSn
	goods.MarketPrice = req.MarketPrice
	goods.ShopPrice = req.ShopPrice
	goods.GoodsBrief = req.GoodsBrief
	goods.ShipFree = req.ShipFree
	goods.Images = req.Images //前端传来的字符串
	goods.DescImages = req.DescImages
	goods.GoodsFrontImage = req.GoodsFrontImage
	goods.IsNew = req.IsNew
	goods.IsHot = req.IsHot
	goods.OnSale = req.OnSale
	global.DB.Save(&goods)
	return &emptypb.Empty{}, nil
}
