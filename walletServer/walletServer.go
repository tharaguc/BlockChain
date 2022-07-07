package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gobc/block"
	"gobc/definition"
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

//walletを作成
func (wsv *WalletServer) Wallet(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		w.Header().Add(definition.CONTENT_TYPE, definition.APP_JSON)
		newWallet := wallet.NewWallet() //walletの作成
		m, _ := newWallet.MarshalJSON()
		io.WriteString(w, string(m[:]))
	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error: Invalid http method")
	}
}

//clientからのrequestをもとにノードへrequestを送る
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
			return
		}
		value32 := float32(value)

		w.Header().Add(definition.CONTENT_TYPE, definition.APP_JSON)

		transaction := wallet.NewTransaction(priKey, pubKey, *t.SenderAddress, *t.RecipientAddress, value32)
		signature := transaction.GenSignature()
		signStr := signature.String()

		//ノードへのrequest
		req := &block.TransactionRequest{
			SenderPrivateKey: t.SenderPrivateKey,
			SenderPublicKey:  t.SenderPublicKey,
			SenderAddress:    t.SenderAddress,
			RecipientAddress: t.RecipientAddress,
			Value:            &value32,
			Signature:        &signStr,
		}
		m, _ := json.Marshal(req)
		buff := bytes.NewBuffer(m)
		res, _ := http.Post(wsv.Gateway()+"/transactions", definition.APP_JSON, buff)
		if res.StatusCode == 201 {
			io.WriteString(w, string(utils.JsonStatus("success")))
			return
		}
		io.WriteString(w, string(utils.JsonStatus("fail")))

	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error: Invalid http method")
	}
}

func (wsv *WalletServer) WalletAmount(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		address := req.URL.Query().Get("address")
		endpoint := fmt.Sprintf("%s/amount", wsv.Gateway())

		//request情報の作成
		client := &http.Client{}
		req, _ := http.NewRequest("GET", endpoint, nil)
		q := req.URL.Query()
		q.Add("address", address)
		req.URL.RawQuery = q.Encode()

		res, err := client.Do(req)
		if err != nil {
			log.Printf("Error: %v", err)
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}

		w.Header().Add(definition.CONTENT_TYPE, definition.APP_JSON)
		if res.StatusCode == 200 {
			dec := json.NewDecoder(res.Body)
			var amountRes block.AmountResponse
			err := dec.Decode(&amountRes)
			if err != nil {
				log.Printf("Error: %v", err)
				io.WriteString(w, string(utils.JsonStatus("fail")))
				return
			}

			m, _ := json.Marshal(struct {
				Message string  `json:"message"`
				Amount  float32 `json:"amount"`
			}{
				Message: "success",
				Amount:  amountRes.Amount,
			})
			io.WriteString(w, string(m[:]))
		} else {
			io.WriteString(w, string(utils.JsonStatus("fail")))
		}

	default:
		log.Println("Error: Invalid http method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (wsv *WalletServer) Run() {
	http.HandleFunc("/", wsv.Index)
	http.HandleFunc("/wallet", wsv.Wallet)
	http.HandleFunc("/wallet/amount", wsv.WalletAmount)
	http.HandleFunc("/transaction", wsv.CreateTransaction)
	fmt.Printf("Wallet Server started on PORT: %v\n", wsv.Port())
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(wsv.Port())), nil))
}
