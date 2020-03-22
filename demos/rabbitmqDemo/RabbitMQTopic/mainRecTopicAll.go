// Author : rexdu
// Time : 2020-03-22 17:40
package main

import "seckill/demos/rabbitmqDemo/RabbitMQ"

func main() {
	var rabbitmq *RabbitMQ.RabbitMQ
	rabbitmq = RabbitMQ.NewRabbitMQTopic("exTopic", "#")
	rabbitmq.ReceiveTopic()
}
