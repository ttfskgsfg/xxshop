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
}

func (SecKills)TableName() string  {
	return "product_seckills"

}
