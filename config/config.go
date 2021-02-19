package config

import "github.com/shopspring/decimal"

var Host = "api.huobi.pro"
var AccessKey = "e31fc594-0b8423f3-hrf5gdfghe-a4540"
var SecretKey = "fa646ab6-40e9658b-73b8b610-92699"
var AccountId = "6643727"
var SubUid int64 = 5678
var SubUids string = "5678"

//购买几次不重复的币
var Count = 2

//最大最小涨幅
var MinIncrease = decimal.NewFromFloat(0.15)
var MaxIncrease = decimal.NewFromFloat(0.3)

//购买多少usdt的币
var Utotal = 100.0
