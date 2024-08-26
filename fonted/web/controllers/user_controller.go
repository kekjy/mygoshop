package controllers

import (
	"fmt"
	"mygoshop/common"
	"mygoshop/datamodels"
	"mygoshop/services"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type UserController struct {
	Ctx         iris.Context
	UserService services.IUserService
}

func (c *UserController) Get() mvc.View {
	return mvc.View{
		Layout: "",
		Name:   "user/welcome.html",
	}
}

func (c *UserController) GetRegister() mvc.View {
	return mvc.View{
		Layout: "",
		Name:   "user/register.html",
	}
}

func (c *UserController) GetLogin() mvc.View {
	return mvc.View{
		Layout: "",
		Name:   "user/login.html",
	}
}

func (c *UserController) GetLogout() {
	common.GlobalCookie(c.Ctx, "uid", "")
	common.GlobalCookie(c.Ctx, "sign", "")
	c.Ctx.Redirect("/user/login")
}

func (c *UserController) GetLoginerror() mvc.View {
	return mvc.View{
		Layout: "",
		Name:   "user/login.html",
		Data: iris.Map{
			"showMessage": "用户名或密码错误，请重试",
		},
	}
}

func (c *UserController) PostRegister() mvc.View {
	var (
		nickName = c.Ctx.FormValue("nickName")
		userName = c.Ctx.FormValue("userName")
		pwd      = c.Ctx.FormValue("password")
	)
	user := &datamodels.User{
		NickName:     nickName,
		UserName:     userName,
		HashPassword: pwd,
	}
	fmt.Println(user)
	_, err := c.UserService.AddUser(user)
	fmt.Println(err)
	if err != nil {
		c.Ctx.Redirect("/user/error")
		return mvc.View{
			Layout: "",
			Name:   "user/register.html",
			Data: iris.Map{
				"showMessage": "用户名已存在或发生未知错误",
			},
		}
	}
	return c.GetLogin()
}

func (c *UserController) PostLogin() mvc.Response {
	// 1.获取用户提交的表单信息
	var (
		userName = c.Ctx.FormValue("userName")
		pwd      = c.Ctx.FormValue("password")
	)

	// 2.验证用户账号密码是否正确
	user, isOk := c.UserService.IsLoginSuccess(userName, pwd)

	// Login Failed
	if !isOk {
		return mvc.Response{
			Path: "loginerror",
		}
	}

	// 写入用户 ID 到 Cookie 中
	common.GlobalCookie(c.Ctx, "uid", strconv.FormatInt(user.ID, 10))
	uidByte := strconv.FormatInt(user.ID, 10)
	uidString, err := common.EnPwdCode([]byte(uidByte))
	if err != nil {
		fmt.Println(err)
	}
	common.GlobalCookie(c.Ctx, "sign", uidString)

	return mvc.Response{
		Path: "/product",
	}
}
