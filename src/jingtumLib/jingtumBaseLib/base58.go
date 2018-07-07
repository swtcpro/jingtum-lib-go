/**
 *
 * base58编码
 *
 * @FileName: base58.go
 * @Auther : 杨雪波
 * @Email : yangxuebo@yeah.net
 * @CreateTime: 2018-07-03 10:44:32
 * @UpdateTime: 2018-07-03 10:44:54
 * Copyright@2013 版权所有
 */

package jingtumBaseLib

import (
    "errors"
    "fmt"
)

// Errors
var (
	ErrorInvalidBase58String = errors.New("invalid base58 string")
)

var (    
    //井通字母表
	JingTumAlphabet = NewAlphabet("jpshnaf39wBUDNEGHJKLM4PQRST7VWXYZ2bcdeCg65rkm8oFqi1tuvAxyz")
)

// base58 字母表结构。
type Alphabet struct {
	encodeTable        [58]rune
	decodeTable        [256]int
	unicodeDecodeTable []rune
}

/**
 *  创建base58 字母表结构,字母表必须是58个字符
 *  params:
 *      自定义的字母表字符
 *  return:
 *      *Alphabet
 *
 */
func NewAlphabet(alphabet string) *Alphabet {
	alphabetRunes := []rune(alphabet)
	if len(alphabetRunes) != 58 {
		panic(fmt.Sprintf("Base58 Alphabet length must 58, but %d", len(alphabetRunes)))
	}

	ret := new(Alphabet)
	for i := range ret.decodeTable {
		ret.decodeTable[i] = -1
	}
	ret.unicodeDecodeTable = make([]rune, 0, 58*2)
	for idx, ch := range alphabetRunes {
		ret.encodeTable[idx] = ch
		if ch >= 0 && ch < 256 {
			ret.decodeTable[byte(ch)] = idx
		} else {
			ret.unicodeDecodeTable = append(ret.unicodeDecodeTable, ch)
			ret.unicodeDecodeTable = append(ret.unicodeDecodeTable, rune(idx))
		}
	}
	return ret
}

/**
 *  以传入的字母表对输入值进行base58编码。
 *  params:
 *      input:待编码字节数组
 *      alphabet:base58编码字母表
 *  return:
 *      string
 */
func Base58Encode(input []byte, alphabet *Alphabet) string {
	inputLength := len(input)
	prefixZeroes := 0
	for prefixZeroes < inputLength && input[prefixZeroes] == 0 {
		prefixZeroes++
	}

	capacity := inputLength*138/100 + 1 // log256 / log58
	output := make([]byte, capacity)
	outputReverseEnd := capacity - 1

	for inputPos := prefixZeroes; inputPos < inputLength; inputPos++ {
		carry := uint32(input[inputPos])

		outputIdx := capacity - 1
		for ; carry != 0 || outputIdx > outputReverseEnd; outputIdx-- {
			carry += (uint32(output[outputIdx]) << 8) // XX << 8 same as: 256 * XX
			output[outputIdx] = byte(carry % 58)
			carry /= 58
		}
		outputReverseEnd = outputIdx
	}

	encodeTable := alphabet.encodeTable
	// when not contains unicode, use []byte to improve performance
	if len(alphabet.unicodeDecodeTable) == 0 {
		retStrBytes := make([]byte, prefixZeroes+(capacity-1-outputReverseEnd))
		for i := 0; i < prefixZeroes; i++ {
			retStrBytes[i] = byte(encodeTable[0])
		}
		for i, n := range output[outputReverseEnd+1:] {
			retStrBytes[prefixZeroes+i] = byte(encodeTable[n])
		}
		return string(retStrBytes)
	}
	retStrRunes := make([]rune, prefixZeroes+(capacity-1-outputReverseEnd))
	for i := 0; i < prefixZeroes; i++ {
		retStrRunes[i] = encodeTable[0]
	}
	for i, n := range output[outputReverseEnd+1:] {
		retStrRunes[prefixZeroes+i] = encodeTable[n]
	}
	return string(retStrRunes)
}

/**
 *  以传入的字母表对输入值进行base58解码。
 *  params:
 *      input:待解码字符串
 *      alphabet:base58编码字母表
 *  return:
 *      []byte, error
 */
func Base58Decode(input string, alphabet *Alphabet) ([]byte, error) {
	inputBytes := []rune(input)
	inputLength := len(inputBytes)
	capacity := inputLength*733/1000 + 1 // log(58) / log(256)
	output := make([]byte, capacity)
	outputReverseEnd := capacity - 1

	// prefix 0
	zero58Byte := alphabet.encodeTable[0]
	prefixZeroes := 0
	for prefixZeroes < inputLength && inputBytes[prefixZeroes] == zero58Byte {
		prefixZeroes++
	}

	for inputPos := 0; inputPos < inputLength; inputPos++ {
		carry := -1
		target := inputBytes[inputPos]
		if target >= 0 && target < 256 {
			carry = alphabet.decodeTable[target]
		} else { // unicode
			for i := 0; i < len(alphabet.unicodeDecodeTable); i += 2 {
				if alphabet.unicodeDecodeTable[i] == target {
					carry = int(alphabet.unicodeDecodeTable[i+1])
					break
				}
			}
		}
		if carry == -1 {
			return nil, ErrorInvalidBase58String
		}

		outputIdx := capacity - 1
		for ; carry != 0 || outputIdx > outputReverseEnd; outputIdx-- {
			carry += 58 * int(output[outputIdx])
			output[outputIdx] = byte(uint32(carry) & 0xff) // same as: byte(uint32(carry) % 256)
			carry >>= 8                                    // same as: carry /= 256
		}
		outputReverseEnd = outputIdx
	}

	retBytes := make([]byte, prefixZeroes+(capacity-1-outputReverseEnd))
	for i, n := range output[outputReverseEnd+1:] {
		retBytes[prefixZeroes+i] = n
	}
	return retBytes, nil
}
