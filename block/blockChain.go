package block

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"gobc/utils"
	"log"
	"strings"
	"sync"
	"time"
)

const (
	MINING_DIFFICULTY  = 3
	MINING_SENDER      = "NETWORK"
	MINER_ADDRESS      = "Miner"
	MINING_REWARD      = 1.00
	MINIG_INTERVAL_SEC = 10

	PORT_RANGE_START       = 3000
	PORT_RANGE_END         = 3005
	IP_RANGE_START         = 0
	IP_RANGE_END           = 1
	NEIGHBOR_SYNC_TIME_SEC = 20
)

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
	mutexMinig      sync.Mutex

	neighbors    []string
	mutexNeibors sync.Mutex
}

//chainのMarshal
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

//他のノードを取得するメソッド
func (bc *BlockChain) SetNeighbors() {
	bc.neighbors = utils.FindNeighbors(utils.GetHost(), bc.port, IP_RANGE_START, IP_RANGE_END, PORT_RANGE_START, PORT_RANGE_END)
	fmt.Println(bc.neighbors)
}

func (bc *BlockChain) SyncNeighbors() {
	bc.mutexNeibors.Lock()
	defer bc.mutexNeibors.Unlock()
	bc.SetNeighbors()
}

func (bc *BlockChain) StartSyncNeighbors() {
	bc.SetNeighbors()
	_ = time.AfterFunc(time.Second*NEIGHBOR_SYNC_TIME_SEC, bc.StartSyncNeighbors)
}

func (bc *BlockChain) Run() {
	bc.StartSyncNeighbors()
}

//マイニングメソッド
func (bc *BlockChain) Mining() bool {
	bc.mutexMinig.Lock() //必要？ -> defficulty上がって時間かかる場合
	defer bc.mutexMinig.Unlock()

	if len(bc.transactionPool) == 0 {
		log.Println("Pool is empty")
		return false
	}

	//ネットワークからマイナーへのTransaction追加
	bc.AddTransaction(MINING_SENDER, bc.minerAddress, MINING_REWARD, nil, nil)
	//PoW
	nonce := bc.ProofOfWork()
	preHash := bc.LastBlock().Hash()
	bc.AddBlock(nonce, preHash)
	log.Println("action=mining, status=success")
	return true
}

//intervalごとにminingを開始
func (bc *BlockChain) StartMining() {
	bc.Mining()
	_ = time.AfterFunc(time.Second*MINIG_INTERVAL_SEC, bc.StartMining)
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

type AmountResponse struct {
	Amount float32 `json:"amount"`
}

func (ar *AmountResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Amount float32 `json:"amount"`
	}{
		Amount: ar.Amount,
	})
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

//transactionPoolを返すメソッド
func (bc *BlockChain) TransactionPool() []*Transaction {
	return bc.transactionPool
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

//transactionのSignを認証するメソッド
func (bc *BlockChain) VerifyTransactionSign(senderPubKey *ecdsa.PublicKey, s *utils.Signature, t *Transaction) bool {
	m, _ := json.Marshal(t)
	h := sha256.Sum256([]byte(m))
	return ecdsa.Verify(senderPubKey, h[:], s.R, s.S)
}

//Transactionを追加し他のノードとシンクさせるメソッド
func (bc *BlockChain) CreateTransaction(sender string, recipient string, value float32, senderPubKey *ecdsa.PublicKey, s *utils.Signature) bool {
	isTransacted := bc.AddTransaction(sender, recipient, value, senderPubKey, s)
	//todo
	//sync other nodes
	return isTransacted
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
