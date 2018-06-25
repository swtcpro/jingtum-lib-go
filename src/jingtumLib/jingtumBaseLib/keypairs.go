/**
 *
 * 文件功能介绍
 *
 * @FileName: keypairs.go
 * @Auther : 杨雪波
 * @Email : yangxuebo@yeah.net
 * @CreateTime: 2013-09-16 10:44:32
 * @UpdateTime: 2013-09-16 10:44:54
 * Copyright@2013 版权所有
 */

package jingtumBaseLib

import (
      "crypto/sha256"
      "errors"
      "fmt"

      "github.com/shengdoushi/base58"
)

var (
    ACCOUNT_PREFIX uint8 = 0
    ALPHABET = "jpshnaf39wBUDNEGHJKLM4PQRST7VWXYZ2bcdeCg65rkm8oFqi1tuvAxyz"
)

func sha256Util(sbytes []byte) ([]byte) {
    h := sha256.New()
    h.Write(sbytes)
    return h.Sum(nil)
}

func __decode(version uint8, input string) (decodedBytes []byte, err error) {
    myAlphabet := base58.NewAlphabet(ALPHABET)
    decodedBytes, err = base58.Decode(input, myAlphabet)
    if (err != nil || decodedBytes[0] != ACCOUNT_PREFIX || len(decodedBytes) < 5) {
        err = errors.New("invalid input size")
		return
	}
    
    computed := sha256Util(sha256Util(decodedBytes[0:len(decodedBytes) - 4]))[0:4]
    checksum := decodedBytes[len(decodedBytes) - 4:]

    for i := 0; i != 4; i++ {
        if computed[i] != checksum[i] {
            err = errors.New("invalid checksum")
		    return
        }
    }

    decodedBytes = decodedBytes[1:len(decodedBytes) - 4]
    return
}

func CheckAddress(address string) bool {
    _, err := __decode(ACCOUNT_PREFIX, address)

    if err != nil {
        fmt.Println(err)
        return false
    }

    return true
}
