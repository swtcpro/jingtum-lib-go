/***  产出随机数
 *** random.go
 *** 主要用于产生随机数
 *** author:              1416205324@qq.com
 *** last_modified_time:  2018-5-25 13:13:23
 */
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
