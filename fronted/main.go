// Author : rexdu
// Time : 2020-03-26 00:25
package main

import (
	"context"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"log"
	"seckill/common"
	"seckill/fronted/middleware"
	"seckill/fronted/web/controllers"
	"seckill/rabbitmq"
	"seckill/repositories"
	"seckill/services"
)

func main() {
	// 1.创建iris实例
	app := iris.New()
	// 2.设置日志等级
	app.Logger().SetLevel("debug")
	// 3.注册模板
	template := iris.HTML("./fronted/web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(template)
	// 4.设置模板目标
	app.StaticWeb("/public", "./fronted/web/public")

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

	//sess := sessions.New(
	//	sessions.Config{
	//		Cookie:  "helloworld",
	//		Expires: 60 * time.Minute,
	//	})
	rabbitmq := rabbitmq.NewRabbitMQSimple("imoocProduct")
	// 5.注册控制器

	userRepository := repositories.NewUserRepository("user", db)
	userService := services.NewUserService(userRepository)
	userParty := app.Party("/user")
	user := mvc.New(userParty)
	user.Register(userService, ctx)
	user.Handle(new(controllers.UserController))

	product := repositories.NewProductManager("product", db)
	productService := services.NewProductService(product)
	order := repositories.NewOrderManager("order", db)
	orderService := services.NewOrderService(order)
	proProduct := app.Party("/product")
	proProduct.Use(middleware.AuthConProduct)
	pro := mvc.New(proProduct)
	pro.Register(productService, orderService, rabbitmq)
	pro.Handle(new(controllers.ProductController))

	// 6.启动服务
	app.Run(
		iris.Addr("localhost:8080"),
		//iris.withoutVersionChecker,
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)

}
