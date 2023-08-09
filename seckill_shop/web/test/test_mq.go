package main

import (
	"zhiliao_web/rabbitmq"
	"zhiliao_web/utils"
)

func main() {

	qe := rabbitmq.QueueAndExchange{
		QueueName:"test_queue",
		ExchangeName:"test_exchange",
		ExchangeType:"direct",
		RoutingKey:"test_routingKey",

	}
	mq := rabbitmq.NewRabbitMq(qe)

	mq.ConnMq()
	mq.OpenChan()

	defer func() {
		mq.CloseMq()
	}()
	defer func() {
		mq.CloseChan()
	}()

	order_map := map[string]interface{}{
		"uemail":"xxx",
		"pid":6,
	}

	mq.PublishMsg(utils.MapToStr(order_map))
}
