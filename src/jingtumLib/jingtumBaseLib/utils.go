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
    "bytes"
    "encoding/binary"
    "math/big"
    "crypto/sha256"
)


func BytesToBigInt(b []byte) *big.Int {
    b_buf  :=  bytes.NewBuffer(b)
    var x big.Int 
    binary.Read(b_buf, binary.BigEndian, &x)
    return &x
}

func Sha256Util(sbytes []byte) ([]byte) {
    h := sha256.New()
    h.Write(sbytes)
    return h.Sum(nil)
}