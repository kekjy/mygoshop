package services

import (
	"mygoshop/datamodels"
	"mygoshop/repositories"
)

type IOrderService interface {
	GetOrderById(int64) (*datamodels.Order, error)
	DeleteOrderById(int64) bool
	UpdateOrder(order *datamodels.Order) error
	InsertOrder(order *datamodels.Order) (int64, error)
	GetAllOrder() ([]*datamodels.Order, error)
	GetAllOrderInfo() (map[int]map[string]string, error)
	InsertOrderByMessage(message *datamodels.Message) (int64, error)
}

type OrderService struct {
	OrderRepository repositories.IOrder
}

func NewOrderService(repository repositories.IOrder) IOrderService {
	return &OrderService{
		OrderRepository: repository,
	}
}

func (o *OrderService) GetOrderById(orderId int64) (*datamodels.Order, error) {
	return o.OrderRepository.SelectByKey(orderId)
}

func (o *OrderService) DeleteOrderById(orderId int64) bool {
	return o.OrderRepository.Delete(orderId)
}

func (o *OrderService) UpdateOrder(order *datamodels.Order) error {
	return o.OrderRepository.Update(order)
}

func (o *OrderService) InsertOrder(order *datamodels.Order) (int64, error) {
	return o.OrderRepository.Insert(order)
}

func (o *OrderService) GetAllOrder() ([]*datamodels.Order, error) {
	return o.OrderRepository.SelectAll()
}

func (o *OrderService) GetAllOrderInfo() (map[int]map[string]string, error) {
	return o.OrderRepository.SelectAllWithMap()
}

// 根据消息创建订单
func (o *OrderService) InsertOrderByMessage(message *datamodels.Message) (int64, error) {
	order := &datamodels.Order{
		UserId:      message.UserID,
		ProductId:   message.ProductID,
		OrderStatus: datamodels.OrderSuccess,
	}
	return o.InsertOrder(order)
}
