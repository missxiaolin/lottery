package controllers

import (
	"fmt"
	"github.com/kataras/iris"
	"lottery/comm"
	"lottery/conf"
	"lottery/models"
	"lottery/services"
	"lottery/web/utils"
	"strconv"
	"time"
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

func (c *IndexController) Get() string {
	c.Ctx.Header("Content-Type", "text/html")
	return "welcome to 使用Go语言实现抽奖系统，<a href='/public/home.html'>开始抽奖</a>"
}

// http://localhost:8080/gifts
func (c *IndexController) GetGifts() map[string]interface{} {
	rs := make(map[string]interface{})
	rs["code"] = 0
	rs["msg"] = ""
	dataList := c.ServiceGift.GetAll(true)
	list := make([]models.LtGift, 0)
	for _, data := range dataList {
		// 正常状态的才需要放进来
		if data.SysStatus == 0 {
			list = append(list, data)
		}
	}

	rs["gifts"] = list
	return rs
}

// http://localhost:8080/newprize
func (c *IndexController) GetNewprize() map[string]interface{} {
	rs := make(map[string]interface{})
	rs["code"] = 0
	rs["msg"] = ""
	gifts := c.ServiceGift.GetAll(true)
	giftIds := []int{}
	for _, data := range gifts {
		// 虚拟券或者实物奖才需要放到外部榜单中展示
		if data.Gtype > 1 {
			giftIds = append(giftIds, data.Id)
		}
	}
	list := c.ServiceResult.GetNewPrize(50, giftIds)
	rs["prize_list"] = list
	return rs
}

// http://localhost:8080/myprize
func (c *IndexController) GetMyprize() map[string]interface{} {
	rs := make(map[string]interface{})
	rs["code"] = 0
	rs["msg"] = ""
	// 验证登录
	loginuser := comm.GetLoginUser(c.Ctx.Request())
	if loginuser == nil || loginuser.Uid < 1 {
		rs["code"] = 101
		rs["msg"] = "请先登录，再来抽奖"
		return rs
	}
	// 只读取出来最新的100次中奖记录
	list := c.ServiceResult.SearchByUser(loginuser.Uid, 1, 100)
	rs["prize_list"] = list
	// 今天抽奖次数
	day, _ := strconv.Atoi(comm.FormatFromUnixTimeShort(time.Now().Unix()))
	num := c.ServiceUserday.Count(loginuser.Uid, day)
	rs["prize_num"] = conf.UserPrizeMax - num
	return rs
}

// 登录 GET /login
func (c *IndexController) GetLogin() {
	// 随机生成登录信息uid
	uid := comm.Random(100000)
	loginuser := models.ObjLoginuser{
		Uid: uid,
		Username: fmt.Sprintf("admin-%d", uid),
		Now: comm.NowUnix(),
		Ip: comm.ClientIP(c.Ctx.Request()),
	}
	refer := c.Ctx.GetHeader("Referer")
	if refer == "" {
		refer = "/public/home.html?from=login"
	}
	comm.SetLoginuser(c.Ctx.ResponseWriter(), &loginuser)
	comm.Redirect(c.Ctx.ResponseWriter(), refer)
}

// 退出 /logout
func (c * IndexController) GetLogout() {
	refer := c.Ctx.GetHeader("Referer")
	if refer == "" {
		refer = "/public/home.html?from=logout"
	}
	comm.SetLoginuser(c.Ctx.ResponseWriter(), nil)
	comm.Redirect(c.Ctx.ResponseWriter(), refer)
}

// 验证用户的IP，今天的抽奖次数是否超过每天最大允许的参与次数
func (c *IndexController) GetCheckLimitIpday(ip string) bool {
	num := utils.IncrIpLucyNum(ip)
	if num > conf.IpLimitMax {
		return false
	} else if num > conf.IpPrizeMax {
		return false
	}
	return true
}

// 验证用户的IP，今天的抽奖次数是否超过每天最大抽奖次数
func (c *IndexController) GetCheckIpday(ip string) bool {
	num := utils.IncrIpLucyNum(ip)
	if num > conf.IpPrizeMax {
		return false
	}
	return true
}