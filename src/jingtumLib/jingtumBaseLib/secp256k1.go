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
	"golang.org/x/crypto/ripemd160"
	"crypto/sha256"
	"fmt"
	"io"
	"math/big"
	"strings"
    "encoding/binary"
)

/******************************************************************************/
/* ECDSA Keypair Generation */
/******************************************************************************/

var secp256k1 EllipticCurve

func init() {
	/* See Certicom's SEC2 2.7.1, pg.15 */
	/* secp256k1 elliptic curve parameters */
	secp256k1.P, _ = new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F", 16)
	secp256k1.A, _ = new(big.Int).SetString("0000000000000000000000000000000000000000000000000000000000000000", 16)
	secp256k1.B, _ = new(big.Int).SetString("0000000000000000000000000000000000000000000000000000000000000007", 16)
	secp256k1.G.X, _ = new(big.Int).SetString("79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798", 16)
	secp256k1.G.Y, _ = new(big.Int).SetString("483ADA7726A3C4655DA4FBFC0E1108A8FD17B448A68554199C47D08FFB10D4B8", 16)
	secp256k1.N, _ = new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 16)
	secp256k1.H, _ = new(big.Int).SetString("01", 16)
}

// PublicKey represents a Bitcoin public key.
type PublicKey struct {
	Point
}

// PrivateKey represents a Bitcoin private key.
type PrivateKey struct {
	PublicKey
	D *big.Int
}

func ScalarMultiple(bytes []byte, discrim uint32) big.Int {
    var privateGen big.Int
	for i := 0; i <= 0xFFFFFFFF; i++ {
		// We hash the bytes to find a 256 bit number, looping until we are sure it
	    // is less than the order of the curve.
	    sh512 := NewSha512()
        sh512.Add(bytes)
	    // If the optional discriminator index was passed in, update the hash.
        sh512.Add32(discrim)
	    sh512.Add32(i)
        privateGenBytes := sh512.Finish256()
        privateGen = BytesToBigInt(privateGenBytes)
	    if (privateGen > 0 && privateGen < secp256k1.N) {
	      return privateGen
	    }
	}

    return privateGen
}

func ScalarMultiple(bytes []byte) big.Int {
    var privateGen big.Int
	for i := 0; i <= 0xFFFFFFFF; i++ {
		// We hash the bytes to find a 256 bit number, looping until we are sure it
	    // is less than the order of the curve.
	    sh512 := NewSha512()
        sh512.Add(bytes)
	    sh512.Add32(i)
        privateGenBytes := sh512.Finish256()
        privateGen = BytesToBigInt(privateGenBytes)
	    if (privateGen > 0 && privateGen < secp256k1.N) {
            return privateGen
        }
	}
    return privateGen
}
