package main

import (
	"fmt"
	"html/template"
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

func (wsv *WalletServer) Run() {
	http.HandleFunc("/", wsv.Index)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(wsv.Port())), nil))
}
