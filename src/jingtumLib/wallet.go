// Package jingtumlib 钱包类，用于创建和导入钱包等功能。
// @FileName: wallet.go
// @Auther : 杨雪波
// @Email : yangxuebo@yeah.net
// @CreateTime: 2018-07-26 10:44:32
// @UpdateTime: 2018-07-26 10:44:54
package jingtumlib

import (
	"fmt"
	"crypto/ecdsa"
	
	"jingtumlib/constant"
	"jingtumlib/crypto/secp256k1"
	"jingtumlib/utils"

	"github.com/btcsuite/btcd/btcec"
)

//Wallet 钱包结构体
type Wallet struct {
	priv   *secp256k1.PrivateKey
	secret string
}

//IsValidAddress 钱包地址合法性验证
func IsValidAddress(address string) bool {
	if address == "" {
		return false
	}

	return utils.IsValidAddress(address)
}

//IsValidSecret 钱包私钥合法性验证
func IsValidSecret(secret string) bool {
	if secret == "" {
		return false
	}

	keyPair := &secp256k1.Secp256KeyPair{}
	_, err := keyPair.DeriveKeyPair(secret)
	if nil != err {
		return false
	}

	return true
}

//Generate 生成钱包
func Generate() (*Wallet, error) {
	keyPair := &secp256k1.Secp256KeyPair{}
	secret, err := keyPair.GenerateSeed()
	if err != nil {
		return nil, err
	}

	return FromSecret(secret)
}

//FromSecret 根据井通私钥创建钱包
func FromSecret(secret string) (*Wallet, error) {
	if secret == "" {
		return nil, constant.ERR_EMPTY_PARAM
	}
	keyPair := &secp256k1.Secp256KeyPair{}
	priv, err := keyPair.DeriveKeyPair(secret)
	if nil != err {
		return nil, err
	}
	wallet := new(Wallet)
	wallet.priv = priv
	wallet.secret = secret
	return wallet, nil
}

//GetPublicKey 获取16进制公钥
func (wallet *Wallet) GetPublicKey() string {
	return wallet.priv.PublicKey.BytesToHex()
}

//GetSecret 获取私钥
func (wallet *Wallet) GetSecret() string {
	return wallet.secret
}

//GetAddress 获取钱包地址
func (wallet *Wallet) GetAddress() string {
	return wallet.priv.PublicKey.ToAddress()
}

//signTx 对交易数据签名
func (wallet *Wallet) signTx(hash []byte) (string, error) {
	priv := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: btcec.S256(),
			X:     wallet.priv.X,
			Y:     wallet.priv.Y,
		},
		D: wallet.priv.D,
	}

	signature, err := (*btcec.PrivateKey)(priv).Sign(hash)
	if err != nil {
		return "", nil
	}

	return fmt.Sprintf("%X", signature.Serialize()), nil
}
