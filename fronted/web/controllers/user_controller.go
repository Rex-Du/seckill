// Author : rexdu
// Time : 2020-03-26 00:09
package controllers

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"seckill/datamodels"
	"seckill/encrypt"
	"seckill/services"
	"seckill/tool"
	"strconv"
)

type UserController struct {
	Ctx     iris.Context
	Service services.IUserService
	//Session *sessions.Session		// session对服务器压力太大，去掉
}

func (u *UserController) GetRegister() mvc.View {
	return mvc.View{
		Name: "/user/register.html",
	}
}

func (c *UserController) PostRegister() {
	var (
		nickName = c.Ctx.FormValue("nickName")
		userName = c.Ctx.FormValue("userName")
		password = c.Ctx.FormValue("password")
	)

	user := &datamodels.User{
		NickName:     nickName,
		UserName:     userName,
		HashPassword: password,
	}
	_, err := c.Service.AddUser(user)
	c.Ctx.Application().Logger().Debug(err)
	if err != nil {
		c.Ctx.Redirect("/user/error")
		return
	}
	c.Ctx.Redirect("/user/login")
	return
}

func (u *UserController) GetLogin() mvc.View {
	return mvc.View{
		Name: "/user/login.html",
	}
}

func (u *UserController) PostLogin() mvc.Response {
	// 获取用户提交的用户名密码
	var (
		userName = u.Ctx.FormValue("userName")
		password = u.Ctx.FormValue("password")
	)
	// 2. 验证账号密码是否正确
	user, isOK := u.Service.IsPwdSuccess(userName, password)

	if !isOK {
		return mvc.Response{
			Path: "/user/login",
		}
	}
	uidByte := []byte(strconv.FormatInt(user.ID, 10))
	uidString, err := encrypt.EnPwdCode(uidByte)
	if err != nil {
		fmt.Println(err)
	}
	tool.GlobalCookie(u.Ctx, "uid", strconv.FormatInt(user.ID, 10), 18000)
	tool.GlobalCookie(u.Ctx, "sign", uidString, 18000)
	//u.Session.Set("userID", strconv.FormatInt(user.ID, 10))

	return mvc.Response{
		Path: "/product/detail",
	}
}
