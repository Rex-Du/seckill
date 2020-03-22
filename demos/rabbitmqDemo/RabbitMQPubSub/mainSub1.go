// Author : rexdu
// Time : 2020-03-22 02:18
package main

import "seckill/demos/rabbitmqDemo/RabbitMQ"

func main() {
	var rabbitmq *RabbitMQ.RabbitMQ
	rabbitmq = RabbitMQ.NewRabbitMQSub("newProduct")
	rabbitmq.ReceiveSub()
}
