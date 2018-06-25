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
)

func Test_sha256Util(t *testing.T) {
    s2 := ([]byte)("ddddddd")
    sbyte := sha256Util(s2)
    t.Log(sbyte)
}

func Test_CheckAddress(t *testing.T) {
	s2 := "ddddddd"
    ok := CheckAddress(s2)
    t.Log(ok)
}

func TestMain(m *testing.M) {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "/tmp")
	flag.Set("v", "3")
	flag.Parse()

	ret := m.Run()
	os.Exit(ret)
}

