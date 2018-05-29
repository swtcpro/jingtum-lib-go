/*** 钱包类
 *** wallet.go
 *** 主要用于创建和管理井通钱包
 *** author:            1416205324@qq.com
 *** last_modified_time:  2018-5-25 12:30:01
 */

package jingtumLib

import (
	"common/goMath"
	"crypto/sha256"
	"fmt"
)

const (
	RandomLen = 16
)

//创建一个新钱包
func Generate() {
	randBytes := goMath.GetRandomStr(RandomLen)
	sha := sha256.New()
	sha.Write([]byte(randBytes))
	sRandBytes := sha.Sum(nil)
	fmt.Println(sRandBytes)
	Debug(sRandBytes)
}
