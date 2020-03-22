// Author : rexdu
// Time : 2020-03-22 14:57
package main

import "seckill/RabbitMQ"

func main() {
	var rabbitmq *RabbitMQ.RabbitMQ
	rabbitmq = RabbitMQ.NewRabbitMQRouting("exRexdu", "rexdu_two")
	rabbitmq.ReceiveRouting()
}
