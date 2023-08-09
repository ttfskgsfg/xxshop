package controller

import (
	"context"

	"zhiliao_seckill_srv/proto/seckill"
	"fmt"
	"zhiliao_seckill_srv/data_source"
	"zhiliao_seckill_srv/models"
	"time"
	"github.com/streadway/amqp"
	"zhiliao_seckill_srv/utils"
	"zhiliao_seckill_srv/redis_lib"
)

type SecKill struct {

}

func (s *SecKill)FrontSecKill(ctx context.Context, in *zhiliao_seckill_srv.SecKillRequest, out *zhiliao_seckill_srv.SecKillResponse) error  {
	id := in.Id

	fmt.Println("===============")
	fmt.Println(id)


	/*
	秒杀逻辑：
		减库存

	限制：
		1.开始时间：还没有开始就不能购买，
			能正常购买的情况：
				当前时间 >= start_time
		2.结束时间：
			能正常购买的情况：
				当前时间 < 结束时间
		3.数量：活动的数量如果为0就能下单

		4.每个用户只能购买一个：依赖于“订单”
	 */
	 seckill := models.SecKills{}
	 now_time := time.Now()
	 result := data_source.Db.Where("id = ?",id).Find(&seckill)

	 if result.Error != nil {
		 out.Code = 500
		 out.Msg = "下单失败，请重试"
		 return nil
	 }

	 fmt.Println("======================now_time")
	fmt.Println(now_time)
	 // 当前时间 >= 开始时间 可以正常抢购
	// 如果查询到数据了，不会报错,能正常下单
	// 如果没有查询到数据，会报错，但是不能正常下单
	 ret_start_time := result.Where("start_time <= ?",now_time).Find(&seckill)

	 if ret_start_time.Error != nil{
	 	out.Code = 500
	 	out.Msg = "抢购还未开始"
		 return nil
	 }

	 // 当前时间 < 结束时间   可以正常抢购

	ret_end_time := ret_start_time.Where("end_time > ?",now_time).Find(&seckill)

	if ret_end_time.Error != nil {
		out.Code = 500
		out.Msg = "抢购已经结束"
		return nil
	}

	// 活动库存必须大于0    可以正常抢购
	ret_num := result.Where("num > 0").Find(&seckill)
	if ret_num.Error != nil {
		out.Code = 500
		out.Msg = "你来晚了，已抢完"
		return nil
	}



	// 获取用户信息
	front_user_email := in.FrontUserEmail


	// 每个用户只能购买一个  购买数量<1  可以正常抢购
	order_rep := models.Orders{}
	ret_uemail := data_source.Db.Where("uemail = ?",front_user_email).Where("sid = ?",id).Find(&order_rep)

	// 如果查询到数据了，不会报错，但是不能正常下单
	// 如果没有查询到数据，会报错，但是能正常下单
	if ret_uemail.Error == nil {
		out.Code = 500
		out.Msg = "下单失败，每个用户只能购买一次"
		return nil

	}

	// 正常的
	ret := result.Update("num",seckill.Num-1)

	// 生成订单：购买成功才往订单表中插入数据

	order := models.Orders{
		Uemail:front_user_email,
		Sid:int(id),
		CreateTime:time.Now(),
	}
	ret_order := data_source.Db.Create(&order)

	if ret_order.Error != nil {
		out.Code = 500
		out.Msg = "下单失败，请重试"
		return nil

	}

	if ret.Error != nil {
		out.Code = 500
		out.Msg = "下单失败，请重试"
		return nil
	}

	out.Code = 200
	out.Msg = "下单成功"
	return nil

}


// 从队列中读取任务并消费

func init()  {

	conn, err := amqp.Dial("amqp://admin:admin@192.168.0.105:5672//myvhost")
	fmt.Println(err)
	defer conn.Close()


	ch,err_ch := conn.Channel()
	fmt.Println(err_ch)
	defer ch.Close()

	ch.Qos(1,0,false)


	deliveries,err := ch.Consume("zhiliao_web.order_queue","order_consumer",false,false,false,false,nil)
	//deliveries,err := ch.Consume("first_queue","first_consumer",false,false,false,false,nil)
	if err != nil {
		fmt.Println(err)
	}

	for delivery := range deliveries {
		//fmt.Println(delivery.ContentType)
		//fmt.Println(string(delivery.Body))
		//delivery.Ack(true)
		fmt.Println("接收到任务")
		go Orderapply(delivery)
	}
}

func Orderapply(delivery amqp.Delivery)  {
	body := string(delivery.Body)
	fmt.Println("============")
	fmt.Println(body)
	fmt.Printf("%T\n",body)

	request_data := utils.StrToMap(body)


	id := request_data["pid"]
	front_user_email := request_data["uemail"].(string)

	order_mq := models.Orders{}
	var count_mq int
	data_source.Db.Where("uemail = ?",front_user_email).Find(&order_mq).Count(&count_mq)

	if count_mq > 0 {
		delivery.Ack(true)
		return
	}

	seckill := models.SecKills{}
	now_time := time.Now()
	result := data_source.Db.Where("id = ?",id).Find(&seckill)


	if result.Error != nil {
		delivery.Ack(true)
		map_data := map[string]interface{}{
			"code" : 500,
			"msg" : "下单失败，请重试",
		}
		is_ok,_ := redis_lib.Conn.Do("GET",front_user_email)
		if is_ok == "OK"{
			return
		}
		redis_lib.Conn.Do("SET",front_user_email,utils.MapToStr(map_data))
		return
	}

	fmt.Println("======================now_time")
	fmt.Println(now_time)
	// 当前时间 >= 开始时间 可以正常抢购
	// 如果查询到数据了，不会报错,能正常下单
	// 如果没有查询到数据，会报错，但是不能正常下单
	ret_start_time := result.Where("start_time <= ?",now_time).Find(&seckill)

	if ret_start_time.Error != nil{
		delivery.Ack(true)
		map_data := map[string]interface{}{
			"code" : 500,
			"msg" : "抢购还未开始",
		}
		is_ok,_ := redis_lib.Conn.Do("GET",front_user_email)
		if is_ok == "OK"{
			return
		}
		redis_lib.Conn.Do("SET",front_user_email,utils.MapToStr(map_data))
		return
	}

	// 当前时间 < 结束时间   可以正常抢购

	ret_end_time := ret_start_time.Where("end_time > ?",now_time).Find(&seckill)

	if ret_end_time.Error != nil {
		delivery.Ack(true)
		map_data := map[string]interface{}{
			"code" : 500,
			"msg" : "抢购已经结束",
		}
		is_ok,_ := redis_lib.Conn.Do("GET",front_user_email)
		if is_ok == "OK"{
			return
		}
		redis_lib.Conn.Do("SET",front_user_email,utils.MapToStr(map_data))
		return
	}

	// 活动库存必须大于0    可以正常抢购
	ret_num := result.Where("num > 0").Find(&seckill)
	if ret_num.Error != nil {
		delivery.Ack(true)
		map_data := map[string]interface{}{
			"code" : 500,
			"msg" : "你来晚了，已抢完",
		}
		is_ok,_ := redis_lib.Conn.Do("GET",front_user_email)
		if is_ok == "OK"{
			return
		}
		redis_lib.Conn.Do("SET",front_user_email,utils.MapToStr(map_data))
		return
	}

	// 获取用户信息

	// 每个用户只能购买一个  购买数量<1  可以正常抢购
	order_rep := models.Orders{}
	ret_uemail := data_source.Db.Where("uemail = ?",front_user_email).Where("sid = ?",id).Find(&order_rep)

	// 如果查询到数据了，不会报错，但是不能正常下单
	// 如果没有查询到数据，会报错，但是能正常下单
	if ret_uemail.Error == nil {
		delivery.Ack(true)
		map_data := map[string]interface{}{
			"code" : 500,
			"msg" : "下单失败，每个用户只能购买一次",
		}
		is_ok,_ := redis_lib.Conn.Do("GET",front_user_email)
		if is_ok == "OK"{
			return
		}
		redis_lib.Conn.Do("SET",front_user_email,utils.MapToStr(map_data))


	}

	// 正常的
	ret := result.Update("num",seckill.Num-1)
	// 生成订单：购买成功才往订单表中插入数据

	fmt.Println("=========+++++++++++++++")
	fmt.Println(id.(string))
	order := models.Orders{
		Uemail:front_user_email,
		Sid:utils.StrToInt(id.(string)),
		CreateTime:time.Now(),
	}
	ret_order := data_source.Db.Create(&order)

	if ret_order.Error != nil {
		delivery.Ack(true)
		map_data := map[string]interface{}{
			"code" : 500,
			"msg" : "下单失败，请重试",
		}
		is_ok,_ := redis_lib.Conn.Do("GET",front_user_email)
		if is_ok == "OK"{
			return
		}
		redis_lib.Conn.Do("SET",front_user_email,utils.MapToStr(map_data))
		return

	}

	if ret.Error != nil {
		delivery.Ack(true)
		map_data := map[string]interface{}{
			"code" : 500,
			"msg" : "下单失败，请重试",
		}
		is_ok,_ := redis_lib.Conn.Do("GET",front_user_email)
		if is_ok == "OK"{
			return
		}
		redis_lib.Conn.Do("SET",front_user_email,utils.MapToStr(map_data))
		return
	}

	// 成功的
	delivery.Ack(true)
	map_data := map[string]interface{}{
		"code" : 200,
		"msg" : "下单成功",
	}
	is_ok,_ := redis_lib.Conn.Do("GET",front_user_email)
	if is_ok == "OK"{
		return
	}
	redis_lib.Conn.Do("SET",front_user_email,utils.MapToStr(map_data))
	return



}




