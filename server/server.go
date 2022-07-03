package main

import (
	"gobc/block"
	"gobc/wallet"
	"io"
	"log"
	"net/http"
	"strconv"
)

//毎回reqest出さない
var cache map[string]*block.BlockChain = make(map[string]*block.BlockChain)

type Server struct {
	port uint16
}

//create server
func NewServer(port uint16) *Server {
	return &Server{port: port}
}

//return port
func (sv *Server) Port() uint16 {
	return sv.port
}

func (sv *Server) GetBlockChain() *block.BlockChain {
	//キャッシュにあるか確認
	bc, ok := cache["chain"]
	if !ok {
		minerWallet := wallet.NewWallet()
		bc = block.NewBlockChain(minerWallet.Address(), sv.Port())
		cache["chain"] = bc
		log.Printf("priKey  : %v", minerWallet.PrivateKeyStr())
		log.Printf("pubKey  : %v", minerWallet.PublicKeyStr())
		log.Printf("address : %v", minerWallet.Address())
	}
	return bc
}

func (sv *Server) GetChain(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")
		bc := sv.GetBlockChain()
		m, _ := bc.MarshalJSON()
		io.WriteString(w, string(m[:]))
	default:
		log.Println("Error: Invalid Method")
	}
}

func (sv *Server) Run() {
	http.HandleFunc("/", sv.GetChain)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(sv.Port())), nil))
}
