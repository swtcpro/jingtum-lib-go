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
      "math/big"

      "github.com/shengdoushi/base58"
)

var (
    ACCOUNT_PREFIX uint8 = 0
    ALPHABET = "jpshnaf39wBUDNEGHJKLM4PQRST7VWXYZ2bcdeCg65rkm8oFqi1tuvAxyz"
    SEED_PREFIX uint8 = 33
)

func sha256Util(sbytes []byte) ([]byte) {
    h := sha256.New()
    h.Write(sbytes)
    return h.Sum(nil)
}

/**
 * concat an item and a buffer
 * @param {integer} item1, should be an integer
 * @param {buffer} buf2, a buffer
 * @returns {buffer} new Buffer
 */
func bufCat0 (item1 uint8, buf2 []byte) ([]byte) {
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
func __encode(version uint8, bytes []byte) (string) {
	buffer := bufCat0(version, bytes)
	checksum := sha256Util(sha256Util(buffer))[0:4]
	ret := bufCat1(buffer, checksum);
    myAlphabet := base58.NewAlphabet(ALPHABET)
    encodedString := base58.Encode(ret, myAlphabet)
	return encodedString
}

func __decode(version uint8, input string) (decodedBytes []byte, err error) {
    myAlphabet := base58.NewAlphabet(ALPHABET)
    decodedBytes, err = base58.Decode(input, myAlphabet)
    if (err != nil || decodedBytes[0] != version || len(decodedBytes) < 5) {
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

func derivePrivateKey(seed []byte) *big.Int {
  order := secp256k1.N
  fmt.Println("start private gen ..")
  privateGen := ScalarMultiple(seed)
  fmt.Println("private Gen ",privateGen)
  Q := secp256k1.ScalarBaseMult(privateGen)
  publickGen := compression(Q)
   fmt.Println("public Gen ",privateGen)
  pb := ScalarMultiple2(publickGen, 0)
  return pb.Add(pb, privateGen).Mod(pb,order)
}

func deriveKeyPair(secret string) (error,*big.Int,*big.Int) {
    myAlphabet := base58.NewAlphabet(ALPHABET)
    decodedBytes, err := base58.Decode(secret, myAlphabet)
    if (err != nil || decodedBytes[0] != SEED_PREFIX || len(decodedBytes) < 5) {
        err := errors.New("invalid input size")
		return err,nil,nil
	}

	entropy := decodedBytes[1:len(decodedBytes) - 4]
    privateKey := derivePrivateKey(entropy)
    fmt.Println(privateKey)

    Q := secp256k1.ScalarBaseMult(privateKey)
    publicKey := new(big.Int).SetBytes(compression(Q))

    return nil,privateKey,publicKey
}

func address(pub *big.Int) string {
    return ToAddress2(pub)
}

func CheckAddress(address string) bool {
    _, err := __decode(ACCOUNT_PREFIX, address)

    if err != nil {
        fmt.Println(err)
        return false
    }

    return true
}

//func GenerateSeed () {
	//var randBytes = brorand(16);
	//return __encode(SEED_PREFIX, randBytes);
//}
