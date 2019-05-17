package controllers

import (
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

	// 4 验证IP今日的参与次数

	// 5 验证IP黑名单

	// 6 验证用户黑名单

	// 7 获得抽奖编码

	// 8 匹配奖品是否中奖

	// 9 有限制奖品发放

	// 10 不同编码的优惠券的发放

	// 11 记录中奖记录

	// 12 返回抽奖结果
	return 0, "", nil
}