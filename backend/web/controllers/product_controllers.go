// Author : rexdu
// Time : 2020-03-23 23:34
package controllers

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"seckill/common"
	"seckill/datamodels"
	"seckill/services"
)

type ProductController struct {
	Ctx            iris.Context
	ProductService services.ProductService
}

func (p *ProductController) GetAll() mvc.View {
	products, _ := p.ProductService.GetAllProduct()
	return mvc.View{
		Name: "product/view.html",
		Data: iris.Map{
			"productArray": products,
		},
	}
}

func (p *ProductController) PostUpdate() {
	product := &datamodels.Product{}
	p.Ctx.Request().ParseForm()
	dec := common.NewDecoder(&common.DecoderOptions{TagName: "imooc"})
	if err := dec.Decode(p.Ctx.Request().Form, product); err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	err := p.ProductService.UpdateProduct(product)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	p.Ctx.Redirect("/product/all")
}
