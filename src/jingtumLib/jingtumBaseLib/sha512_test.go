/**
 *
 * 文件功能介绍
 *
 * @FileName: .go
 * @Auther : Pandao
 * @Email : 272383090@qq.com
 * @CreateTime: 2013-09-16 10:44:32
 * @UpdateTime: 2013-09-16 10:44:54
 * Copyright@2013 版权所有
 */

package jingtumBaseLib

import (
     "testing"
     "encoding/hex"
     _"os"
     _"flag"
)

func Test_sha512Util(t *testing.T) {
   sh512 := NewSha512()
   sh512.Add([]byte("6666我"))
   sh512.Add([]byte("888"))
   sh512.Add32(4294967295)
   t.Log(hex.EncodeToString(sh512.Finish128()))
}
