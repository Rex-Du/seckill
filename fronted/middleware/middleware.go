// Author : rexdu
// Time : 2020-03-31 23:45
package middleware

import (
	"github.com/kataras/iris"
	"seckill/encrypt"
	"strings"
)

func AuthConProduct(ctx iris.Context) {
	uidString := ctx.GetCookie("sign")
	// 不知道是什么原因cookie中的+变成了空格，所以要先变回来才行
	uidString = strings.Replace(uidString, " ", "+", 1)
	if uidString == "" {
		ctx.Application().Logger().Debug("login first")
		ctx.Redirect("/user/login")
		return
	}
	uidByte, err := encrypt.DePwdCode(uidString)
	if err != nil {
		//if uidString == ""{
		ctx.Application().Logger().Debug(err)
		ctx.Application().Logger().Debug("login first")
		ctx.Redirect("/user/login")
		return
	}
	ctx.Application().Logger().Debug(string(uidByte), " already login")
	ctx.Next()
}
