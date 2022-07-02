package main

import (
	"fmt"
	"log"
	"time"
)

type Block struct {
	nonce        int
	previousHash string
	timestamp    int64

	//複数のトランザクション
	transactions []string
}

//ブロックのプリント用メソッド
func (b *Block) Print() {
	fmt.Printf("timestamp    %d\n", b.timestamp)
	fmt.Printf("noce         %d\n", b.nonce)
	fmt.Printf("previousHash %s\n", b.previousHash)
	fmt.Printf("transactions %s\n", b.transactions)
}

func NewBlock(nonce int, previousHash string) *Block {
	b := new(Block)
	b.timestamp = time.Now().UnixNano()
	b.nonce = nonce
	b.previousHash = previousHash
	return b
}

func init() {
	log.SetPrefix("BlockChain: ")
}

func main() {
	b := NewBlock(0, "testHash")
	b.Print()
}
