package controllers

import (
	"fmt"
	"log"
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
