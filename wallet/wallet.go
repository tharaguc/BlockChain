package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"gobc/utils"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

//Walletの情報
type Wallet struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
	address    string
}

func (w *Wallet) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Privatekey string `json:"private_key"`
		PublicKey  string `json:"public_key"`
		Adddress   string `json:"address"`
	}{
		Privatekey: w.PrivateKeyStr(),
		PublicKey:  w.PublicKeyStr(),
		Adddress:   w.Address(),
	})
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
	return fmt.Sprintf("%064x%064x", w.publicKey.X.Bytes(), w.publicKey.Y.Bytes())
}

//walletからのtransaction情報
type Transaction struct {
	senderPrivateKey *ecdsa.PrivateKey
	senderPublicKey  *ecdsa.PublicKey
	senderAddress    string
	recipientAddress string
	value            float32
}

//transactionを作成するメソッド
func NewTransaction(priKey *ecdsa.PrivateKey, pubKey *ecdsa.PublicKey, sender string, recipient string, value float32) *Transaction {
	return &Transaction{priKey, pubKey, sender, recipient, value}
}

//Signature生成メソッド
func (t *Transaction) GenSignature() *utils.Signature {
	m, _ := json.Marshal(t)
	h := sha256.Sum256([]byte(m))
	r, s, _ := ecdsa.Sign(rand.Reader, t.senderPrivateKey, h[:])
	return &utils.Signature{R: r, S: s}
}

//marshalメソッドカスタム
func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"sender_address"`
		Recipient string  `json:"recipient_address"`
		Value     float32 `json:"value"`
	}{
		Sender:    t.senderAddress,
		Recipient: t.recipientAddress,
		Value:     t.value,
	})
}

//requestの情報
type TransactionRequest struct {
	SenderPrivateKey *string `json:"sender_private_key"`
	SenderPublicKey  *string `json:"sender_public_key"`
	SenderAddress    *string `json:"sender_address"`
	RecipientAddress *string `json:"recipient_address"`
	Value            *string `json:"value"`
}

//requestのValidate
func (req *TransactionRequest) Validate() bool {
	if *req.SenderPrivateKey == "" ||
		*req.SenderPublicKey == "" ||
		*req.SenderAddress == "" ||
		*req.RecipientAddress == "" ||
		*req.Value == "" {
		return false
	}
	return true
}
