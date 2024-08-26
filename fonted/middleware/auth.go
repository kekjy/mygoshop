package middleware

import (
	"mygoshop/common"

	"github.com/kataras/iris/v12"
)

func AuthProduct(ctx iris.Context) {
	__uid, _ := ctx.GetRequestCookie("uid")
	__sign, _ := ctx.GetRequestCookie("sign")
	sign := __sign.Value
	uid := __uid.Value
	if uid == "" || sign == "" {
		ctx.Application().Logger().Debug("Must log in first")
		ctx.Redirect("/user/login")
		return
	}
	__signString, err := common.EnPwdCode([]byte(uid))
	if err != nil {
		common.GlobalCookie(ctx, "uid", "")
		common.GlobalCookie(ctx, "sign", "")
		ctx.Application().Logger().Debug("Sign Error!")
		ctx.Redirect("/user/login")
		return
	}
	if __signString != sign {
		common.GlobalCookie(ctx, "uid", "")
		common.GlobalCookie(ctx, "sign", "")
		ctx.Application().Logger().Debug("Cookie Error!")
		ctx.Redirect("/user/login")
		return
	}
	ctx.Application().Logger().Debug("Already logged in")
	ctx.Next()
	ctx.Redirect("/user/login")
}

func IsLogIn(ctx iris.Context) {
	__uid, _ := ctx.GetRequestCookie("uid")
	__sign, _ := ctx.GetRequestCookie("sign")
	sign := __sign.Value
	uid := __uid.Value
	if uid != "" && sign != "" {
		if __signString, err := common.EnPwdCode([]byte(uid)); err == nil && __signString == sign {
			ctx.Redirect("/product")
			return
		}
	}
	common.GlobalCookie(ctx, "uid", "")
	common.GlobalCookie(ctx, "sign", "")
	ctx.Next()
}
