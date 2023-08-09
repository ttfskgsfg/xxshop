package models

import "time"

type SecKills struct {
	Id int
	Name string
	Price float32 `gorm:"type:decimal(11,2)"`
	Num int
	PId int
	StartTime time.Time
	EndTime time.Time
	// 1表示下架，0表示未下架
	Status int
	CreateTime time.Time
	Orders []Orders `gorm:"ForeignKey:Sid;AssiciationForeignKey:Id"`
}

func (SecKills)TableName() string  {
	return "product_seckills"

}

type Orders struct {
	Id int
	Uemail string

	// 一订单只能属于一个活动，一个活动可以有多个订单
	Sid int
	CreateTime time.Time
	// 0:未支付，1：已支付
	//PayStatus int

}

func (Orders)TableName()string  {
	return "orders"

}

