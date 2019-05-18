package utils

import (
	"fmt"
	"log"
	"lottery/comm"
	"lottery/datasource"
	"math"
)

const ipFrameSize = 2

func IncrIpLuckyNum(strIp string) int64 {
	ip := comm.Ip4toInt(strIp)
	i := ip % ipFrameSize
	key := fmt.Sprintf("day_ips_%d", i)
	cacheObj := datasource.InstanceCache()
	rs, err := cacheObj.Do("HINCRBY", key, ip, 1)
	if err != nil {
		log.Println("ip_day_lucky redis HINCRBY error=", err)
		return math.MaxInt32
	} else {
		return rs.(int64)
	}
}