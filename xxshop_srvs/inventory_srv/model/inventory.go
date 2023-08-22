package model

//仓库
//type Stock struct {
//	BaseModel
//	Name string
//
//}

type Inventory struct {
	BaseModel
	Goods   int32 `gorm:"type:int;index"`
	Stocks  int32 `gorm:"type:int"`
	Version int32 `gorm:"type:int"` //分布式的乐观锁
}

//归还订单
//type InventoryHistory struct {
//	user int32
//	goods int32
//	nums int32
//	order int32
//	status int32 //1. 表示库存是预扣减， 幂等性， 2. 表示已经支付
//}
