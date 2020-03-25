// Author : rexdu
// Time : 2020-03-22 23:40
package main

import (
	"context"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"log"
	"seckill/backend/web/controllers"
	"seckill/common"
	"seckill/repositories"
	"seckill/services"
)

func main() {
	// 1.创建iris实例
	app := iris.New()
	// 2.设置日志等级
	app.Logger().SetLevel("debug")
	// 3.注册模板
	template := iris.HTML("./backend/web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(template)
	// 4.设置模板目标
	app.StaticWeb("/assets", "./backend/web/assets")

	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "访问的页面出错！"))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})
	// 连接数据库
	db, err := common.NewMysqlConn()
	if err != nil {
		log.Println(err)
		return
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	// 5.注册控制器
	productRepository := repositories.NewProductManager("product", db)
	productService := services.NewProductService(productRepository)
	productParty := app.Party("/product")
	product := mvc.New(productParty)
	product.Register(ctx, productService)
	product.Handle(new(controllers.ProductController))

	orderRepository := repositories.NewOrderManager("order", db)
	orderService := services.NewOrderService(orderRepository)
	orderParty := app.Party("/order")
	order := mvc.New(orderParty)
	order.Register(ctx, orderService)
	order.Handle(new(controllers.OrderController))

	// 6.启动服务
	app.Run(
		iris.Addr("localhost:8080"),
		//iris.withoutVersionChecker,
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)

}
