package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

//Walletの情報
type Wallet struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
	address    string
}

//Wallet作成
func NewWallet() *Wallet {
	w := new(Wallet)
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	w.privateKey = privateKey
	w.publicKey = &w.privateKey.PublicKey

	//bitcoinと同じアドレス生成手順
	h2 := sha256.New()
	h2.Write(w.publicKey.X.Bytes())
	h2.Write(w.publicKey.Y.Bytes())
	result2 := h2.Sum(nil)

	h3 := ripemd160.New()
	h3.Write(result2)
	result3 := h3.Sum(nil)

	vd4 := make([]byte, 21)
	vd4[0] = 0x00
	copy(vd4[1:], result3[:])

	h5 := sha256.New()
	h5.Write(vd4)
	result5 := h5.Sum(nil)

	h6 := sha256.New()
	h6.Write(result5)
	result6 := h6.Sum(nil)

	chsum := result6[:4]

	dc8 := make([]byte, 25)
	copy(dc8[:21], vd4[:])
	copy(dc8[21:], chsum[:])

	address := base58.Encode(dc8)
	w.address = address

	return w
}

//アドレスを返すメソッド
func (w *Wallet) Address() string {
	return w.address
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
