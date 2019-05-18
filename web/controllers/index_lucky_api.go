package controllers

import (
	"fmt"
	"lottery/conf"
	"lottery/models"
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

	// 8 匹配奖品是否中奖

	// 9 有限制奖品发放

	// 10 不同编码的优惠券的发放

	// 11 记录中奖记录

	// 12 返回抽奖结果
	return 0, "", nil
}