// Author : rexdu
// Time : 2020-03-31 23:45
package middleware

import "github.com/kataras/iris"

func AuthConProduct(ctx iris.Context) {
	uid := ctx.GetCookie("uid")
	if uid == "" {
		ctx.Application().Logger().Debug("login first")
		ctx.Redirect("/user/login")
		return
	}
	ctx.Application().Logger().Debug("already login")
	ctx.Next()
}
