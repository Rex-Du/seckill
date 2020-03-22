// Author : rexdu
// Time : 2020-03-22 01:45 
package main

import (
	"seckill/RabbitMQ"
)

func main() {
	var (
		rabbitmq *RabbitMQ.RabbitMQ
	)
	rabbitmq = RabbitMQ.NewRabbitMQSimple("rexduSimple")
	rabbitmq.ConsumeSimple()
}
