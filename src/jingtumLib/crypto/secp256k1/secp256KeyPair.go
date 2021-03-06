/**
 *
 * 文件功能介绍
 *
 * @FileName: secp256k1.go
 * @Auther : 杨雪波
 * @Email : yangxuebo@yeah.net
 * @CreateTime: 2018-07-04 10:44:32
 * @UpdateTime: 2018-07-04 10:44:54
 * Copyright@2018 版权所有
 */

package secp256k1

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
	"strings"

	jtConst "jingtumlib/constant"
	jtEncode "jingtumlib/encoding"
	jtUtils "jingtumlib/utils"

	"github.com/btcsuite/btcd/btcec"
	"golang.org/x/crypto/ripemd160"
)

/******************************************************************************/
/* ECDSA Keypair Generation */
/******************************************************************************/

var (
	ec EllipticCurve
)

/**
 *  初始化椭圆曲线参数
 */
func init() {
	ec.P, _ = new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F", 16)
	ec.A, _ = new(big.Int).SetString("0000000000000000000000000000000000000000000000000000000000000000", 16)
	ec.B, _ = new(big.Int).SetString("0000000000000000000000000000000000000000000000000000000000000007", 16)
	ec.G.X, _ = new(big.Int).SetString("79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798", 16)
	ec.G.Y, _ = new(big.Int).SetString("483ADA7726A3C4655DA4FBFC0E1108A8FD17B448A68554199C47D08FFB10D4B8", 16)
	ec.N, _ = new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 16)
	ec.H, _ = new(big.Int).SetString("01", 16)
}

//Secp256KeyPair Secp256KeyPair
type Secp256KeyPair struct{}

// PublicKey represents a Bitcoin public key.
type PublicKey struct {
	Point
}

// PrivateKey represents a Bitcoin private key.
type PrivateKey struct {
	PublicKey
	D *big.Int
}

/**
 *  根据 seed byte 获得私钥
 */
func derivePrivateKey(seed []byte) *big.Int {
	privateGen := scalarMultiple(seed)
	publickGen := ec.ScalarBaseMult(privateGen).Compression()
	pb := scalarMultipleDiscrim(publickGen, 0)
	return addMod(pb, privateGen, ec.N)
}

//DeriveKeyPair 根据私钥生成秘钥对
func (*Secp256KeyPair) DeriveKeyPair(secret string) (*PrivateKey, error) {
	decodedBytes, err := jtEncode.Base58Decode(secret, jtEncode.JingTumAlphabet)
	if err != nil || decodedBytes[0] != jtConst.SeedPrefix || len(decodedBytes) < 5 {
		err = fmt.Errorf("invalid input size")
		return nil, err
	}
	var priv PrivateKey
	entropy := decodedBytes[1 : len(decodedBytes)-4]
	priv.D = derivePrivateKey(entropy)
	Q := ec.ScalarBaseMult(priv.D)
	priv.X = Q.X
	priv.Y = Q.Y
	return &priv, nil
}

//GenerateSeed 生成私钥
func (*Secp256KeyPair) GenerateSeed() (string, error) {
	seedBytes := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, seedBytes)
	if err != nil {
		return "", fmt.Errorf("Reading random reader: %s", err.Error())
	}
	return jtUtils.EncodeB58(jtConst.SeedPrefix, seedBytes), nil
}

//CheckAddress 验证地址
func (*Secp256KeyPair) CheckAddress(address string) bool {
	_, err := jtUtils.DecodeB58(jtConst.AccountPrefix, address)

	if err != nil {
		return false
	}

	return true
}

//ToBytes 将椭圆点压缩成(02+X 如Y 偶), 或(03+X 如Y奇),得到 33 字节的 public key
func (pub *PublicKey) ToBytes() (b []byte) {
	x := pub.X.Bytes()

	paddedx := append(bytes.Repeat([]byte{0x00}, 32-len(x)), x...)

	if pub.Y.Bit(0) == 0 {
		return append([]byte{0x02}, paddedx...)
	}

	return append([]byte{0x03}, paddedx...)
}

//ToBytes 私钥转成32字节
func (priv *PrivateKey) ToBytes() (b []byte) {
	d := priv.D.Bytes()

	/* Pad D to 32 bytes */
	paddedD := append(bytes.Repeat([]byte{0x00}, 32-len(d)), d...)

	return paddedD
}

//BytesToHex BytesToHex
func (pub *PublicKey) BytesToHex() string {
	return strings.ToUpper(hex.EncodeToString(pub.ToBytes()))
}

//ToAddress 公钥转成钱包地址
func (pub *PublicKey) ToAddress() (address string) {
	pubBytes := pub.ToBytes()

	/* SHA256 Hash */
	sha256H := sha256.New()
	sha256H.Reset()
	sha256H.Write(pubBytes)
	pubHash1 := sha256H.Sum(nil)

	/* RIPEMD-160 Hash */
	ripemd160H := ripemd160.New()
	ripemd160H.Reset()
	ripemd160H.Write(pubHash1)
	pubHash2 := ripemd160H.Sum(nil)
	address = jtUtils.EncodeB58(jtConst.AccountPrefix, pubHash2)

	return address
}

func scalarMultipleDiscrim(bytes []byte, discrim uint32) *big.Int {
	var privateGen *big.Int
	var i uint32
	for i = 0; i <= 0xFFFFFFFF; i++ {
		// We hash the bytes to find a 256 bit number, looping until we are sure it
		// is less than the order of the curve.
		sh512 := jtUtils.NewSha512()
		sh512.Add(bytes)
		// If the optional discriminator index was passed in, update the hash.
		sh512.Add32(discrim)
		sh512.Add32(i)
		privateGenBytes := sh512.Finish256()
		privateGen = new(big.Int).SetBytes(privateGenBytes)
		if privateGen.Cmp(big.NewInt(0)) == 1 && privateGen.Cmp(ec.N) == -1 {
			return privateGen
		}
	}

	return privateGen
}

func scalarMultiple(bytes []byte) *big.Int {
	var privateGen *big.Int
	var i uint32
	for i = 0; i <= 0xFFFFFFFF; i++ {
		// We hash the bytes to find a 256 bit number, looping until we are sure it
		// is less than the order of the curve.
		sh512 := jtUtils.NewSha512()
		sh512.Add(bytes)
		sh512.Add32(i)
		privateGenBytes := sh512.Finish256()
		privateGen = new(big.Int).SetBytes(privateGenBytes)
		if privateGen.Cmp(big.NewInt(0)) == 1 && privateGen.Cmp(ec.N) == -1 {
			return privateGen
		}
	}
	return privateGen
}

//PrivKeyFromBytes PrivKeyFromBytes
func PrivKeyFromBytes(curve elliptic.Curve, secret string) (*btcec.PrivateKey,
	*btcec.PublicKey) {
	keyPair := &Secp256KeyPair{}
	pri, _ := keyPair.DeriveKeyPair(secret)
	priv := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: curve,
			X:     pri.X,
			Y:     pri.Y,
		},
		D: pri.D,
	}

	return (*btcec.PrivateKey)(priv), (*btcec.PublicKey)(&priv.PublicKey)
}
