// Author : rexdu
// Time : 2020-03-22 01:42
package main

import (
	"fmt"
	"seckill/demos/rabbitmqDemo/RabbitMQ"
	"time"
)

func main() {
	var (
		rabbitmq *RabbitMQ.RabbitMQ
	)
	rabbitmq = RabbitMQ.NewRabbitMQSimple("rexduSimple")
	start := time.Now().Unix()
	fmt.Println("开始发送")
	rabbitmq.PublishSimple("hello rexdu")
	fmt.Println("发送成功！")
	fmt.Println("总共耗时：", time.Now().Unix()-start)
}
