/**
 * 钱包类，用于创建和导入钱包等功能.
 *
 * @FileName: wallet.go
 * @Auther : 杨雪波
 * @Email : yangxuebo@yeah.net
 * @CreateTime: 2018-07-26 10:44:32
 * @UpdateTime: 2018-07-26 10:44:54
 */
package jingtumLib

import (
	jtConst "jingtumLib/constant"
	"jingtumLib/crypto/secp256k1"
)

//钱包结构体
type Wallet struct {
	priv *secp256k1.PrivateKey
}

/**
 * 根据井通私钥创建钱包
 */
func FromSecret(secret string) (*Wallet, error) {
	if secret == "" {
		return nil, jtConst.ERR_EMPTY_PARAM
		//fmt.Errorf("invalid merkle root (remote: %x local: %x)", header.Root, root)
	}
	keyPair := &secp256k1.Secp256KeyPair{}
	priv, err := keyPair.DeriveKeyPair(secret)
	if nil != err {
		return nil, err
	}
	wallet := new(Wallet)
	wallet.priv = priv
	return wallet, nil
}
