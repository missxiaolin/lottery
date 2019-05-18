package controllers

import (
	"fmt"
	"log"
	"lottery/comm"
	"lottery/conf"
	"lottery/models"
	"lottery/services"
	"lottery/web/utils"
	"strconv"
	"time"
)

func (api *LuckyApi) checkUserday(uid int, num int64) bool {
	userdayService := services.NewUserdayService()
	userdayInfo := userdayService.GetUserToday(uid)
	if userdayInfo != nil && userdayInfo.Uid == uid {
		// 今天存在抽奖记录
		if userdayInfo.Num >= conf.UserPrizeMax {
			if int(num) < userdayInfo.Num {
				utils.InitUserLuckyNum(uid, int64(userdayInfo.Num))
			}
			return false
		} else {
			userdayInfo.Num++
			if int(num) < userdayInfo.Num {
				utils.InitUserLuckyNum(uid, int64(userdayInfo.Num))
			}
			err103 := userdayService.Update(userdayInfo, nil)
			if err103 != nil {
				log.Println("index_lucky_check_userday ServiceUserDay.Update " +
					"err103=", err103)
			}
		}
	} else {
		// 创建今天的用户参与记录
		y, m, d := time.Now().Date()
		strDay := fmt.Sprintf("%d%02d%02d", y, m, d)
		day, _ := strconv.Atoi(strDay)
		userdayInfo = &models.LtUserday{
			Uid:        uid,
			Day:        day,
			Num:        1,
			SysCreated: int(time.Now().Unix()),
		}
		err103 := userdayService.Create(userdayInfo)
		if err103 != nil {
			log.Println("index_lucky_check_userday ServiceUserDay.Create " +
				"err103=", err103)
		}
		utils.InitUserLuckyNum(uid, 1)
	}
	return true
}

func (api *LuckyApi) checkBlackip(ip string) (bool, *models.LtBlackip) {
	info := services.NewBlackipService().GetByIp(ip)
	if info == nil || info.Ip == "" {
		return true, nil
	}
	if info.Blacktime > int(time.Now().Unix()) {
		// IP黑名单存在，并且还在黑名单有效期内
		return false, info
	}
	return true, info
}

func (api *LuckyApi) checkBlackUser(uid int) (bool, *models.LtUser) {
	info := services.NewUserService().Get(uid)
		if info != nil && info.Blacktime > int(time.Now().Unix()) {
		// 黑名单存在并且有效
		return false, info
	}
	return true, info
}

func (api *LuckyApi) prize(prizeCode int, limitBlack bool) *models.ObjGiftPrize {
	var prizeGift *models.ObjGiftPrize
	giftList := services.NewGiftService().GetAllUse(true)
	for _, gift := range giftList {
		if gift.PrizeCodeA <= prizeCode &&
			gift.PrizeCodeB >= prizeCode {
			// 中奖编码区间满足条件，说明可以中奖
			if !limitBlack || gift.Gtype < conf.GtypeGiftSmall {
				prizeGift = &gift
				break
			}
		}
	}
	return prizeGift
}

func (api *LuckyApi) prizeLarge(ip string,
	uid int, username string,
	userinfo *models.LtUser,
	blackipInfo *models.LtBlackip) {
	userService := services.NewUserService()
	blackipService := services.NewBlackipService()
	nowTime := comm.NowUnix()
	blackTime := 30 * 86400
	// 更新用户的黑名单信息
	if userinfo == nil || userinfo.Id <= 0 {
		userinfo = &models.LtUser{
			Id:			uid,
			Username:   username,
			Blacktime:  nowTime+blackTime,
			SysCreated: nowTime,
			SysIp:      ip,
		}
		userService.Create(userinfo)
	} else {
		userinfo = &models.LtUser{
			Id: uid,
			Blacktime:nowTime+blackTime,
			SysUpdated:nowTime,
		}
		userService.Update(userinfo, nil)
	}
	// 更新要IP的黑名单信息
	if blackipInfo == nil || blackipInfo.Id <= 0 {
		blackipInfo = &models.LtBlackip{
			Ip:         ip,
			Blacktime:  nowTime+blackTime,
			SysCreated: nowTime,
		}
		blackipService.Create(blackipInfo)
	} else {
		blackipInfo.Blacktime = nowTime + blackTime
		blackipInfo.SysUpdated = nowTime
		blackipService.Update(blackipInfo, nil)
	}
}
