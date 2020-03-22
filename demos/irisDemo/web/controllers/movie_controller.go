// Author : rexdu
// Time : 2020-03-22 23:22
package controllers

import (
	"github.com/kataras/iris/mvc"
	"seckill/demos/irisDemo/repositories"
	"seckill/demos/irisDemo/services"
)

type MovieController struct {
}

func (c *MovieController) Get() mvc.View {
	movieRepository := repositories.NewMovieManager()
	movieService := services.NewMovieServiceManger(movieRepository)
	MovieResult := movieService.ShowMovieName()
	return mvc.View{
		Name: "/movie/index.html",
		Data: MovieResult,
	}
}
