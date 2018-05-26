package jingtumLib

import (
	"common/goMath"
	"crypto/sha256"
	"fmt"
)

const (
	RandomLen = 16
)

func Generate() {
	randBytes := goMath.GetRandomStr(RandomLen)
	sha := sha256.New()
	sha.Write([]byte(randBytes))
	sRandBytes := sha.Sum(nil)
	fmt.Println(sRandBytes)
	//Debugf("log in")
}
