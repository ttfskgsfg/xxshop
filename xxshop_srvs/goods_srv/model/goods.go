package model

type Category struct {
	BaseModel
	//json序列化的时候自定义名称 即在后面加json的tag
	Name             string    `gorm:"type:varchar(20);not null" json:"name"`
	ParentCategoryID int32     `json:"parent"`
	ParentCategory   *Category `json:"-"` //不想被序列化 要加-
	//定义指向自己的子分类  必须指明外键指向的字段
	//foreignKey表示的是指向另外一张表的哪个键， references指的是外键对应user的哪个键
	SubCategory []*Category `gorm:"foreignKey:ParentCategoryID;references:ID" json:"sub_category"`
	Level       int32       `gorm:"type:int;not null;default:1" json:"level"`
	IsTab       bool        `gorm:"default:false;not null" json:"is_tab"`
}

//上面的SubCategory []*Category类似于
//type Category2 struct {
//	Name             string `gorm:"type:varchar(20);not null" `
//	ParentCategoryID int32
//	ParentCategory   *Category
//	Level            int32 `gorm:"type:int;not null;default:1"`
//	IsTab            bool  `gorm:"default:false;not null"`
//}

// 商品品牌和商品是多对多关系
type Brands struct {
	BaseModel        //组合
	Name      string `gorm:"type:varchar(20);not null"`
	Logo      string `gorm:"type:varchar(200);default:'';not null"`
}

// 自己定义商品分类和品牌表
type GoodsCategoryBrand struct {
	BaseModel
	//对品牌和分类id建立联合唯一的索引 只要索引名称一样，就能建成联合唯一索引
	CategoryID int32    `gorm:"type:int;index:idx_category_brand,unique"`
	Category   Category //分类的外键

	BrandsID int32  `gorm:"type:int;index:idx_category_brand,unique"`
	Brands   Brands //品牌的外键
}

// 重载表名
func (GoodsCategoryBrand) TableName() string {
	return "goodscategorybrand"
}

// 轮播图
type Banner struct {
	BaseModel
	Image string `gorm:"type:varchar(200);not null"`
	Url   string `gorm:"type:varchar(200);not null"`
	Index int32  `gorm:"type:int;default:1;not null"`
}

type Goods struct {
	BaseModel
	//不需要索引是因为商品分类和品牌是由ID来确定，而不是由外键来确定的
	CategoryID int32 `gorm:"type:int;not null"`
	Category   Category
	BrandsID   int32 `gorm:"type:int;not null"`
	Brands     Brands

	OnSale   bool `gorm:"default:false;not null"` //是否上架
	ShipFree bool `gorm:"default:false;not null"` //是否免运费
	IsNew    bool `gorm:"default:false;not null"` //是否是新品
	IsHot    bool `gorm:"default:false;not null"` //是否是热卖品

	Name        string  `gorm:"type:varchar(50);not null"`
	GoodsSn     string  `gorm:"type:varchar(50);not null"`   //商家自己的商品编号
	ClickNum    int32   `gorm:"type:int;default:0;not null"` //商品点击数
	SoldNum     int32   `gorm:"type:int;default:0;not null"` //销量
	FavNum      int32   `gorm:"type:int;default:0;not null"` //收藏数量
	MarketPrice float32 `gorm:"not null"`                    //商品市场价格
	ShopPrice   float32 `gorm:"not null"`                    //商品的优惠价格
	GoodsBrief  string  `gorm:"type:varchar(100);not null"`  //商品简介

	Images          GormList `gorm:"type:varchar(1000);not null"` //商品购物界面图
	DescImages      GormList `gorm:"type:varchar(1000);not null"` //商品详情界面图
	GoodsFrontImage string   `gorm:"type:varchar(200);not null"`  //商品封面图
}
