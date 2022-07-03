package main

import (
	"bc/wallet"
	"fmt"
	"log"
)

func init() {
	log.SetPrefix("BlockChain: ")
}

func main() {
	w := wallet.NewWallet()
	fmt.Println(w.PrivateKeyStr())
	fmt.Println(w.PublicKeyStr())
	// myAddress := MINER_ADDRESS
	// blockChain := NewBlockChain(myAddress)

	// blockChain.AddTransaction("A", "B", 3.0)
	// blockChain.Mining()

	// blockChain.AddTransaction("B", "C", 4.2)
	// blockChain.AddTransaction("C", "A", 3.34)
	// blockChain.Mining()
	// blockChain.Print()

	// fmt.Printf("my %.2f\n", blockChain.CalculateTotalAmount(myAddress))
	// fmt.Printf("A  %.2f\n", blockChain.CalculateTotalAmount("A"))
	// fmt.Printf("B  %.2f\n", blockChain.CalculateTotalAmount("B"))
	// fmt.Printf("C  %.2f\n", blockChain.CalculateTotalAmount("C"))
}
