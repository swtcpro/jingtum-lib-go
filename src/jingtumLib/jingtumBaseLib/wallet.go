/*** 钱包类
 *** wallet.go
 *** 主要用于创建和管理井通钱包
 *** author:            1416205324@qq.com
 *** last_modified_time:  2018-5-25 12:30:01
 */

package jingtumBaseLib

import (
	"common/goMath"
	"crypto/sha256"
	"fmt"
    "errors"
)

const (
	RandomLen = 16
)

type Wallet struct {
    priv *PrivateKey
}

//创建一个新钱包
func Generate() {
	randBytes := goMath.GetRandomStr(RandomLen)
	sha := sha256.New()
	sha.Write([]byte(randBytes))
	sRandBytes := sha.Sum(nil)
	fmt.Println(sRandBytes)
	//Debug(sRandBytes)
}

func IsValidAddress(address string) bool {
    return CheckAddress(address)
}

func FromSecret(secret string) *Wallet,error {
    keyPair KeyPair = &Secp256KeyPair{}
    priv,err := keyPair.DeriveKeyPair("snsYqv2FsYLuibE9TGHdG5x5V5Qcn")
    
    if nil != err {
        return nil, err
    }
    wallet Wallet = &Wallet{}
    wallet.priv = priv
    return wallet,nil
}

func (wallet *Wallet) GetPublicKey() PublicKey {
    return wallet.priv.PublicKey
}
