package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

//Blockの情報
type Block struct {
	nonce        int
	previousHash [32]byte
	timestamp    int64
	//複数のトランザクション
	transactions []*Transaction
}

//Blockのプリント用メソッド
func (b *Block) Print() {
	fmt.Printf("timestamp    : %d\n", b.timestamp)
	fmt.Printf("noce         : %d\n", b.nonce)
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

//適切にJSONMarshalするメソッドオーバーライド（json.Marshalの上書き）小文字のメンバはmarshalできないがjsonでは小文字で扱いたい
func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp    int64          `json:"timestamp"`
		Nonce        int            `json:"nonce"`
		PreviousHash [32]byte       `json:"previous_hash"`
		Transactions []*Transaction `json:"transactions"`
	}{
		Timestamp:    b.timestamp,
		Nonce:        b.nonce,
		PreviousHash: b.previousHash,
		Transactions: b.transactions,
	})
}

//新規Block作成
func NewBlock(nonce int, previousHash [32]byte, transactions []*Transaction) *Block {
	b := new(Block)
	b.timestamp = time.Now().UnixNano()
	b.nonce = nonce
	b.previousHash = previousHash
	b.transactions = transactions
	return b
}

//BlockChainの情報
type BlockChain struct {
	//Transactionの配列
	transactionPool []*Transaction
	//Blockの配列
	chain []*Block
}

//BlockChainの作成（初期化）
func NewBlockChain() *BlockChain {
	//Genesis Block
	b := &Block{}
	bc := new(BlockChain)
	bc.AddBlock(0, b.Hash())
	return bc
}

//BlockをChainに追加するメソッド
func (bc *BlockChain) AddBlock(nonce int, previousHash [32]byte) *Block {
	b := NewBlock(nonce, previousHash, bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*Transaction{} //Poolを初期化
	return b
}

//最後のBlockを返すメソッド
func (bc *BlockChain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

//BlockChainのプリント用メソッド
func (bc *BlockChain) Print() {
	for i, block := range bc.chain {
		if i == 0 {
			fmt.Printf("%s Genesis Block %s\n", strings.Repeat("=", 25), strings.Repeat("=", 25))
		} else {
			fmt.Printf("\n%s Block %d %s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		}
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 25))
}

//Transactionの情報
type Transaction struct {
	senderAddress    string
	recipientAddress string
	value            float32
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

//適切にJSONMarshalするメソッドオーバーライド（json.Marshalの上書き）小文字のメンバはmarshalできないがjsonでは小文字で扱いたい
func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		SenderAddress    string  `json:"senderAddress"`
		RecipientAddress string  `json:"recipientAddress"`
		Value            float32 `json:"value"`
	}{
		SenderAddress:    t.senderAddress,
		RecipientAddress: t.recipientAddress,
		Value:            t.value,
	})
}

//TransactionをPoolに追加するメソッド
func (bc *BlockChain) AddTransaction(sender string, recipient string, value float32) {
	t := NewTransaction(sender, recipient, value)
	bc.transactionPool = append(bc.transactionPool, t)
}

func init() {
	log.SetPrefix("BlockChain: ")
}

func main() {
	blockChain := NewBlockChain()

	blockChain.AddTransaction("A", "B", 3.0)
	preHash := blockChain.LastBlock().Hash()
	blockChain.AddBlock(1, preHash)

	blockChain.AddTransaction("C", "D", 4.2)
	blockChain.AddTransaction("B", "C", 3.34)
	preHash = blockChain.LastBlock().Hash()
	blockChain.AddBlock(2, preHash)
	blockChain.Print()
}
