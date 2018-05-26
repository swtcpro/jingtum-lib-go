package goMath

import (
	"math/rand"
)

const (
	RandomStr = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)

//获取指定位数的随机字符串
func GetRandomStr(lenth int) string {
	sour := make([]byte, lenth)
	for i := range sour {
		sour[i] = RandomStr[rand.Intn(len(RandomStr))]
	}
	return string(sour)
}
