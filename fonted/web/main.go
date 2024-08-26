package main

import (
	"context"
	"fmt"
	"mygoshop/config"
	"mygoshop/db"
	"mygoshop/fonted/middleware"
	"mygoshop/fonted/web/controllers"
	"mygoshop/repositories"
	"mygoshop/services"

	_ "github.com/go-sql-driver/mysql"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func main() {
	app := iris.New()

	app.Get("/", func(ctx iris.Context) {
		ctx.Redirect("/user/")
	})

	app.Logger().SetLevel("debug")

	tempalte := iris.HTML("./fonted/web/views", ".html").Reload(true)
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
	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	userParty := app.Party("/user")
	userParty.Use(middleware.IsLogIn)
	user := mvc.New(userParty)
	user.Register(ctx, userService)
	user.Handle(new(controllers.UserController))

	productRepository := repositories.NewProductManager(db)
	productService := services.NewProductService(productRepository)
	productParty := app.Party("/product")

	// 验证中间件
	productParty.Use(middleware.AuthProduct)

	orderRepository := repositories.NewOrderManager(db)
	orderService := services.NewOrderService(orderRepository)

	product := mvc.New(productParty)
	product.Register(ctx, productService, orderService)
	product.Handle(new(controllers.ProductController))

	// order 服务

	// 启动服务
	app.Run(
		iris.Addr(config.FontedSet.Host+":"+config.FontedSet.Port),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}
