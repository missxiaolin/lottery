package controllers

import (
	"github.com/kataras/iris"
	"lottery/services"
)

type IndexController struct {
	Ctx iris.Context
	ServiceUser services.UserService
	ServiceGift services.GiftService
	ServiceCode services.CodeService
	ServiceResult services.ResultService
	ServiceUserday services.UserdayService
	ServiceBlackip services.BlackipService
}

func (c IndexController) Get() string {
	c.Ctx.Header("Content-Type", "text/html")
	return "welcome to 使用Go语言实现抽奖系统，<a href='/public/index.html'>开始抽奖</a>"
}