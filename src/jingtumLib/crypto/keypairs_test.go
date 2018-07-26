/**
 *
 * 加解密测试类
 *
 * @FileName: keypairs_test.go
 * @Auther : 杨雪波
 * @Email : yangxuebo@yeah.net
 * @CreateTime: 2013-09-16 10:44:32
 * @UpdateTime: 2013-09-16 10:44:54
 * Copyright@2013 版权所有
 */

package crypto

import (
	"flag"
	"jingtumLib/constant"
	"jingtumLib/crypto/secp256k1"
	"jingtumLib/utils"
	"math/big"
	"os"
	"testing"
)

var (
	keyPair KeyPair = &secp256k1.Secp256KeyPair{}
)

func Test_sha256Util(t *testing.T) {
	s2 := ([]byte)("ddddddd")
	sbyte := utils.Sha256Util(s2)
	t.Log(sbyte)
}

func Test_CheckAddress(t *testing.T) {
	s2 := "jPFikqnwT44sNDaYa32MX4gNRcXbSxnQJe"
	ok := keyPair.CheckAddress(s2)
	t.Log(ok)
}

func Test_encode(t *testing.T) {
	address := []byte{250, 95, 217, 244, 150, 117, 99, 213, 201, 175, 202, 133, 239, 51, 28, 120, 142, 54, 36, 56}
	adds := utils.EncodeB58(constant.ACCOUNT_PREFIX, address)
	t.Log(adds)
}

func Test_deriveKeyPair(t *testing.T) {
	pri, _ := keyPair.DeriveKeyPair("snsYqv2FsYLuibE9TGHdG5x5V5Qcn")
	t.Log("private key : ", pri.D)
	t.Log("public key : ", new(big.Int).SetBytes(pri.PublicKey.ToBytes()))
	t.Log("public address : ", pri.PublicKey.ToAddress())
	t.Log("public key to hex : ", pri.PublicKey.BytesToHex())
}

func TestMain(m *testing.M) {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "/tmp")
	flag.Set("v", "3")
	flag.Parse()

	ret := m.Run()
	os.Exit(ret)
}
