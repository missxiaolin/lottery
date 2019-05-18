package controllers

import (
	"lottery/models"
	"lottery/services"
	"time"
)

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
