// Author : rexdu
// Time : 2020-03-22 02:18
package main

import (
	"log"
	"seckill/demos/rabbitmqDemo/RabbitMQ"
	"strconv"
	"time"
)

func main() {
	var rabbitmq *RabbitMQ.RabbitMQ
	rabbitmq = RabbitMQ.NewRabbitMQSub("newProduct")
	for i := 0; i < 100; i++ {
		rabbitmq.PublishPub("订阅模式生产第 " + strconv.Itoa(i) + " 条数据")
		log.Printf("订阅模式生产第 " + strconv.Itoa(i) + " 条数据")
		time.Sleep(1 * time.Second)
	}
}
