package main

import (
	"fmt"
	"log"
	"strings"
	"time"
)

//Block
type Block struct {
	nonce        int
	previousHash string
	timestamp    int64
	//複数のトランザクション
	transactions []string
}

//Blockのプリント用メソッド
func (b *Block) Print() {
	fmt.Printf("timestamp    : %d\n", b.timestamp)
	fmt.Printf("noce         : %d\n", b.nonce)
	fmt.Printf("previousHash : %s\n", b.previousHash)
	fmt.Printf("transactions : %s\n", b.transactions)
}

//新規Block作成
func NewBlock(nonce int, previousHash string) *Block {
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
	bc := new(BlockChain)
	bc.AddBlock(0, "initHash")
	return bc
}

//BlockをChainに追加するメソッド
func (bc *BlockChain) AddBlock(nonce int, previousHash string) *Block {
	b := NewBlock(nonce, previousHash)
	bc.chain = append(bc.chain, b)
	return b
}

//BlockChainのプリント用メソッド
func (bc *BlockChain) Print() {
	for i, block := range bc.chain {
		fmt.Printf("%s Chain %d %s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 25))
}

func init() {
	log.SetPrefix("BlockChain: ")
}

func main() {
	blockChain := NewBlockChain()
	blockChain.AddBlock(1, "test 1")
	blockChain.AddBlock(2, "test 2")
	blockChain.Print()
}
