package main

import (
	"context"
	"fmt"
	"mygoshop/backend/web/controllers"
	"mygoshop/config"
	"mygoshop/db"
	"mygoshop/repositories"
	"mygoshop/services"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func main() {
	app := iris.New()

	//设置错误模式
	app.Logger().SetLevel("debug")

	app.Get("/", func(ctx iris.Context) {
		ctx.Redirect("/order/index")
	})

	// 注册模板
	tempalte := iris.HTML("./backend/web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(tempalte)

	app.HandleDir("/assets", "./assets")

	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewLayout("")
		ctx.ViewData("Message", ctx.Values().GetStringDefault("Message", "visit error"))
		ctx.View("shared/error.html")
	})

	// 连接数据库
	db, err := db.NewDbConn()
	if err != nil {
		fmt.Println(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 注册控制器
	productRepository := repositories.NewProductManager(db)
	productService := services.NewProductService(productRepository)
	productParty := app.Party("/product")
	product := mvc.New(productParty)
	product.Register(ctx, productService)
	product.Handle(new(controllers.ProductController))

	orderRepository := repositories.NewOrderManager(db)
	orderService := services.NewOrderService(orderRepository)
	orderParty := app.Party("/order")
	order := mvc.New(orderParty)
	order.Register(ctx, orderService)
	order.Handle(new(controllers.OrderController))

	// 启动服务
	app.Run(
		iris.Addr(config.BackendSet.Host+":"+config.BackendSet.Port),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)

}
