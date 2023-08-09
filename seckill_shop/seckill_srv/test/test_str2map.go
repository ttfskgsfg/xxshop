package main

import (
	"zhiliao_seckill_srv/utils"
	"fmt"
)

func main() {

	order_map := map[string]interface{}{
		"uemail":"1277405413@11.com",
		"pid":6,
	}

	str := utils.MapToStr(order_map)

	map_ret := utils.StrToMap(str)
	fmt.Println(map_ret)


}
