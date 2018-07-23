/**
 *
 * 文件功能介绍
 *
 * @FileName: keypairs.go
 * @Auther : 杨雪波
 * @Email : yangxuebo@yeah.net
 * @CreateTime: 2018-06-16 10:44:32
 * @UpdateTime: 2018-06-16 10:44:54
 * Copyright@2013 版权所有
 */

package jingtumBaseLib

import (
	"fmt"
)

/**
 *  接口定义
 */
type KeyPair interface {
	//根据私钥获取秘钥对
	DeriveKeyPair(secret string) (*PrivateKey, error)

	//地址格式验证
	CheckAddress(address string) bool
}

func DecodeAddress(address string) []byte {
	decodedBytes, err := __decode(ACCOUNT_PREFIX, address)
	if err != nil {
		panic(fmt.Sprintf("Issuer invalid issuer info %v", address))
	}

	return decodedBytes
}

//func GenerateSeed () {
//var randBytes = brorand(16);
//return __encode(SEED_PREFIX, randBytes);
//}
