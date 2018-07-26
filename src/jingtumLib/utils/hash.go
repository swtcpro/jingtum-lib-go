/**
 * Hash 工具类.
 *
 * @FileName: hash.go
 * @Auther : 杨雪波
 * @Email : yangxuebo@yeah.net
 * @CreateTime: 2018-07-26 10:44:32
 * @UpdateTime: 2018-07-26 10:44:54
 */
package utils

import (
	"crypto/sha256"
)

func Sha256Util(sbytes []byte) []byte {
	h := sha256.New()
	h.Write(sbytes)
	return h.Sum(nil)
}
