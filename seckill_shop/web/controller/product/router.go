package product

import (
	"github.com/gin-gonic/gin"
	"zhiliao_web/middle_ware"
)

func Router(router *gin.RouterGroup)  {
	router.GET("/get_product_list",middle_ware.JwtTokenValid,GetProductList)
	router.POST("/product_add",middle_ware.JwtTokenValid,ProductAdd)
	router.POST("/product_del",middle_ware.JwtTokenValid,ProductDel)
	router.GET("/to_product_edit",middle_ware.JwtTokenValid,ProductToEdit)
	router.POST("/do_product_edit",middle_ware.JwtTokenValid,ProductDoEdit)

}
