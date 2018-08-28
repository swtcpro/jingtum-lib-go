/**
 *
 * 公私钥接口定义类
 *
 * @FileName: keypairs.go
 * @Auther : 杨雪波
 * @Email : yangxuebo@yeah.net
 * @CreateTime: 2018-06-16 10:44:32
 * @UpdateTime: 2018-06-16 10:44:54
 * Copyright@2013 版权所有
 */

package crypto

import (
	"jingtumlib/crypto/secp256k1"
)

//KeyPair KeyPair
type KeyPair interface {
	//根据私钥获取秘钥对
	DeriveKeyPair(secret string) (*secp256k1.PrivateKey, error)

	//地址格式验证
	CheckAddress(address string) bool
	//生成私钥
	GenerateSeed() (string, error)
}
