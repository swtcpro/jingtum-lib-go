package secp256k1

// func PrivKeyFromBytes(curve elliptic.Curve, secret string) (*btcec.PrivateKey,
// 	*btcec.PublicKey) {
// 	keyPair := &Secp256KeyPair{}
// 	pri, _ := keyPair.DeriveKeyPair(secret)
// 	fmt.Printf("public key : %d\n", new(big.Int).SetBytes(pri.PublicKey.ToBytes()))
// 	fmt.Printf("public address : %s\n", pri.PublicKey.ToAddress())

// 	priv := &ecdsa.PrivateKey{
// 		PublicKey: ecdsa.PublicKey{
// 			Curve: curve,
// 			X:     pri.X,
// 			Y:     pri.Y,
// 		},
// 		D: pri.D,
// 	}

// 	return (*btcec.PrivateKey)(priv), (*btcec.PublicKey)(&priv.PublicKey)
// }
