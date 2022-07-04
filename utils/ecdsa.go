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
	return &ecdsa.PublicKey{elliptic.P256(), &bix, &biy}
}

func StringToPrivateKey(s string, pubKey *ecdsa.PublicKey) *ecdsa.PrivateKey {
	b, _ := hex.DecodeString(s)
	var bi big.Int
	_ = bi.SetBytes(b)
	return &ecdsa.PrivateKey{*pubKey, &bi}
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
