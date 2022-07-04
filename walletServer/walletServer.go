package main

import (
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
		io.WriteString(w, string(utils.JsonStatus("success")))
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
