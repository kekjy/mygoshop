package controllers

import (
	"fmt"
	"mygoshop/common"
	"mygoshop/config"
	"mygoshop/datamodels"
	"mygoshop/services"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type ProductController struct {
	Ctx            iris.Context
	ProductService services.IProductService
}

func (p *ProductController) Get() mvc.View {
	ProductArray, err := p.ProductService.GetAllProduct()
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	ProductArray = append(ProductArray, make([]*datamodels.Product, 8)...)
	return mvc.View{
		Layout: "",
		Name:   "product/menhu.html",
		Data: iris.Map{
			"product0": ProductArray[0],
			"product1": ProductArray[1],
			"product2": ProductArray[2],
			"product3": ProductArray[3],
			"product4": ProductArray[4],
			"product5": ProductArray[5],
			"product6": ProductArray[6],
			"product7": ProductArray[7],
		},
	}
}

// product/detail?productid=1
func (p *ProductController) GetDetail() mvc.View {
	productId, err := p.Ctx.URLParamInt64("productid")
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	product, err := p.ProductService.GetProductById(productId)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	return mvc.View{
		Layout: "",
		Name:   "product/product.html",
		Data: iris.Map{
			"product": product,
		},
	}
}

// product/order?productid=1
func (p *ProductController) GetOrder() string {
	productId, err := p.Ctx.URLParamInt64("productid")
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	fmt.Println("productID", productId)

	hostUrl := fmt.Sprintf("http://%s:%s/onsale?productid=%s",
		config.ValidateSet.Host,
		config.ValidateSet.Port,
		strconv.FormatInt(productId, 10),
	)

	fmt.Println("validate : ", hostUrl)
	response, body, err := common.GetCurl(hostUrl, p.Ctx.Request())
	if err != nil {
		return "server error"
	}

	// 判断状态
	if response.StatusCode == 200 {
		if string(body) == "true" {
			return "true"
		} else {
			return "false"
		}
	} else {
		return "false"
	}

}
