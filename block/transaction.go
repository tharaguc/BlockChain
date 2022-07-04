package block

import (
	"encoding/json"
	"fmt"
	"strings"
)

//Transactionの情報
type Transaction struct {
	senderAddress    string
	recipientAddress string
	value            float32
}

//適切にJSONMarshalするメソッドオーバーライド（json.Marshalの上書き）小文字のメンバはmarshalできないがjsonでは小文字で扱いたい
func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		SenderAddress    string  `json:"sender_address"`
		RecipientAddress string  `json:"recipient_address"`
		Value            float32 `json:"value"`
	}{
		SenderAddress:    t.senderAddress,
		RecipientAddress: t.recipientAddress,
		Value:            t.value,
	})
}

//Transactionを作成するメソッド
func NewTransaction(sender string, recipient string, value float32) *Transaction {
	return &Transaction{sender, recipient, value}
}

//Transaction情報のプリント用メソッド
func (t *Transaction) Print() {
	fmt.Printf("%s Transaction %s\n", strings.Repeat("-", 6), strings.Repeat("-", 6))
	fmt.Printf("senderAdress     : %s\n", t.senderAddress)
	fmt.Printf("recipientAdress  : %s\n", t.recipientAddress)
	fmt.Printf("value            : %.2f\n", t.value)
	fmt.Println(strings.Repeat("-", 25))
}

//requestの情報
type TransactionRequest struct {
	SenderPrivateKey *string  `json:"sender_private_key"`
	SenderPublicKey  *string  `json:"sender_public_key"`
	SenderAddress    *string  `json:"sender_address"`
	RecipientAddress *string  `json:"recipient_address"`
	Value            *float32 `json:"value"`
	Signature        *string  `json:"signature"`
}

//requestのValidate
func (req *TransactionRequest) Validate() bool {
	if *req.SenderPrivateKey == "" ||
		*req.SenderPublicKey == "" ||
		*req.SenderAddress == "" ||
		*req.RecipientAddress == "" ||
		req.Value == nil {
		return false
	}
	return true
}
