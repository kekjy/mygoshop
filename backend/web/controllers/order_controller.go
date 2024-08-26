package controllers

import (
	"mygoshop/services"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type OrderController struct {
	Ctx          iris.Context
	OrderService services.IOrderService
}

func (o *OrderController) GetAll() mvc.View {
	orderMap, err := o.OrderService.GetAllOrderInfo()

	if err != nil {
		o.Ctx.Application().Logger().Debug(err)
		return mvc.View{
			Name: "shared/error.html",
			Data: iris.Map{
				"Message": "Failed to load orders",
			},
		}
	}

	return mvc.View{
		Name: "order/view.html",
		Data: iris.Map{
			"order": orderMap,
		},
	}
}

func (o *OrderController) GetIndex() mvc.View {
	return mvc.View{
		Name: "order/index.html",
	}
}
