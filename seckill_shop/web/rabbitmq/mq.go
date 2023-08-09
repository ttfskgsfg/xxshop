package rabbitmq

import (
	"github.com/streadway/amqp"
	"fmt"
)

type RabbitMq struct {
	Conn *amqp.Connection
	Ch *amqp.Channel
	QueueName string // 队列名称
	ExchangeName string // 交换机名称
	ExchangeType string // 交换机类型
	RoutingKey string // routingKey
}

type QueueAndExchange struct {
	QueueName string // 队列名称
	ExchangeName string // 交换机名称
	ExchangeType string // 交换机类型
	RoutingKey string // routingKey
}


func (r *RabbitMq)ConnMq()  {
	conn, err := amqp.Dial("amqp://admin:admin@192.168.0.105:5672//myvhost")
	if err != nil {
		fmt.Printf("连接mq出错，错误信息为:%v\n",err)
		return
	}
	r.Conn = conn

}


func (r *RabbitMq)CloseMq()  {
	err := r.Conn.Close()
	if err != nil {
		fmt.Printf("关闭连接出错，错误信息为：%v\n",err)
		return
	}

}

// 开启channel通道
func (r *RabbitMq)OpenChan()  {
	ch,err := r.Conn.Channel()
	if err != nil {
		fmt.Printf("开启channel通道出错，错误信息为：%v\n",err)
		return
	}
	r.Ch = ch
}


// 关闭channnel通道
func (r *RabbitMq)CloseChan()  {
	err := r.Ch.Close()
	if err != nil {
		fmt.Printf("关闭channel通道出错，错误信息为：%v\n",err)
	}
}

// 生产者
func (r *RabbitMq)PublishMsg(body string)  {

	ch := r.Ch
	// 创建队列
	ch.QueueDeclare(r.QueueName,true,false,false,false,nil)

	// 创建交换机
	ch.ExchangeDeclare(r.ExchangeName,r.ExchangeType,true,false,false,false,nil)


	// 队列绑定交换机
	ch.QueueBind(r.QueueName,r.RoutingKey,r.ExchangeName,false,nil)

	// 生产任务
	ch.Publish(r.ExchangeName,r.RoutingKey,false,false,amqp.Publishing{
		ContentType:"text/plain",
		Body:[]byte(body),
		DeliveryMode:amqp.Persistent,

	})

}

// 创建实例
func NewRabbitMq(qe QueueAndExchange) RabbitMq {
	return RabbitMq{
		QueueName:qe.QueueName,
		ExchangeName:qe.ExchangeName,
		ExchangeType:qe.ExchangeType,
		RoutingKey:qe.RoutingKey,
	}
}







