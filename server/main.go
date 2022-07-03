package main

import (
	"flag"
	"log"
)

func init() {
	log.SetPrefix("NETWORK: ")
}

func main() {
	port := flag.Uint("p", 5001, "TCP Port Number for Server")
	flag.Parse()
	app := NewServer(uint16(*port))
	app.Run()
}
