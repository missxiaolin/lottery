package utils

import (
	"fmt"
	"log"
	"lottery/datasource"
	"math"
)

const userFrameSize = 2

// 今天的用户抽奖次数递增，返回递增后的数值
func IncrUserLuckyNum(uid int) int64 {
	i := uid % userFrameSize
	// 集群的redis统计数递增
	return incrServUserLucyNum(i, uid)
}

func incrServUserLucyNum(i, uid int) int64 {
	key := fmt.Sprintf("day_users_%d", i)
	cacheObj := datasource.InstanceCache()
	rs, err := cacheObj.Do("HINCRBY", key, uid, 1)
	if err != nil {
		log.Println("user_day_lucky redis HINCRBY key=", key,
			", uid=", uid, ", err=", err)
		return math.MaxInt32
	} else {
		num := rs.(int64)
		return num
	}
}

// 从给定的数据直接初始化用户的参与次数
func InitUserLuckyNum(uid int, num int64) {
	if num <= 1 {
		return
	}
	i := uid % userFrameSize
	// 集群
	initServUserLuckyNum(i, uid, num)
}

func initServUserLuckyNum(i, uid int, num int64) {
	key := fmt.Sprintf("day_users_%d", i)
	cacheObj := datasource.InstanceCache()
	_, err := cacheObj.Do("HSET", key, uid, num)
	if err != nil {
		log.Println("user_day_lucky redis HSET key=", key,
			", uid=", uid, ", err=", err)
	}
}