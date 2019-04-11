package conf

import "time"

const SysTimeform = "2006-01-02 15:04:05"
const SysTimeformShort = "2006-01-02"

// 中国时区
var SysTimeLocation, _ = time.LoadLocation("Asia/Chongqing")

// ObjSalesign 签名密钥
var SignSecret = []byte("0123456789abcdef")

// cookie中的加密验证密钥
var CookieSecret = "hellolottery"