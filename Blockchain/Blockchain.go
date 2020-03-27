package assignment01IBC

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
)

//Transaction ...
type Transaction struct {
	Sender   string
	Receiver string
	Amount   int
}

//HashPointer holds pointer to previos block ans also hash of previous block
type HashPointer struct {
	PreviousBlockHash string
	PreviousBlockPtr  *Block
}

//Block structure represents each block of Blockchain
type Block struct {
	BlockTransactions []Transaction
	HashPtr           HashPointer
	BlockNumber       int
}

//GetHash take an argument as string and returns the SHA256 hash of the argument
func GetHash(data string) string {
	shaHash := sha256.New()
	shaHash.Write([]byte(data))
	return hex.EncodeToString(shaHash.Sum(nil))
}

func getAllTransactions(chainHead *Block) string {
	var allTrans string
	for _, transaction := range chainHead.BlockTransactions {
		allTrans += transaction.Sender + transaction.Receiver + strconv.Itoa(transaction.Amount)
	}
	return allTrans
}

//InsertBlock function allows the miner to to insert the block
func InsertBlock(transaction []Transaction, chainHead *Block) *Block {
	var blockChain *Block
	if chainHead == nil {
		newHashPtr := HashPointer{PreviousBlockHash: "0000000000000000000000000000000000000000000000000000000000000000", PreviousBlockPtr: nil}
		genesisBlock := Block{BlockTransactions: transaction, HashPtr: newHashPtr, BlockNumber: 1}
		blockChain = &genesisBlock
	} else {
		data := getAllTransactions(chainHead)
		data += chainHead.HashPtr.PreviousBlockHash + strconv.Itoa(chainHead.BlockNumber)
		newHashPtr := HashPointer{PreviousBlockHash: GetHash(data), PreviousBlockPtr: chainHead}
		newBlock := Block{BlockTransactions: transaction, HashPtr: newHashPtr, BlockNumber: chainHead.BlockNumber + 1}
		blockChain = &newBlock
	}
	return blockChain
}

// ValidateTransaction : Gets the current balance of a node
func ValidateTransaction(chainHead *Block, transactional Transaction) bool {
	var balance int
	for chainHead != nil {
		for _, trans := range chainHead.BlockTransactions {
			if trans.Sender == transactional.Sender {
				balance -= trans.Amount

			} else if trans.Receiver == transactional.Sender {
				balance += trans.Amount
			}
		}
		chainHead = chainHead.HashPtr.PreviousBlockPtr
	}
	if transactional.Amount <= balance && transactional.Amount > 0 {
		return true
	}
	return false
}

//ListBlocks shows the complete blockchain
func ListBlocks(chainHead *Block) {
	if chainHead == nil {
		println("No Blocks Found")
	} else {
		fmt.Println("\n-----------------------------")
		fmt.Println("---------BLOCKCHAIN----------")
		fmt.Println("-----------------------------")
		for chainHead != nil {
			fmt.Println("Block # "+strconv.Itoa(chainHead.BlockNumber)+"  => Transactions  : ", chainHead.BlockTransactions, "  ::  Previous Hash  : "+chainHead.HashPtr.PreviousBlockHash)
			chainHead = chainHead.HashPtr.PreviousBlockPtr
		}
	}
}

//VerifyChain allows the miner to verify all the transactions
func VerifyChain(chainHead *Block) {
	if chainHead == nil {
		fmt.Println("No Blocks Found")
	} else {
		for chainHead.HashPtr.PreviousBlockPtr != nil {
			data := getAllTransactions(chainHead.HashPtr.PreviousBlockPtr)
			data += chainHead.HashPtr.PreviousBlockPtr.HashPtr.PreviousBlockHash + strconv.Itoa(chainHead.HashPtr.PreviousBlockPtr.BlockNumber)
			blockHash := GetHash(data)
			if chainHead.HashPtr.PreviousBlockHash != blockHash {
				fmt.Println("Blockchain currentPtrered, Block # " + strconv.Itoa(chainHead.BlockNumber-1) + " had been altered")
				return
			}
			chainHead = chainHead.HashPtr.PreviousBlockPtr
		}
		fmt.Println("\n-------------Blockchain Verified---------\n")
	}
}
