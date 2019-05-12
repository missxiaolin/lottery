package controllers

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"lottery/services"
)

type AdminController struct {
	Ctx iris.Context
	ServiceUser services.UserService
	ServiceGift services.GiftService
	ServiceCode services.CodeService
	ServiceResult services.ResultService
	ServiceUserday services.UserdayService
	ServiceBlackip services.BlackipService
}

func (c *AdminController) Get() mvc.Result {
	return mvc.View{
		Name: "admin/index.html",
		Data: iris.Map{
			"Title": "管理后台",
			"Channel":"",
		},
		Layout: "admin/layout.html",
	}
}