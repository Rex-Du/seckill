// Author : rexdu
// Time : 2020-03-31 22:53
package controllers

import (
	"encoding/json"
	"github.com/kataras/golog"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
	"seckill/datamodels"
	"seckill/rabbitmq"
	"seckill/services"
	"strconv"
)

type ProductController struct {
	Ctx            iris.Context
	ProductService services.IProductService
	OrderService   services.IOrderService
	RabbitMQ       *rabbitmq.RabbitMQ
	Session        *sessions.Session
}

func (p *ProductController) GetDetail() mvc.View {
	//id := p.Ctx.URLParam("productid")
	product, _ := p.ProductService.GetProductByID(1)
	return mvc.View{
		Layout: "shared/productLayout.html",
		Name:   "product/view.html",
		Data: iris.Map{
			"product": product,
		},
	}
}

func (p *ProductController) GetOrder() []byte {
	var (
		productID int64
		err       error
	)
	userIDSting := p.Ctx.GetCookie("uid")
	productIDstr := p.Ctx.URLParam("productID")
	productID, err = strconv.ParseInt(productIDstr, 10, 64)
	if err != nil {
		golog.Error("数据格式错误", err)
		return []byte("false")
	}
	userID, err := strconv.ParseInt(userIDSting, 10, 64)
	if err != nil {
		golog.Error("userID转换int64错误", err)
		return []byte("false")

	}
	message := datamodels.NewMessage(productID, userID)
	messageByte, err := json.Marshal(message)
	if err != nil {
		golog.Error("message转换为json错误", err)
		return []byte("false")
	}
	err = p.RabbitMQ.PublishSimple(string(messageByte))
	if err != nil {
		golog.Error(err)
		return []byte("false")
	}
	return []byte("true")

	//product, err := p.ProductService.GetProductByID(int64(productID))
	//if err != nil {
	//	p.Ctx.Application().Logger().Debug(err)
	//}
	//if product.ProductNum > 0 {
	//	// 扣除商品数量
	//	product.ProductNum -= 1
	//	err := p.ProductService.UpdateProduct(product)
	//	if err != nil {
	//		p.Ctx.Application().Logger().Debug(err)
	//	}
	//	// 创建订单
	//	userID, err := strconv.Atoi(userIDSting)
	//	order := &datamodels.Order{
	//		UserId:      int64(userID),
	//		ProductId:   int64(productID),
	//		OrderStatus: datamodels.OrderSuccess,
	//	}
	//
	//	orderID, err = p.OrderService.InsertOrder(order)
	//	if err != nil {
	//		p.Ctx.Application().Logger().Debug(err)
	//	} else {
	//		showMessage = "抢购成功"
	//	}
	//}
	//
	//return mvc.View{
	//	Layout: "shared/productLayout.html",
	//	Name:   "product/result.html",
	//	Data: iris.Map{
	//		"orderID":     orderID,
	//		"showMessage": showMessage,
	//	},
	//}
}
