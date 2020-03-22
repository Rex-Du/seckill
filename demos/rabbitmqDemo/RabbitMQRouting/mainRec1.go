// Author : rexdu
// Time : 2020-03-22 14:56
package main

import "seckill/demos/rabbitmqDemo/RabbitMQ"

func main() {
	var rabbitmq *RabbitMQ.RabbitMQ
	rabbitmq = RabbitMQ.NewRabbitMQRouting("exRexdu", "rexdu_one")
	rabbitmq.ReceiveRouting()
}
