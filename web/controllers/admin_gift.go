package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"lottery/comm"
	"lottery/models"
	"lottery/services"
	"lottery/web/utils"
	"lottery/web/viewmodels"
	"time"
)

type AdminGiftController struct {
	Ctx iris.Context
	ServiceUser services.UserService
	ServiceGift services.GiftService
	ServiceCode services.CodeService
	ServiceResult services.ResultService
	ServiceUserday services.UserdayService
	ServiceBlackip services.BlackipService
}

func (c *AdminGiftController) Get() mvc.Result {
	datalist := c.ServiceGift.GetAll(false)
	for i, giftInfo := range datalist {
		// 奖品发放的计划数据
		prizedata := make([][2]int, 0)
		err := json.Unmarshal([]byte(giftInfo.PrizeData), &prizedata)
		if err != nil || prizedata == nil || len(prizedata) < 1 {
			datalist[i].PrizeData = "[]"
		} else {
			newpd := make([]string, len(prizedata))
			for index, pd := range prizedata {
				ct := comm.FormatFromUnixTime(int64(pd[0]))
				newpd[index] = fmt.Sprintf("【%s】: %d", ct , pd[1])
			}
			str, err := json.Marshal(newpd)
			if err == nil && len(str) > 0 {
				datalist[i].PrizeData = string(str)
			} else {
				datalist[i].PrizeData = "[]"
			}
		}
		// 奖品当前的奖品池数量
		num := utils.GetGiftPoolNum(giftInfo.Id)
		datalist[i].Title = fmt.Sprintf("【%d】%s", num, datalist[i].Title)
	}
	total := len(datalist)
	return mvc.View{
		Name: "admin/gift.html",
		Data: iris.Map{
			"Title": "管理后台",
			"Channel":"gift",
			"Datalist": datalist,
			"Total":    total,
		},
		Layout: "admin/layout.html",
	}
}

func (c *AdminGiftController) GetEdit() mvc.Result {
	id := c.Ctx.URLParamIntDefault("id", 0)
	giftInfo := viewmodels.ViewGift{}
	if id > 0 {
		data := c.ServiceGift.Get(id, false)
		if data != nil {
			giftInfo.Id = data.Id
			giftInfo.Title = data.Title
			giftInfo.PrizeNum = data.PrizeNum
			giftInfo.PrizeCode = data.PrizeCode
			giftInfo.PrizeTime = data.PrizeTime
			giftInfo.Img = data.Img
			giftInfo.Displayorder = data.Displayorder
			giftInfo.Gtype = data.Gtype
			giftInfo.Gdata = data.Gdata
			giftInfo.TimeBegin = comm.FormatFromUnixTime(int64(data.TimeBegin))
			giftInfo.TimeEnd = comm.FormatFromUnixTime(int64(data.TimeEnd))
		}
	}
	return mvc.View{
		Name: "admin/giftEdit.html",
		Data: iris.Map{
			"Title":   "管理后台",
			"Channel": "gift",
			"info":    giftInfo,
		},
		Layout: "admin/layout.html",
	}
}

func (c *AdminGiftController) PostSave() mvc.Result {
	data := viewmodels.ViewGift{}
	err := c.Ctx.ReadForm(&data)
	//fmt.Printf("%v\n", info)
	if err != nil {
		fmt.Println("admin_gift.PostSave ReadForm error=", err)
		return mvc.Response{
			Text: fmt.Sprintf("ReadForm转换异常, err=%s", err),
		}
	}
	giftInfo := models.LtGift{}
	giftInfo.Id = data.Id
	giftInfo.Title = data.Title
	giftInfo.PrizeNum = data.PrizeNum
	giftInfo.PrizeCode = data.PrizeCode
	giftInfo.PrizeTime = data.PrizeTime
	giftInfo.Img = data.Img
	giftInfo.Displayorder = data.Displayorder
	giftInfo.Gtype = data.Gtype
	giftInfo.Gdata = data.Gdata
	t1, err1 := comm.ParseTime(data.TimeBegin)
	t2, err2 := comm.ParseTime(data.TimeEnd)
	if err1 != nil || err2 != nil {
		return mvc.Response{
			Text: fmt.Sprintf("开始时间、结束时间的格式不正确, err1=%s, err2=%s", err1, err2),
		}
	}
	giftInfo.TimeBegin = int(t1.Unix())
	giftInfo.TimeEnd = int(t2.Unix())
	fmt.Println(giftInfo.Id)
	if giftInfo.Id > 0 {
		datainfo := c.ServiceGift.Get(giftInfo.Id, false)
		if datainfo != nil {
			giftInfo.SysUpdated = int(time.Now().Unix())
			giftInfo.SysIp = comm.ClientIP(c.Ctx.Request())
			// 对比修改的内容项
			if datainfo.PrizeNum != giftInfo.PrizeNum {
				// 奖品总数量发生了改变
				giftInfo.LeftNum = datainfo.LeftNum - datainfo.PrizeNum - giftInfo.PrizeNum
				if giftInfo.LeftNum < 0 || giftInfo.PrizeNum <= 0 {
					giftInfo.LeftNum = 0
				}
				giftInfo.SysStatus = datainfo.SysStatus
				utils.ResetGiftPrizeData(&giftInfo, c.ServiceGift)
			} else {
				giftInfo.LeftNum = giftInfo.PrizeNum
			}
			if datainfo.PrizeTime != giftInfo.PrizeTime {
				// 发奖周期发生了变化
				utils.ResetGiftPrizeData(&giftInfo, c.ServiceGift)
			}
			c.ServiceGift.Update(&giftInfo, []string{"title", "prize_num", "left_num", "prize_code", "prize_time",
				"img", "displayorder", "gtype", "gdata", "time_begin", "time_end", "sys_updated"})
		} else {
			giftInfo.Id = 0
		}
	}
	if giftInfo.Id == 0 {
		giftInfo.LeftNum = giftInfo.PrizeNum
		giftInfo.SysIp = comm.ClientIP(c.Ctx.Request())
		giftInfo.SysCreated =  int(time.Now().Unix())
		c.ServiceGift.Create(&giftInfo)
		// 更新奖品的发奖计划
		utils.ResetGiftPrizeData(&giftInfo, c.ServiceGift)
	}
	return mvc.Response{
		Path: "/admin/gift",
	}
}

func (c *AdminGiftController) GetDelete() mvc.Result {
	id, err := c.Ctx.URLParamInt("id")
	if err == nil {
		c.ServiceGift.Delete(id)
	}
	return mvc.Response{
		Path: "/admin/gift",
	}
}

func (c *AdminGiftController) GetReset() mvc.Result {
	id, err := c.Ctx.URLParamInt("id")
	if err == nil {
		c.ServiceGift.Update(&models.LtGift{Id:id, SysStatus:0}, []string{"sys_status"})
	}
	return mvc.Response{
		Path: "/admin/gift",
	}
}