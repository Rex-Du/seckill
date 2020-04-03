// Author : rexdu
// Time : 2020-03-25 22:32
package services

import (
	"seckill/datamodels"
	"seckill/repositories"
)

type IOrderService interface {
	GetOrderByID(orderID int64) (order *datamodels.Order, err error)
	DeleteOrderByID(orderID int64) bool
	UpdateOrder(order *datamodels.Order) (err error)
	InsertOrder(order *datamodels.Order) (orderID int64, err error)
	GetAllOrders() (orders []*datamodels.Order, err error)
	GetAllOrderInfo() (orderInfo map[int]map[string]string, err error)
	InsertOrderByMessage(*datamodels.Message) (int64, error)
}

type OrderService struct {
	OrderRepository repositories.IOrderRepository
}

func NewOrderService(repo repositories.IOrderRepository) IOrderService {
	return &OrderService{OrderRepository: repo}
}

func (o *OrderService) GetOrderByID(orderID int64) (order *datamodels.Order, err error) {
	return o.OrderRepository.SelectByKey(orderID)
}

func (o *OrderService) DeleteOrderByID(orderID int64) bool {
	return o.OrderRepository.Delete(orderID)
}

func (o *OrderService) UpdateOrder(order *datamodels.Order) (err error) {
	return o.OrderRepository.Update(order)
}

func (o *OrderService) InsertOrder(order *datamodels.Order) (orderID int64, err error) {
	return o.OrderRepository.Insert(order)
}

func (o *OrderService) GetAllOrders() (orders []*datamodels.Order, err error) {
	return o.OrderRepository.SelectAll()
}

func (o *OrderService) GetAllOrderInfo() (orderInfo map[int]map[string]string, err error) {
	return o.OrderRepository.SelectAllWithInfo()
}

func (o *OrderService) InsertOrderByMessage(message *datamodels.Message) (int64, error) {
	order := &datamodels.Order{
		UserId:      message.UserId,
		ProductId:   message.ProductID,
		OrderStatus: datamodels.OrderSuccess,
	}
	return o.OrderRepository.Insert(order)
}
