package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

//Block
type Block struct {
	nonce        int
	previousHash [32]byte
	timestamp    int64
	//複数のトランザクション
	transactions []string
}

//Blockのプリント用メソッド
func (b *Block) Print() {
	fmt.Printf("timestamp    : %d\n", b.timestamp)
	fmt.Printf("noce         : %d\n", b.nonce)
	fmt.Printf("previousHash : %x\n", b.previousHash)
	fmt.Printf("transactions : %s\n", b.transactions)
}

//BlockのHash化
func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b)
	return sha256.Sum256(m)
}

//適切にJSONMarshalするメソッド（json.Marshalの上書き）
func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp    int64    `json:"timestamp"`
		Nonce        int      `json:"nonce"`
		PreviousHash [32]byte `json:"previous_hash"`
		Transactions []string `json:"transactions"`
	}{
		Timestamp:    b.timestamp,
		Nonce:        b.nonce,
		PreviousHash: b.previousHash,
		Transactions: b.transactions,
	})
}

//新規Block作成
func NewBlock(nonce int, previousHash [32]byte) *Block {
	b := new(Block)
	b.timestamp = time.Now().UnixNano()
	b.nonce = nonce
	b.previousHash = previousHash
	return b
}

//BlockChain
type BlockChain struct {
	transactionPool []string
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
	b := NewBlock(nonce, previousHash)
	bc.chain = append(bc.chain, b)
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
			fmt.Printf("%s Block %d %s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		}
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 25))
}

func init() {
	log.SetPrefix("BlockChain: ")
}

func main() {
	blockChain := NewBlockChain()

	preHash := blockChain.LastBlock().Hash()
	blockChain.AddBlock(1, preHash)
	preHash = blockChain.LastBlock().Hash()
	blockChain.AddBlock(2, preHash)
	blockChain.Print()
}
