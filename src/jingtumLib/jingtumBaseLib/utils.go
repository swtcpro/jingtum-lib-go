/**
 *
 * 文件功能介绍
 *
 * @FileName: utils.go
 * @Auther : 杨雪波
 * @Email : yangxuebo@yeah.net
 * @CreateTime: 2018-07-10 10:44:32
 * @UpdateTime: 2018-07-10 10:44:54
 * Copyright@2018 版权所有
 */

package jingtumBaseLib

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"math/big"
	"regexp"
)

func BytesToBigInt(b []byte) *big.Int {
	b_buf := bytes.NewBuffer(b)
	var x big.Int
	binary.Read(b_buf, binary.BigEndian, &x)
	return &x
}

func Sha256Util(sbytes []byte) []byte {
	h := sha256.New()
	h.Write(sbytes)
	return h.Sum(nil)
}

func IsValidCurrency(currency string) bool {
	if currency == "" {
		return false
	}

	match, _ := regexp.MatchString("^([a-zA-Z0-9]{3,6}|[A-F0-9]{40})$", currency)

	return match
}
