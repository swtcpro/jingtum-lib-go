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

package jingtumBaseLib

import (
     "testing"
     "os"
     "flag"
     _"math/big"
)

func Test_sha256Util(t *testing.T) {
    s2 := ([]byte)("ddddddd")
    sbyte := sha256Util(s2)
    t.Log(sbyte)
}

func Test_CheckAddress(t *testing.T) {
	s2 := "jPFikqnwT44sNDaYa32MX4gNRcXbSxnQJe"
    ok := CheckAddress(s2)
    t.Log(ok)
}

func Test_encode(t *testing.T) {
   address := []byte{250,95,217,244,150,117,99,213,201,175,202,133,239,51,28,120,142,54,36,56}
   adds := __encode(ACCOUNT_PREFIX, address)
   t.Log(adds)
}

func Test_deriveKeyPair(t *testing.T) {
    _,pri,pub := deriveKeyPair("shvLwmy5oLuFUuSss7L2PTTh7513J")
    t.Log("private key : ",pri)
    t.Log("public key : ",pub)
    t.Log("public address : ",address(pub))
}

func TestMain(m *testing.M) {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "/tmp")
	flag.Set("v", "3")
	flag.Parse()

	ret := m.Run()
	os.Exit(ret)
}

