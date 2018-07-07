/**
 *
 * 文件功能介绍
 *
 * @FileName: base58_test.go
 * @Auther : 杨雪波
 * @Email : yangxuebo@yeah.net
 * @CreateTime: 2013-09-16 10:44:32
 * @UpdateTime: 2013-09-16 10:44:54
 * Copyright@2013 版权所有
 */

package jingtumBaseLib

import (
     "testing"
)

func Test_base58Encode(t *testing.T) {
    t.Log(Base58Encode([]byte("ddddd"),JingTumAlphabet))
}
