// Author : rexdu
// Time : 2020-03-22 18:16
package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"seckill/demos/irisDemo/web/controllers"
)

func main() {
	var (
		app *iris.Application
	)
	app = iris.New()
	app.Logger().SetLevel("debug")
	// 注册模型层
	app.RegisterView(iris.HTML("./web/views", ".html"))
	// 注册控制器
	mvc.New(app.Party("/hello")).Handle(new(controllers.MovieController))
	app.Run(iris.Addr("localhost:8080"))
}
