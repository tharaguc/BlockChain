package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"math/big"
)

//Signatureの情報
type Signature struct {
	R *big.Int
	S *big.Int
}

func (s *Signature) String() string {
	return fmt.Sprintf("%064x%064x", s.R, s.S)
}

func StringToSignature(s string) *Signature {
	x, y := StringToBigInts(s)
	return &Signature{&x, &y}
}

func StringToPublicKey(s string) *ecdsa.PublicKey {
	bix, biy := StringToBigInts(s)
	res := new(ecdsa.PublicKey)
	res.Curve = elliptic.P256()
	res.X = &bix
	res.Y = &biy
	return res
}

func StringToPrivateKey(s string, pubKey *ecdsa.PublicKey) *ecdsa.PrivateKey {
	b, _ := hex.DecodeString(s)
	var bi big.Int
	_ = bi.SetBytes(b)
	res := new(ecdsa.PrivateKey)
	res.PublicKey = *pubKey
	res.D = &bi
	return res
}

func StringToBigInts(s string) (big.Int, big.Int) {
	bx, _ := hex.DecodeString(s[:64])
	by, _ := hex.DecodeString(s[64:])

	var bix big.Int
	var biy big.Int

	_ = bix.SetBytes(bx)
	_ = biy.SetBytes(by)
	return bix, biy
}
