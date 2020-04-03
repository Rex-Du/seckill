package main

import (
	"github.com/kataras/golog"
	"seckill/common"
	"seckill/rabbitmq"
	"seckill/repositories"
	"seckill/services"
)

func main() {
	golog.SetLevel("debug")
	db, err := common.NewMysqlConn()
	if err != nil {
		golog.Error("数据库连接失败", err)
		return
	}
	orderRepository := repositories.NewOrderManager("order", db)
	orderService := services.NewOrderService(orderRepository)
	productRepository := repositories.NewProductManager("product", db)
	productService := services.NewProductService(productRepository)

	// rabbit消费端
	rabbitmq := rabbitmq.NewRabbitMQSimple("imoocProduct")
	rabbitmq.ConsumeSimple(orderService, productService)
}
