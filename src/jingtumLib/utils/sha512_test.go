/**
 *
 * 文件功能介绍
 *
 * @FileName: sha512_test.go
 * @Auther : 杨雪波
 * @Email : yangxuebo@yeah.net
 * @CreateTime: 2018-07-16 10:44:32
 * @UpdateTime: 2018-07-16 10:44:54
 * Copyright@2018 版权所有
 */

package utils

import (
	"encoding/hex"
	"testing"
)

func Test_sha512Util(t *testing.T) {
	sh512 := NewSha512()
	sh512.Add([]byte("6666我"))
	sh512.Add([]byte("888"))
	sh512.Add32(4294967295)
	t.Log(hex.EncodeToString(sh512.Finish128()))
}
