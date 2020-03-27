package main

import (
	"encoding/gob"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	a1 "github.com/nauraizmushtaq/assignment03IBC/assignment01IBC_i160106"

	"fmt"
)

type stoshiData struct {
	blockChain     *a1.Block
	connectedNodes []string
	quorumNumber   int
	nmCoins        int
}
type transactionParameter struct {
	Amount   int
	Receiver string
}

var transParameters transactionParameter
var stoshiObject stoshiData
var mutex = make(chan bool)
var port = os.Args[1]

//sendData ... Send BLokchain and Adddress
func sendData() {
	stoshiObject.connectedNodes = append(stoshiObject.connectedNodes, port)
	for i := 0; i < len(stoshiObject.connectedNodes); i++ {
		conn, _ := net.Dial("tcp", "localhost:"+stoshiObject.connectedNodes[i])
		enc := gob.NewEncoder(conn)
		enc.Encode("Receive_Data")
		enc.Encode(stoshiObject.blockChain)
		enc.Encode(stoshiObject.connectedNodes)
		conn.Close()
	}
	fmt.Println("Blockchain Sent to all Connected Nodes along with Peers Address  ...... ")
}

//preming done by stoshi
func preMineBlock() {
	fmt.Println("New Blcok Mined")
	fmt.Println("Transaction Mined .... Stoshi gets 100 nmCoins for Mining")
	stoshiTrans := a1.Transaction{Sender: "Miner", Receiver: port, Amount: 100}
	var trans []a1.Transaction
	trans = append(trans, stoshiTrans)
	stoshiObject.blockChain = a1.InsertBlock(trans, stoshiObject.blockChain)
	stoshiObject.nmCoins += 100
}

//MakeTransaction ...
func MakeTransaction() {
	var r1 = rand.New(rand.NewSource(time.Now().UnixNano()))
	flag := true
	var rndMinerPort string
	for flag == true {
		rndMinerPort = stoshiObject.connectedNodes[r1.Intn(len(stoshiObject.connectedNodes))]
		if rndMinerPort != port {
			flag = false
		}
	}
	fmt.Println("Miner Port :: ", rndMinerPort)
	transToBeValidated := a1.Transaction{Sender: port, Receiver: transParameters.Receiver, Amount: transParameters.Amount}
	conn, _ := net.Dial("tcp", "localhost:"+rndMinerPort)
	enc := gob.NewEncoder(conn)
	_ = enc.Encode("Validate_Transaction")
	_ = enc.Encode(transToBeValidated)
	_ = conn.Close()
}

//routeConnecitonRequests ...        Receive the Blockchain Structure and Network Peers addresses
func performCoinbaseTransation(transToBeValidated a1.Transaction) {
	coinBaseTransaction := a1.Transaction{Sender: "Miner", Receiver: port, Amount: 100}
	var trans []a1.Transaction
	fmt.Println("Perfroming Coins Base Transaction")
	fmt.Println(coinBaseTransaction)
	trans = append(trans, transToBeValidated)
	trans = append(trans, coinBaseTransaction)
	stoshiObject.blockChain = a1.InsertBlock(trans, stoshiObject.blockChain)

}

//Flood ... Send BLokchain and Adddress
func Flood() {
	a1.ListBlocks(stoshiObject.blockChain)
	for i := 0; i < len(stoshiObject.connectedNodes); i++ {
		if stoshiObject.connectedNodes[i] != port {
			conn, _ := net.Dial("tcp", "localhost:"+stoshiObject.connectedNodes[i])
			enc := gob.NewEncoder(conn)
			enc.Encode("Receive_Blockchain")
			enc.Encode(stoshiObject.blockChain)
			conn.Close()

		}

	}
	fmt.Println("Blockchain Sent to all Connected Nodes ----------------------- ")
}

//mineBlock ... mine block
func mineBlock(dec *gob.Decoder) {
	var transToBeValidated a1.Transaction
	_ = dec.Decode(&transToBeValidated)
	fmt.Println(transToBeValidated)
	if a1.ValidateTransaction(stoshiObject.blockChain, transToBeValidated) == true {
		fmt.Println("transaction Validated")
		performCoinbaseTransation(transToBeValidated)
		fmt.Println("Block Mined Forwarding Blockchain to Other Peers")
		Flood()
	} else {
		fmt.Println("transaction is not verfied, aviod double spending")
	}
}

//routeConnRequest ... Get Node Address and wait till qourum Completes
func routeConnRequest(c net.Conn) {
	var data string
	dec := gob.NewDecoder(c)
	dec.Decode(&data)
	if data == "Address" {
		dec.Decode(&data)
		address := strings.Split(data, ":")
		stoshiObject.quorumNumber--
		fmt.Println("Connected to ", address[0], ":", address[1])
		stoshiObject.connectedNodes = append(stoshiObject.connectedNodes, address[1])
		preMineBlock()
		fmt.Println("Stoshi Current Balance : ", stoshiObject.nmCoins)
		if stoshiObject.quorumNumber == 0 {
			fmt.Println("Quorum Completed ...... ")
			sendData()
			mutex <- true
			return
		}
		mutex <- false
		return
	}
	if data == "Receive_Blockchain" {
		fmt.Println("Recieive Blockchain")
		dec.Decode(&stoshiObject.blockChain)
		a1.VerifyChain(stoshiObject.blockChain)
		a1.ListBlocks(stoshiObject.blockChain)
	}
	if data == "Validate_Transaction" {
		fmt.Println("Validating the Transaction")
		mineBlock(dec)
	}
}

func listen() {
	ln, _ := net.Listen("tcp", ":"+port)
	for {
		conn, _ := ln.Accept()
		go routeConnRequest(conn)
	}
}

func main() {
	stoshiObject.quorumNumber, _ = strconv.Atoi(os.Args[2])
	ln, _ := net.Listen("tcp", ":"+port)
	for {
		fmt.Println("Satoshi Server Listening at ", port)
		conn, _ := ln.Accept()
		go routeConnRequest(conn)
		if <-mutex {
			break
		}
		fmt.Println("Nodes Remained : ", stoshiObject.quorumNumber)
	}
	_ = ln.Close()
	go listen()
	for {
		fmt.Println("------------ To Perfrom Transaction ------------")
		fmt.Println("Peers Address : ", stoshiObject.connectedNodes)
		fmt.Println("Node Address : ", port)
		fmt.Println("Enter Amount to Transfer : ")
		_, _ = fmt.Scan(&transParameters.Amount)
		fmt.Println("Enter Reciever Port : ")
		_, _ = fmt.Scan(&transParameters.Receiver)
		fmt.Println(transParameters.Receiver, "+", transParameters.Amount)
		MakeTransaction()
	}
}
