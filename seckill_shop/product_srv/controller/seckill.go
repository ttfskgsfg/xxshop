package controller

import (
	"context"
	"zhiliao_product_srv/proto/seckill"
	"zhiliao_product_srv/models"
	"zhiliao_product_srv/data_source"
	"github.com_bak/pkg/errors"
	"time"
	"zhiliao_product_srv/utils"
	"fmt"
)

type SecKills struct {
}

func (s *SecKills)SecKillList(ctx context.Context, in *zhiliao_product_srv.SecKillsRequest, out *zhiliao_product_srv.SecKillsResponse) error  {

	currentPage := in.CurrentPage
	pagesize := in.Pagesize

	seckills := []models.SecKills{}

	offsetNum := pagesize * (currentPage - 1)
	result := data_source.Db.Limit(pagesize).Offset(offsetNum).Find(&seckills)

	if result.Error != nil {
		out.Code = 500
		out.Msg = "没有查询到数据"
		return errors.New("没有查询到数据")
	}


	var count int32
	seckills_count := []models.SecKills{}
	data_source.Db.Find(&seckills_count).Count(&count)

	seckills_rep := []*zhiliao_product_srv.SecKill{}

	for _,seckill := range seckills {
		seckill_rep := zhiliao_product_srv.SecKill{}
		seckill_rep.Id = int32(seckill.Id)
		seckill_rep.Name = seckill.Name
		seckill_rep.Price = seckill.Price
		seckill_rep.Num = int32(seckill.Num)
		product := models.Products{
			Id:seckill.PId,
		}
		data_source.Db.First(&product)
		seckill_rep.Pid = int32(seckill.PId)
		seckill_rep.Pname = product.Name
		seckill_rep.StartTime = seckill.StartTime.Format("2006-01-02 15:04:05")
		seckill_rep.EndTime = seckill.EndTime.Format("2006-01-02 15:04:05")
		seckill_rep.CreateTime = seckill.CreateTime.Format("2006-01-02 15:04:05")

		seckills_rep = append(seckills_rep,&seckill_rep)
	}
	out.Code = 200
	out.Msg = "查询成功"
	out.Current = currentPage
	out.PageSize = pagesize
	out.Seckills = seckills_rep
	out.Total = count

	return nil
	
}

func (s *SecKills)GetProducts(ctx context.Context, in *zhiliao_product_srv.ProductRequest, out *zhiliao_product_srv.ProductResponse) error  {


	products := []models.Products{}

	result := data_source.Db.Find(&products)

	products_rep := []*zhiliao_product_srv.Product{}



	if result.Error != nil {
		out.Code = 500
		out.Msg = "没有查询到商品信息"
		return errors.New("没有查询到商品信息")
	}

	for _,product := range products {
		product_rep := zhiliao_product_srv.Product{}
		product_rep.Id = int32(product.Id)
		product_rep.Pname = product.Name
		products_rep = append(products_rep, &product_rep)
	}

	out.Code = 200
	out.Msg = "查询成功"
	out.Products = products_rep
	return nil


}

func (s *SecKills)SecKillAdd(ctx context.Context, in *zhiliao_product_srv.SecKill, out *zhiliao_product_srv.SecKillResponse) error  {

	name := in.Name
	price := in.Price
	num := in.Num
	pid := in.Pid
	start_time := in.StartTime
	end_time := in.EndTime

	time_start_time,_ := time.Parse("2006-01-02 15:04:05",start_time)
	time_end_time,_ := time.Parse("2006-01-02 15:04:05",end_time)
	seckill := models.SecKills{
		Name:name,
		Price:price,
		Num:int(num),
		PId:int(pid),
		StartTime:time_start_time,
		EndTime:time_end_time,
		Status:0,
		CreateTime:time.Now(),

	}

	result := data_source.Db.Create(&seckill)

	if result.Error != nil {
		out.Code = 500
		out.Msg = "添加失败"
	}

	out.Code = 200
	out.Msg = "添加成功"
	return  nil

}

func (s *SecKills)SecKillDel(ctx context.Context, in *zhiliao_product_srv.SecKillDelRequest, out *zhiliao_product_srv.SecKillResponse) error  {
	id := in.Id
	// 删除数据库数据操作
	seckill := models.SecKills{
		Id:int(id),
	}
	result := data_source.Db.Delete(&seckill)

	if result.Error != nil {
		out.Code = 500
		out.Msg = "删除失败"
		return nil
	}

	out.Code = 200
	out.Msg = "删除成功"
	return nil

}

func (s *SecKills)SecKillToEdit(ctx context.Context, in *zhiliao_product_srv.SecKillDelRequest, out *zhiliao_product_srv.SecKillToEditResponse) error  {
	id := in.Id

	seckill := models.SecKills{
		Id:int(id),
	}


	reslt := data_source.Db.First(&seckill)

	if reslt.Error != nil {
		out.Code = 500
		out.Msg = "没有查询到数据"
		return errors.New("没有查询到数据")
	}

	seckill_rep := &zhiliao_product_srv.SecKill{}
	seckill_rep.Id = int32(seckill.Id)
	seckill_rep.Name = seckill.Name
	seckill_rep.Price = seckill.Price
	seckill_rep.Num = int32(seckill.Num)
	seckill_rep.Pid = int32(seckill.PId)
	product := models.Products{}
	data_source.Db.Where("id = ?",seckill.PId).Find(&product)
	seckill_rep.Pname = product.Name
	seckill_rep.StartTime = seckill.StartTime.Format("2006-01-02 15:04:05")
	seckill_rep.EndTime = seckill.EndTime.Format("2006-01-02 15:04:05")

	products_no_rep := []*zhiliao_product_srv.Product{}

	products_no := []models.Products{}
	data_source.Db.Where("id != ?",seckill.PId).Find(&products_no)

	for _,product_no := range products_no{
		product_no_rep := zhiliao_product_srv.Product{}
		product_no_rep.Id = int32(product_no.Id)
		product_no_rep.Pname = product_no.Name
		products_no_rep = append(products_no_rep, &product_no_rep)
	}

	out.ProductsNo = products_no_rep

	out.Code = 200
	out.Msg = "成功"
	out.Seckill = seckill_rep
	return nil
}

func (s *SecKills)SecKillDoEdit(ctx context.Context, in *zhiliao_product_srv.SecKill, out *zhiliao_product_srv.SecKillResponse) error  {
	id := in.Id
	name := in.Name
	price := in.Price
	num := in.Num
	Pid := in.Pid
	start_time := in.StartTime
	end_time :=in.EndTime

	time_start_time,_ := time.Parse("2006-01-02 15:04:05",start_time)
	time_end_time,_ := time.Parse("2006-01-02 15:04:05",end_time)
	seckill := models.SecKills{
		Name:name,
		Price:price,
		Num:int(num),
		PId:int(Pid),
		StartTime:time_start_time,
		EndTime:time_end_time,
	}
	result := data_source.Db.Where("id = ?",int(id)).Find(&models.SecKills{}).Update(seckill)

	if result.Error != nil {
		out.Code = 500
		out.Msg = "更新活动失败"
		return nil
	}

	out.Code = 200
	out.Msg = "更新活动成功"
	return nil

}

func (s *SecKills)FrontSecKillList(ctx context.Context, in *zhiliao_product_srv.FrontSecKillRequest, out *zhiliao_product_srv.FrontSecKillResponse) error  {

	/*
	活动显示过滤：
		1.只显示未来一天要做的活动  当前时间+1天  >=  开始时间
		2.一页6条
	 */

	 tomorrow_time := utils.AddHour(24)
	 fmt.Println("==============")
	 fmt.Println(tomorrow_time)

	currentPage := in.CurrentPage
	pagesize := in.Pagesize

	offsetNum := pagesize * (currentPage - 1)
	seckills := []models.SecKills{}
	result := data_source.Db.Where("start_time <= ?",tomorrow_time).Where("status = ?",0).Limit(pagesize).Offset(offsetNum).Find(&seckills)

	if result.Error !=nil {
		out.Code = 500
		out.Msg = "查询不到数据"
		return errors.New("查询不到数据")
	}

	seckills_rep := []*zhiliao_product_srv.SecKill{}

	for _,seckill := range seckills {
		seckill_rep := zhiliao_product_srv.SecKill{}
		seckill_rep.Id = int32(seckill.Id)
		seckill_rep.Name = seckill.Name
		seckill_rep.Price = seckill.Price
		seckill_rep.Num = int32(seckill.Num)
		seckill_rep.Pid = int32(seckill.PId)
		product := models.Products{}
		data_source.Db.Where("id = ?",seckill.PId).Find(&product)
		seckill_rep.Pname = product.Name
		seckill_rep.Pic = product.Pic
		seckill_rep.PPrice = product.Price
		seckill_rep.Pdesc = product.Desc
		seckill_rep.StartTime = seckill.StartTime.Format("2006-01-02 15:04:05")
		seckill_rep.EndTime = seckill.EndTime.Format("2006-01-02 15:04:05")
		seckill_rep.CreateTime = seckill.CreateTime.Format("2006-01-02 15:04:05")

		seckills_rep = append(seckills_rep, &seckill_rep)
	}


	seckills_count := []models.SecKills{}
	var count int32
	data_source.Db.Where("start_time <= ?",tomorrow_time).Where("status = ?",0).Find(&seckills_count).Count(&count)



	out.Code = 200
	out.Msg = "成功"
	out.Current = currentPage
	out.PageSize = pagesize
	out.TotalPage = (count + pagesize -1) / pagesize
	out.SeckillList = seckills_rep

	return nil

}

func (s *SecKills)FrontSecKillDetail(ctx context.Context, in *zhiliao_product_srv.SecKillDelRequest, out *zhiliao_product_srv.FrongSecKillDetailResponse) error  {

	id := in.Id
	fmt.Println(id)

	seckill := models.SecKills{}
	result := data_source.Db.Where("id = ?",id).Find(&seckill)

	if result.Error != nil {
		return errors.New("没有查询到数据")
	}
	product := models.Products{}

	data_source.Db.Where("id = ?",seckill.PId).Find(&product)
	seckill_rep := &zhiliao_product_srv.SecKill{
		Id:int32(seckill.Id),
		Name:seckill.Name,
		Num:int32(seckill.Num),
		Price:seckill.Price,
		Pid:int32(seckill.PId),
		Pname:product.Name,
		Pic:product.Pic,
		PPrice:product.Price,
		Pdesc:product.Desc,
		Unit:product.Unit,
		StartTime:seckill.StartTime.Format("2006-01-02 15:04:05"),
		EndTime:seckill.EndTime.Format("2006-01-02 15:04:05"),
		CreateTime:seckill.CreateTime.Format("2006-01-02 15:04:05"),



	}
	out.Code = 200
	out.Msg = "查询成功"
	out.Seckill = seckill_rep
	return nil
}
