package model

import "time"

type ShoppingCart struct { //购物车的表
	BaseModel
	User    int32 `gorm:"type:int;index"` //在购物车列表中我们需要查询当前用户的购物车记录
	Goods   int32 `gorm:"type:int;index"` //加索引：我们需要查询时候， 1. 会影响插入性能 2. 会占用磁盘
	Nums    int32 `gorm:"type:int"`       //商品加了多少件在购物车和中
	Checked bool  //商品是否选中
}

func (ShoppingCart) TableName() string {
	return "shoppingcart"
}

type OrderInfo struct {
	BaseModel

	User    int32  `gorm:"type:int;index"`
	OrderSn string `gorm:"type:varchar(30);index"` //订单号，我们平台自己生成的订单号
	PayType string `gorm:"type:varchar(20) comment 'alipay(支付宝)， wechat(微信)'"`

	//status大家可以考虑使用iota来做
	Status     string     `gorm:"type:varchar(20)  comment 'PAYING(待支付), TRADE_SUCCESS(成功)， TRADE_CLOSED(超时关闭), WAIT_BUYER_PAY(交易创建), TRADE_FINISHED(交易结束)'"`
	TradeNo    string     `gorm:"type:varchar(100) comment '交易号'"` //交易号就是支付宝的订单号 查账
	OrderMount float32    //交易金额
	PayTime    *time.Time `gorm:"type:datetime"`

	Address      string `gorm:"type:varchar(100)"`
	SignerName   string `gorm:"type:varchar(20)"`
	SingerMobile string `gorm:"type:varchar(11)"`
	Post         string `gorm:"type:varchar(20)"`
}

func (OrderInfo) TableName() string {
	return "orderinfo"
}

// 要根据商品查询订单时 如果放入 OrderInfo中查询效率较低
// 查看某一件商品被多少个订单包含 这样方便
type OrderGoods struct {
	BaseModel

	Order int32 `gorm:"type:int;index"`
	Goods int32 `gorm:"type:int;index"`

	//把商品的信息保存下来了 ， 字段冗余， 高并发系统中我们一般都不会遵循三范式  做镜像 记录
	//避免跨服务调用  通过商品名称查订单这样方便些
	GoodsName  string `gorm:"type:varchar(100);index"`
	GoodsImage string `gorm:"type:varchar(200)"`
	GoodsPrice float32
	Nums       int32 `gorm:"type:int"`
}

func (OrderGoods) TableName() string {
	return "ordergoods"
}

//归还订单
//type InventoryHistory struct {
//	user int32
//	goods int32
//	nums int32
//	order int32
//	status int32 //1. 表示库存是预扣减， 幂等性， 2. 表示已经支付
//}
