package controller

import (
	"context"
	"zhiliao_product_srv/proto/product"
	"fmt"
	"zhiliao_product_srv/models"
	"zhiliao_product_srv/data_source"
	"github.com/pkg/errors"
	"time"
)
type Products struct {
}

func (p *Products) ProductList(ctx context.Context, in *zhiliao_product_srv.ProductsRequest, out *zhiliao_product_srv.ProductsResponse) error  {

	currentPage := in.CurrentPage
	pagesize := in.Pagesize

	products := []models.Products{}

	offsetNum := pagesize * (currentPage - 1)
	result := data_source.Db.Limit(pagesize).Offset(offsetNum).Find(&products)

	if result.Error != nil {
		out.Code = 500
		out.Msg = "没有查询到数据"
		return errors.New("没有查询到数据")
	}


	var count int32
	products_count := []models.Products{}
	data_source.Db.Find(&products_count).Count(&count)

	products_rep := []*zhiliao_product_srv.Product{}

	for _,product := range products {
		product_rep := zhiliao_product_srv.Product{}
		product_rep.Id = int32(product.Id)
		product_rep.Name = product.Name
		product_rep.Price = product.Price
		product_rep.Num = int32(product.Num)
		product_rep.Unit = product.Unit
		product_rep.Pic = product.Pic
		product_rep.Desc = product.Desc
		product_rep.CreateTime = product.CreateTime.Format("2006-01-02 15:04:05")

		products_rep = append(products_rep,&product_rep)
	}
	out.Code = 200
	out.Msg = "查询成功"
	out.Current = currentPage
	out.PageSize = pagesize
	out.Products = products_rep
	out.Total = count

	return nil

}

func (p *Products)ProductAdd(ctx context.Context, in *zhiliao_product_srv.ProductAddRequest, out *zhiliao_product_srv.ProductAddResponse) error  {

	name := in.Name
	price := in.Price
	num := in.Num
	unit := in.Unit
	pic_path := in.Pic
	desc :=in.Desc

	product := models.Products{
		Name:name,
		Price:price,
		Num:int(num),
		Unit:unit,
		Desc:desc,
		Pic:pic_path,
		CreateTime:time.Now(),


	}
	result := data_source.Db.Create(&product)

	if result.Error != nil {
		out.Code = 500
		out.Msg = "添加商品失败"
		return nil
	}

	out.Code = 200
	out.Msg = "添加商品成功"
	return nil


}

func (p *Products) ProductDel(ctx context.Context, in *zhiliao_product_srv.ProductDelRequest, out *zhiliao_product_srv.ProductAddResponse) error  {

	id := in.Id
	fmt.Println(id)
	// 删除数据库数据操作
	product := models.Products{
		Id:int(id),
	}
	result := data_source.Db.Delete(&product)

	if result.Error != nil {
		out.Code = 500
		out.Msg = "删除失败"
		return nil
	}

	out.Code = 200
	out.Msg = "删除成功"
	return nil

}

func (p *Products) ProductToEdit(ctx context.Context, in *zhiliao_product_srv.ProductToEditRequest, out *zhiliao_product_srv.ProductToEditResponse) error  {

	id := in.Id
	fmt.Println(id)

	product := models.Products{
		Id:int(id),
	}




	reslt := data_source.Db.First(&product)

	if reslt.Error != nil {
		out.Code = 500
		out.Msg = "没有查询到数据"
		return errors.New("没有查询到数据")
	}

	product_rep := &zhiliao_product_srv.Product{}
	product_rep.Id = int32(product.Id)
	product_rep.Name = product.Name
	product_rep.Price = product.Price
	product_rep.Num = int32(product.Num)
	product_rep.Unit = product.Unit
	product_rep.Pic = product.Pic
	product_rep.Desc = product.Desc


	out.Code = 200
	out.Msg = "成功"
	out.Product = product_rep
	return nil
}

func (p *Products)ProductDoEdit(ctx context.Context, in *zhiliao_product_srv.ProductEditRequest, out *zhiliao_product_srv.ProductEditResponse) error  {

	id := in.Id
	name := in.Name
	price := in.Price
	num := in.Num
	unit := in.Unit
	pic_path := in.Pic
	desc :=in.Desc

	product := models.Products{}
	if  len(pic_path) < 1 {
		product = models.Products{
			Name:name,
			Price:price,
			Num:int(num),
			Unit:unit,
			Desc:desc,
		}
	}

	product = models.Products{
		Name:name,
		Price:price,
		Num:int(num),
		Unit:unit,
		Desc:desc,
		Pic:pic_path,
	}
	result := data_source.Db.Where("id = ?",int(id)).Find(&models.Products{}).Update(product)

	if result.Error != nil {
		out.Code = 500
		out.Msg = "更新商品失败"
		return nil
	}

	out.Code = 200
	out.Msg = "更新商品成功"
	return nil

}