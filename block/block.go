package block

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

//Blockの情報
type Block struct {
	timestamp    int64
	nonce        int
	previousHash [32]byte
	transactions []*Transaction
}

//適切にJSONMarshalするメソッドオーバーライド（json.Marshalの上書き）小文字のフィールドはmarshalできないがjsonでは小文字で扱いたい
func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp    int64          `json:"timestamp"`
		Nonce        int            `json:"nonce"`
		PreviousHash string         `json:"previous_hash"`
		Transactions []*Transaction `json:"transactions"`
	}{
		Timestamp:    b.timestamp,
		Nonce:        b.nonce,
		PreviousHash: fmt.Sprintf("%x", b.previousHash),
		Transactions: b.transactions,
	})
}

//Blockのプリント用メソッド
func (b *Block) Print() {
	fmt.Printf("timestamp    : %d\n", b.timestamp)
	fmt.Printf("nonce        : %d\n", b.nonce)
	fmt.Printf("previousHash : %x\n", b.previousHash)
	for _, t := range b.transactions {
		t.Print()
	}
}

//BlockのHash化
func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b)
	return sha256.Sum256(m)
}
