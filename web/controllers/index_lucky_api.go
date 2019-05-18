package controllers

import (
	"fmt"
	"log"
	"lottery/comm"
	"lottery/conf"
	"lottery/models"
	"lottery/services"
	"lottery/web/utils"
)

type LuckyApi struct {
}

func (api *LuckyApi)luckyDo(uid int, username, ip string) (int, string, *models.ObjGiftPrize) {
	// 2 用户抽奖分布式锁定
	ok := utils.LockLucky(uid)
	if ok {
		defer utils.UnlockLucky(uid)
	} else {
		return 102, "正在抽奖，请稍后重试", nil
	}

	// 3 验证用户今日参与次数
	userDayNum := utils.IncrUserLuckyNum(uid)
	if userDayNum > conf.UserPrizeMax {
		return 103, "今日的抽奖次数已用完，明天再来吧", nil
	} else {
		ok = api.checkUserday(uid, userDayNum)
		if !ok {
			return 103, "今日的抽奖次数已用完，明天再来吧", nil
		}
	}

	// 4 验证IP今日的参与次数
	ipDayNum := utils.IncrIpLuckyNum(ip)
	if ipDayNum > conf.IpLimitMax {
		return 104, "相同IP参与次数太多，明天再来参与吧", nil
	}

	limitBlack := false // 黑名单
	if ipDayNum > conf.IpPrizeMax {
		limitBlack = true
	}
	// 5 验证IP黑名单
	var blackipInfo *models.LtBlackip
	if !limitBlack {
		ok, blackipInfo = api.checkBlackip(ip)
		if !ok {
			fmt.Println("黑名单中的IP", ip, limitBlack)
			limitBlack = true
		}
	}

	// 6 验证用户黑名单
	var userInfo *models.LtUser
	if !limitBlack {
		ok, userInfo = api.checkBlackUser(uid)
		if !ok {
			limitBlack = true
		}
	}

	// 7 获得抽奖编码
	prizeCode := comm.Random(10000)

	// 8 匹配奖品是否中奖
	prizeGift := api.prize(prizeCode, limitBlack)
	if prizeGift == nil ||
		prizeGift.PrizeNum < 0 ||
		(prizeGift.PrizeNum > 0 && prizeGift.LeftNum <= 0) {
		return 205, "很遗憾，没有中奖，请下次再试", nil
	}

	// 9 有限制奖品发放
	if prizeGift.PrizeNum > 0 {
		if utils.GetGiftPoolNum(prizeGift.Id) <= 0 {
			return 206, "很遗憾，没有中奖，请下次再试", nil
		}
		ok = utils.PrizeGift(prizeGift.Id, prizeGift.LeftNum)
		if !ok {
			return 207, "很遗憾，没有中奖，请下次再试", nil
		}
	}

	// 10 不同编码的优惠券的发放
	if prizeGift.Gtype == conf.GtypeCodeDiff {
		code := utils.PrizeCodeDiff(prizeGift.Id, services.NewCodeService())
		if code == "" {
			return 208, "很遗憾，没有中奖，请下次再试", nil
		}
		prizeGift.Gdata = code
	}

	// 11 记录中奖记录
	result := models.LtResult{
		GiftId:     prizeGift.Id,
		GiftName:   prizeGift.Title,
		GiftType:   prizeGift.Gtype,
		Uid:        uid,
		Username:   username,
		PrizeCode:  prizeCode,
		GiftData:   prizeGift.Gdata,
		SysCreated: comm.NowUnix(),
		SysIp:      ip,
		SysStatus:  0,
	}
	err := services.NewResultService().Create(&result)
	if err != nil {
		log.Println("index_lucky.GetLucky ServiceResult.Create ", result,
			", error=", err)
		return 209, "很遗憾，没有中奖，请下次再试", nil
	}
	if prizeGift.Gtype == conf.GtypeGiftLarge {
		// 如果获得了实物大奖，需要将用户、IP设置成黑名单一段时间
		api.prizeLarge(ip, uid, username, userInfo, blackipInfo)
	}

	// 12 返回抽奖结果
	return 0, "", prizeGift
}