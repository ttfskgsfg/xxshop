package handler

import (
	"context"
	"fmt"
	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"sync"
	"xxshop_srvs/inventory_srv/global"
	"xxshop_srvs/inventory_srv/model"
	"xxshop_srvs/inventory_srv/proto"
)

type InventoryServer struct {
	proto.UnimplementedInventoryServer
}

// 锁一定要全局
var m sync.Mutex

func (*InventoryServer) SetInv(ctx context.Context, req *proto.GoodsInvInfo) (*emptypb.Empty, error) {
	//设置库存 如果要更新库存
	var inv model.Inventory
	//global.DB.First(&inv, req.GoodsId)  //查询主键才能用这种做法 //设置主键可以解决 但是不适合
	global.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inv)
	if inv.Goods == 0 {
		//要加商品信息添加进来
		inv.Goods = req.GoodsId
	}
	inv.Stocks = req.Num
	global.DB.Save(&inv)
	return &emptypb.Empty{}, nil
}

func (*InventoryServer) InvDetail(ctx context.Context, req *proto.GoodsInvInfo) (*proto.GoodsInvInfo, error) {
	var inv model.Inventory
	if result := global.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inv); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "没有库存信息")
	}

	return &proto.GoodsInvInfo{
		GoodsId: inv.Goods,
		Num:     inv.Stocks,
	}, nil
}

func (*InventoryServer) Sell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	//扣减库存 逐一扣减   本地事务(必须全部成功或者失败) 比如要扣减三件商品 结果第二件或者第三件扣减失败
	//先拿到多少件   数据一致性
	//手动事务
	//并发情况下 可能出现超卖
	//m.Lock()  //分布式系统下 如果请求的是不同商品 ，但却用同一把锁 会急剧降低性能
	client := goredislib.NewClient(&goredislib.Options{
		Addr: "127.0.0.1:8089",
	})
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)
	rs := redsync.New(pool)

	tx := global.DB.Begin()
	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory
		//if result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
		//	//参数错误
		//	tx.Rollback() //回滚之前操作
		//	return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		//}
		//for {
		mutex := rs.NewMutex(fmt.Sprintf("goods_%d", goodInfo.GoodsId))
		if err := mutex.Lock(); err != nil {
			return nil, status.Errorf(codes.Internal, "获取redis分布式锁异常")
		}

		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			//参数错误
			tx.Rollback() //回滚之前操作
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}
		if inv.Stocks < goodInfo.Num {
			tx.Rollback()
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}
		//扣减  会出先数据不一致问题 要用分布式锁
		inv.Stocks -= goodInfo.Num
		//update inventory set stocks = stocks-1, version=version+1 where goods=goods and version=version
		//查询inv数据库  要初始化  tx.Model(&inv)这样就是在上边查询后的条件再查询
		//写法有瑕疵 //零值 对于int类型来说 默认值是0 这种会被gorm给忽略掉 必须要加select
		//if result := tx.Model(&model.Inventory{}).Select("Stocks", "Version").Where("goods = ? and version= ?", goodInfo.GoodsId, inv.Version).Updates(model.Inventory{Stocks: inv.Stocks, Version: inv.Version + 1}); result.RowsAffected == 0 {
		//	zap.S().Info("库存扣减失败")
		//} else {
		//	break
		//}
		//使用了事务 只能用开启事务的db来保存
		//global.DB.Save(&inv)
		tx.Save(&inv)
		if ok, err := mutex.Unlock(); !ok || err != nil {
			return nil, status.Errorf(codes.Internal, "释放redis分布式锁异常")
		}

		//}
	}
	tx.Commit() //需要手动提交 修改
	//m.Unlock()
	return &emptypb.Empty{}, nil
}

func (*InventoryServer) Reback(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	//库存归还 1、订单超时归还  2、订单创建失败,归还之前扣减的库存 3、手动归还
	tx := global.DB.Begin()
	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory
		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			//参数错误
			tx.Rollback() //回滚之前操作
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}
		//加库存 会出先数据不一致问题 要用分布式锁
		inv.Stocks += goodInfo.Num
		//使用了事务 只能用开启事务的db来保存
		//global.DB.Save(&inv)
		tx.Save(&inv)
	}
	tx.Commit() //需要手动提交 修改
	return &emptypb.Empty{}, nil
}
