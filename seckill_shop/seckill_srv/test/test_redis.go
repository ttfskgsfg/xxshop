package main


import (
	"github.com/garyburd/redigo/redis"
	"fmt"
	//"zhiliao_seckill_srv/utils"
)

func main() {

	conn,err := redis.Dial("tcp","192.168.0.105:6379")
	if err != nil {
		fmt.Println("连接出错")
	}

	ret,_ := redis.String(conn.Do("get","1277405413@qq.com"))
	fmt.Println(ret)

	//is_ok,err1 := conn.Do("SET","name","hallen")
	//map_data := map[string]interface{}{
	//	"code" : 500,
	//	"msg" : "下单失败，请重试",
	//}
	//is_ok,err1 := conn.Do("SET","127711@123111.com",utils.MapToStr(map_data))
	//fmt.Println(err1)
	//fmt.Println(is_ok)
	//if is_ok == "OK"{
	//	ret,_ := redis.String(conn.Do("get","name"))
	//	fmt.Println(ret)
	//}

}
