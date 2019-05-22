package cron

import (
	"log"
	"lottery/comm"
	"lottery/services"
	"lottery/web/utils"
	"time"
)

/**
 * 只需要一个应用运行的服务
 * 全局的服务
 */
func ConfigueAppOneCron() {
	// 每5分钟执行一次，奖品的发奖计划到期的时候，需要重新生成发奖计划
	go resetAllGiftPrizeData()
	// 每分钟执行一次，根据发奖计划，把奖品数量放入奖品池
	go distributionAllGiftPool()
}

// 重置所有奖品的发奖计划
// 每5分钟执行一次
func resetAllGiftPrizeData() {
	giftService := services.NewGiftService()
	list := giftService.GetAll(false)
	nowTime := comm.NowUnix()
	for _, giftInfo := range list {
		if giftInfo.PrizeTime != 0 &&
			(giftInfo.PrizeData == "" || giftInfo.PrizeEnd <= nowTime) {
			// 立即执行
			log.Println("crontab start utils.ResetGiftPrizeData giftInfo=", giftInfo)
			utils.ResetGiftPrizeData(&giftInfo, giftService)
			// 预加载缓存数据
			giftService.GetAll(true)
			log.Println("crontab end utils.ResetGiftPrizeData giftInfo")
		}
	}

	// 每5分钟执行一次
	time.AfterFunc(5 * time.Minute, resetAllGiftPrizeData)
}

// 根据发奖计划，把奖品数量放入奖品池
// 每分钟执行一次
func distributionAllGiftPool() {
	log.Println("crontab start utils.DistributionGiftPool")
	num := utils.DistributionGiftPool()
	log.Println("crontab end utils.DistributionGiftPool, num=", num)

	// 每分钟执行一次
	time.AfterFunc(time.Minute, distributionAllGiftPool)
}