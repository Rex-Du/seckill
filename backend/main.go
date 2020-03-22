// Author : rexdu
// Time : 2020-03-22 23:40
package main

import "github.com/kataras/iris"

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

	// 5.注册控制器
	// 6.启动服务
	app.Run(
		iris.Addr("localhost:8080"),
		//iris.withoutVersionChecker,
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)

}
