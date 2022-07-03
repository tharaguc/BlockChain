package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
)

//Walletの情報
type Wallet struct {
	privateKey *ecdsa.PrivateKey
	publicKey *ecdsa.PublicKey
}

//Wallet作成
func NewWallet() *Wallet {
	w := new(Wallet)
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	w.privateKey = privateKey
	w.publicKey = &w.privateKey.PublicKey
	return w
}

//privateKeyを返すメソッド
func (w *Wallet) PrivateKey() *ecdsa.PrivateKey {
	return w.privateKey
}

//privateKeyの文字を返すメソッド
func (w *Wallet) PrivateKeyStr() string {
	return fmt.Sprintf("%x", w.privateKey.D.Bytes())
}

//publicKeyを返すメソッド
func (w *Wallet) PublicKey() *ecdsa.PublicKey {
	return w.publicKey
}

//publicKeyの文字を返すメソッド
func (w *Wallet) PublicKeyStr() string {
	return fmt.Sprintf("%x%x", w.privateKey.X.Bytes(), w.publicKey.Y.Bytes())
}