// Author : rexdu
// Time : 2020-03-22 17:40
package main

import (
	"seckill/RabbitMQ"
	"strconv"
)

func main() {
	var rabbitOne, rabbitTwo *RabbitMQ.RabbitMQ
	rabbitOne = RabbitMQ.NewRabbitMQTopic("exTopic", "rexdu.topic.one")
	rabbitTwo = RabbitMQ.NewRabbitMQTopic("exTopic", "rexdu.topic.two")
	for i := 0; i < 10; i++ {
		rabbitOne.PublishTopic("hello world one " + strconv.Itoa(i))
		rabbitTwo.PublishTopic("hello world two " + strconv.Itoa(i))
	}
}
