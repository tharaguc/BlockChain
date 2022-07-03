package block

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"gobc/utils"
	"log"
	"strings"
	"time"
)

const (
	MINING_DIFFICULTY = 3
	MINING_SENDER     = "NETWORK"
	MINER_ADDRESS     = "Miner"
	MINING_REWARD     = 1.00
)

//Blockの情報
type Block struct {
	timestamp    int64
	nonce        int
	previousHash [32]byte
	transactions []*Transaction
}

//Blockのプリント用メソッド
func (b *Block) Print() {
	fmt.Printf("timestamp    : %d\n", b.timestamp)
	fmt.Printf("nonce        : %d\n", b.nonce)
	fmt.Printf("previousHash : %x\n", b.previousHash)
	for _, t := range b.transactions {
		t.Print()
	}
}

//BlockのHash化
func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b)
	return sha256.Sum256(m)
}

//適切にJSONMarshalするメソッドオーバーライド（json.Marshalの上書き）小文字のフィールドはmarshalできないがjsonでは小文字で扱いたい
func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp    int64          `json:"timestamp"`
		Nonce        int            `json:"nonce"`
		PreviousHash string         `json:"previous_hash"`
		Transactions []*Transaction `json:"transactions"`
	}{
		Timestamp:    b.timestamp,
		Nonce:        b.nonce,
		PreviousHash: fmt.Sprintf("%x", b.previousHash),
		Transactions: b.transactions,
	})
}

//新規Block作成
func NewBlock(nonce int, previousHash [32]byte, transactions []*Transaction) *Block {
	b := new(Block)
	b.timestamp = time.Now().UnixNano()
	b.nonce = nonce
	b.previousHash = previousHash
	b.transactions = transactions
	return b
}

//BlockChainの情報
type BlockChain struct {
	transactionPool []*Transaction
	chain           []*Block
	minerAddress    string
	port            uint16
}

func (bc *BlockChain) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Blocks []*Block `json:"chains"`
	}{
		Blocks: bc.chain,
	})
}

//BlockChainの作成（初期化）
func NewBlockChain(minerAddress string, port uint16) *BlockChain {
	//Genesis Block
	b := &Block{}
	bc := new(BlockChain)
	bc.minerAddress = minerAddress
	bc.AddBlock(0, b.Hash())
	bc.port = port
	return bc
}

//マイニングメソッド
func (bc *BlockChain) Mining() bool {
	//ネットワークからマイナーへのTransaction追加
	bc.AddTransaction(MINING_SENDER, bc.minerAddress, MINING_REWARD, nil, nil)
	nonce := bc.ProofOfWork()
	preHash := bc.LastBlock().Hash()
	bc.AddBlock(nonce, preHash)
	log.Println("action=mining, status=success")
	return true
}

//アドレスをもとにtransactionによる差分を計算
func (bc *BlockChain) CalculateTotalAmount(address string) float32 {
	var total float32 = 0.00

	//全てのtransaction参照
	for _, b := range bc.chain {
		for _, t := range b.transactions {
			value := t.value
			if address == t.recipientAddress {
				total += value
			}
			if address == t.senderAddress {
				total -= value
			}
		}
	}
	return total
}

//BlockをChainに追加するメソッド
func (bc *BlockChain) AddBlock(nonce int, previousHash [32]byte) *Block {
	b := NewBlock(nonce, previousHash, bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*Transaction{} //Poolを初期化
	return b
}

//最後のBlockを返すメソッド
func (bc *BlockChain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

//BlockChainのプリント用メソッド
func (bc *BlockChain) Print() {
	for i, block := range bc.chain {
		if i == 0 {
			fmt.Printf("%s Genesis Block %s\n", strings.Repeat("=", 25), strings.Repeat("=", 25))
		} else {
			fmt.Printf("\n%s Block %d %s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		}
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 25))
}

//Transactionの情報
type Transaction struct {
	senderAddress    string
	recipientAddress string
	value            float32
}

//Transactionを作成するメソッド
func NewTransaction(sender string, recipient string, value float32) *Transaction {
	return &Transaction{sender, recipient, value}
}

//Transaction情報のプリント用メソッド
func (t *Transaction) Print() {
	fmt.Printf("%s Transaction %s\n", strings.Repeat("-", 6), strings.Repeat("-", 6))
	fmt.Printf("senderAdress     : %s\n", t.senderAddress)
	fmt.Printf("recipientAdress  : %s\n", t.recipientAddress)
	fmt.Printf("value            : %.2f\n", t.value)
	fmt.Println(strings.Repeat("-", 25))
}

//適切にJSONMarshalするメソッドオーバーライド（json.Marshalの上書き）小文字のメンバはmarshalできないがjsonでは小文字で扱いたい
func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		SenderAddress    string  `json:"sender_address"`
		RecipientAddress string  `json:"recipient_address"`
		Value            float32 `json:"value"`
	}{
		SenderAddress:    t.senderAddress,
		RecipientAddress: t.recipientAddress,
		Value:            t.value,
	})
}

//transactionのSignを認証するメソッド
func (bc *BlockChain) VerifyTransactionSign(senderPubKey *ecdsa.PublicKey, s *utils.Signature, t *Transaction) bool {
	m, _ := json.Marshal(t)
	h := sha256.Sum256([]byte(m))
	return ecdsa.Verify(senderPubKey, h[:], s.R, s.S)
}

//TransactionをPoolに追加するメソッド
func (bc *BlockChain) AddTransaction(sender string, recipient string, value float32, senderPubKey *ecdsa.PublicKey, s *utils.Signature) bool {
	t := NewTransaction(sender, recipient, value)

	//マイニング報酬の場合
	if sender == MINING_SENDER {
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	}

	if bc.VerifyTransactionSign(senderPubKey, s, t) {

		/*
			if bc.CalculateTotalAmount(sender) < value {
				log.Println("Error: Not enough balance in a wallet")
				return false
			}
		*/

		bc.transactionPool = append(bc.transactionPool, t)
		return true
	} else {
		log.Println("Error: Verify TransactionSign")
	}
	return false
}

//PoolのTransactionsをコピーするメソッド
func (bc *BlockChain) CopyTransactionsFromPool() []*Transaction {
	copy := make([]*Transaction, 0)
	for _, t := range bc.transactionPool {
		copy = append(copy, NewTransaction(t.senderAddress, t.recipientAddress, t.value))
	}
	return copy
}

//nonceが正しいかどうか判定するメソッド
func (bc *BlockChain) IsValidProof(nonce int, preHash [32]byte, transactions []*Transaction, difficulty int) bool {
	zeros := strings.Repeat("0", difficulty)
	guessBlock := Block{0, nonce, preHash, transactions}
	guessBlockHash := fmt.Sprintf("%x", guessBlock.Hash()) //byte -> Base16に変換
	return guessBlockHash[:difficulty] == zeros            //最初の{difficulty}文字判定
}

//正しいnonceを求めるメソッド
func (bc *BlockChain) ProofOfWork() int {
	transactions := bc.CopyTransactionsFromPool()
	preHash := bc.LastBlock().Hash()
	nonce := 0
	//正しいnonceになるまでループ
	for !bc.IsValidProof(nonce, preHash, transactions, MINING_DIFFICULTY) {
		nonce += 1
	}
	return nonce
}
