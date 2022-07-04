package main

import (
	"encoding/json"
	"fmt"
	"gobc/block"
	"gobc/definition"
	"gobc/utils"
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
		w.Header().Add(definition.CONTENT_TYPE, definition.APP_JSON)
		bc := sv.GetBlockChain()
		m, _ := bc.MarshalJSON()
		io.WriteString(w, string(m[:]))
	default:
		log.Println("Error: Invalid Method")
	}
}

func (sv *Server) Transactions(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		w.Header().Add(definition.CONTENT_TYPE, definition.APP_JSON)
		bc := sv.GetBlockChain()
		transactions := bc.TransactionPool()
		m, _ := json.Marshal(struct {
			Transactions []*block.Transaction `json:"transactions"`
			Length       int                  `json:"length"`
		}{
			Transactions: transactions,
			Length:       len(transactions),
		})
		io.WriteString(w, string(m[:]))

	case http.MethodPost:
		dec := json.NewDecoder(req.Body)
		var t block.TransactionRequest
		err := dec.Decode(&t)
		if err != nil {
			log.Printf("Error: %v", err)
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}
		if !t.Validate() {
			log.Println("Error: missing fields")
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}

		pubKey := utils.StringToPublicKey(*t.SenderPublicKey)
		signature := utils.StringToSignature(*t.Signature)
		bc := sv.GetBlockChain()
		isCreated := bc.CreateTransaction(*t.SenderAddress, *t.RecipientAddress, *t.Value, pubKey, signature)

		w.Header().Add(definition.CONTENT_TYPE, definition.APP_JSON)
		var msg []byte
		if !isCreated {
			w.WriteHeader(http.StatusBadRequest)
			msg = utils.JsonStatus("fail")
		} else {
			w.WriteHeader(http.StatusCreated)
			msg = utils.JsonStatus("success")
		}
		io.WriteString(w, string(msg))

	default:
		log.Println("Error: Invalid Method")
	}
}

func (sv *Server) Run() {
	http.HandleFunc("/", sv.GetChain)
	http.HandleFunc("/transactions", sv.Transactions)
	fmt.Printf("Blockchain Server started on PORT: %v\n", sv.Port())
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(sv.Port())), nil))
}
