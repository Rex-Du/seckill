// Author : rexdu
// Time : 2020-03-22 14:56
package main

import (
	"log"
	"seckill/demos/rabbitmqDemo/RabbitMQ"
	"strconv"
	"time"
)

func main() {
	var exOne, exTwo *RabbitMQ.RabbitMQ
	exOne = RabbitMQ.NewRabbitMQRouting("exRexdu", "rexdu_one")
	exTwo = RabbitMQ.NewRabbitMQRouting("exRexdu", "rexdu_two")
	for i := 0; i < 10; i++ {
		exOne.PublishRouting("Hello rexdu one! " + strconv.Itoa(i))
		exTwo.PublishRouting("Hello rexdu two! " + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		log.Println(i)
	}
}
