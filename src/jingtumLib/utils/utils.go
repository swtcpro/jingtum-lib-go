/**
 * 通用工具类.
 *
 * @FileName: utils.go
 * @Auther : 杨雪波
 * @Email : yangxuebo@yeah.net
 * @CreateTime: 2018-07-26 10:44:32
 * @UpdateTime: 2018-07-26 10:44:54
 */
package utils

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	jtConst "jingtumLib/constant"
	jtEncode "jingtumLib/encoding"
	"math/big"
	"regexp"
)

const REGEX_CURRENCY = "^([a-zA-Z0-9]{3,6}|[A-F0-9]{40})$"

/**
 * concat an item and a buffer
 * @param {integer} item1, should be an integer
 * @param {buffer} buf2, a buffer
 * @returns {buffer} new Buffer
 */
func bufCat0(item1 uint8, buf2 []byte) []byte {
	var buf []byte
	buf = append(buf, item1)
	buf = append(buf, buf2...)
	return buf
}

/**
 * concat one buffer and another
 * @param {buffer} item1, should be an integer
 * @param {buffer} buf2, a buffer
 * @returns {buffer} new Buffer
 */
func bufCat1(buf1 []byte, buf2 []byte) []byte {
	var buf []byte
	buf = append(buf, buf1...)
	buf = append(buf, buf2...)
	return buf
}

/**
 * encode use jingtum base58 encoding
 * including version + data + checksum
 * @param {integer} version
 * @param {buffer} bytes
 * @returns {string}
 * @private
 */
func EncodeB58(version uint8, bytes []byte) string {
	buffer := bufCat0(version, bytes)
	checksum := Sha256Util(Sha256Util(buffer))[0:4]
	ret := bufCat1(buffer, checksum)
	encodedString := jtEncode.Base58Encode(ret, jtEncode.JingTumAlphabet)
	return encodedString
}

func DecodeB58(version uint8, input string) (decodedBytes []byte, err error) {
	decodedBytes, err = jtEncode.Base58Decode(input, jtEncode.JingTumAlphabet)
	if err != nil || decodedBytes[0] != version || len(decodedBytes) < 5 {
		err = errors.New("invalid input size")
		return
	}

	computed := Sha256Util(Sha256Util(decodedBytes[0 : len(decodedBytes)-4]))[0:4]
	checksum := decodedBytes[len(decodedBytes)-4:]

	for i := 0; i != 4; i++ {
		if computed[i] != checksum[i] {
			err = errors.New("invalid checksum")
			return
		}
	}

	decodedBytes = decodedBytes[1 : len(decodedBytes)-4]

	return
}

func BytesToBigInt(b []byte) *big.Int {
	b_buf := bytes.NewBuffer(b)
	var x big.Int
	binary.Read(b_buf, binary.BigEndian, &x)
	return &x
}

func IsValidCurrency(currency string) bool {
	ok, _ := regexp.MatchString(REGEX_CURRENCY, currency)
	return ok
}

func DecodeAddress(address string) []byte {
	decodedBytes, err := DecodeB58(jtConst.ACCOUNT_PREFIX, address)
	if err != nil {
		panic(fmt.Sprintf("Issuer invalid issuer info %v", address))
	}

	return decodedBytes
}

func IsValidAddress(address string) bool {
	_, err := DecodeB58(jtConst.ACCOUNT_PREFIX, address)
	if err != nil {
		return false
	}
	return true
}
