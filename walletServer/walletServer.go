package main

import (
	"encoding/json"
	"fmt"
	"gobc/utils"
	"gobc/wallet"
	"html/template"
	"io"
	"log"
	"net/http"
	"path"
	"strconv"
)

const tempDir = "templates"

type WalletServer struct {
	port uint16
	//接続するノード
	gateway string
}

func NewWalletServer(port uint16, gateway string) *WalletServer {
	return &WalletServer{port, gateway}
}

func (wsv *WalletServer) Port() uint16 {
	return wsv.port
}

func (wsv *WalletServer) Gateway() string {
	return wsv.gateway
}

func (wsv *WalletServer) Index(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		t, err := template.ParseFiles(path.Join(tempDir, "index.html"))
		if err != nil {
			fmt.Println("temp err")
		}
		t.Execute(w, "")
	default:
		log.Println("Error: Invalid http method")
	}
}

func (wsv *WalletServer) Wallet(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		w.Header().Add("Content-Type", "application/json")
		newWallet := wallet.NewWallet() //walletの作成
		m, _ := newWallet.MarshalJSON()
		io.WriteString(w, string(m[:]))
	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error: Invalid http method")
	}
}

func (wsv *WalletServer) CreateTransaction(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		dec := json.NewDecoder(req.Body)
		var t wallet.TransactionRequest
		err := dec.Decode(&t)
		if err != nil {
			log.Printf("Error: %v\n", err)
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}
		if !t.Validate() {
			log.Println("Error: missing fields")
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}

		pubKey := utils.StringToPublicKey(*t.SenderPublicKey)
		priKey := utils.StringToPrivateKey(*t.SenderPrivateKey, pubKey)
		value, err := strconv.ParseFloat(*t.Value, 32)
		if err != nil {
			log.Println("Error: Parse error")
			io.WriteString(w, string(utils.JsonStatus("fail")))
		}
		value32 := float32(value)

		io.WriteString(w, string(utils.JsonStatus("success")))
		fmt.Println(pubKey)
		fmt.Println(priKey)
		fmt.Println(value32)

	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error: Invalid http method")
	}
}

func (wsv *WalletServer) Run() {
	http.HandleFunc("/", wsv.Index)
	http.HandleFunc("/wallet", wsv.Wallet)
	http.HandleFunc("/transaction", wsv.CreateTransaction)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(wsv.Port())), nil))
}
