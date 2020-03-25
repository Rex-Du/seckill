// Author : rexdu
// Time : 2020-03-26 00:09
package controllers

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
	"seckill/services"
)

type UserController struct {
	Ctx         iris.Context
	UserService services.IUserService
	Session     *sessions.Session
}

func (u *UserController) GetRegister() mvc.View {
	return mvc.View{
		Name: "/user/register.html",
	}
}

//func (u *UserController) PostRegister() mvc.View {
//
//}

func (u *UserController) GetLogin() mvc.View {
	return mvc.View{
		Name: "/user/login.html",
	}
}

//func (u *UserController) PostLogin() mvc.View {
//
//}
