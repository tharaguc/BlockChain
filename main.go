package main

import (
	"fmt"
	"gobc/block"
	"gobc/wallet"
	"log"
)

func init() {
	log.SetPrefix("NETWORK: ")
}

func main() {
	wm := wallet.NewWallet()
	wa := wallet.NewWallet()
	wb := wallet.NewWallet()

	t := wallet.NewTransaction(wa.PrivateKey(), wa.PublicKey(), wa.Address(), wb.Address(), 1.01)

	bc := block.NewBlockChain(wm.Address())
	isAdded := bc.AddTransaction(wa.Address(), wb.Address(), 1.01, wa.PublicKey(), t.GenSignature())
	fmt.Println("Added? ", isAdded)

	bc.Mining()
	bc.Print()

	fmt.Printf("A : %.2f\n", bc.CalculateTotalAmount(wa.Address()))
	fmt.Printf("B : %.2f\n", bc.CalculateTotalAmount(wb.Address()))
	fmt.Printf("W : %.2f\n", bc.CalculateTotalAmount(wm.Address()))
}
