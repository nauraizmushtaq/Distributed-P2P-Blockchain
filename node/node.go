package main

import (
	"encoding/gob"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"

	a1 "github.com/nauraizmushtaq/assignment03IBC/assignment01IBC_i160106"
)

type nodeData struct {
	blockChain           *a1.Block
	connectionNodeAddrss []string
	nmCoins              int
	nodeHostname         string
	nodePort             string
	serverAddress        string
}
type transactionParameter struct {
	Amount   int
	Receiver string
}

var transParameters transactionParameter
var nodeObject nodeData

//Flood ... flood blockchain to all other nodes after adding new Block
func Flood() {
	a1.ListBlocks(nodeObject.blockChain)
	for i := 0; i < len(nodeObject.connectionNodeAddrss); i++ {
		if nodeObject.connectionNodeAddrss[i] != nodeObject.nodePort {
			conn, _ := net.Dial("tcp", "localhost:"+nodeObject.connectionNodeAddrss[i])
			enc := gob.NewEncoder(conn)
			enc.Encode("Receive_Blockchain")
			enc.Encode(nodeObject.blockChain)
			conn.Close()
		}
	}
	fmt.Println("Blockchain Sent to all Connected Nodes ------------ ")
}

//performCoinbaseTransation... Add miner tracions which is coin base transaction
func performCoinbaseTransation(transToBeValidated a1.Transaction) {
	coinBaseTransaction := a1.Transaction{Sender: "Miner", Receiver: nodeObject.nodePort, Amount: 100}
	var trans []a1.Transaction
	fmt.Println("Perfroming Coins Base Transaction")
	fmt.Println(coinBaseTransaction)
	trans = append(trans, transToBeValidated)
	trans = append(trans, coinBaseTransaction)
	nodeObject.blockChain = a1.InsertBlock(trans, nodeObject.blockChain)

}

//mineBlock ... mine block
func mineBlock(dec *gob.Decoder) {
	var transToBeValidated a1.Transaction
	_ = dec.Decode(&transToBeValidated)
	fmt.Println(transToBeValidated)
	if a1.ValidateTransaction(nodeObject.blockChain, transToBeValidated) == true {
		fmt.Println("transaction Validated")
		performCoinbaseTransation(transToBeValidated)
		fmt.Println("Forwarding Chain to All Peers")
		Flood()
	} else {
		fmt.Println("transaction is not verfied, aviod double spending")
	}
}

//routeConnecitonRequests ...        Receive the Blockchain Structure and Network Peers addresses
func routeConnecitonRequests(c net.Conn) {
	var data string
	dec := gob.NewDecoder(c)
	dec.Decode(&data)
	if data == "Validate_Transaction" {
		fmt.Println("Validating Transaction")
		mineBlock(dec)
	}
	if data == "Receive_Data" {
		dec.Decode(&nodeObject.blockChain)
		dec.Decode(&nodeObject.connectionNodeAddrss)
		a1.VerifyChain(nodeObject.blockChain)
		a1.ListBlocks(nodeObject.blockChain)
		fmt.Println("Peers Address to which Current Node needs to establish Connection : ", nodeObject.connectionNodeAddrss)
	}
	if data == "Receive_Blockchain" {
		fmt.Println("Recieive Blockchain")
		dec.Decode(&nodeObject.blockChain)
		a1.VerifyChain(nodeObject.blockChain)
		a1.ListBlocks(nodeObject.blockChain)
	}
}

//Lisntr Thread
func listen() {
	ln, _ := net.Listen("tcp", ":"+nodeObject.nodePort)
	for {
		conn, _ := ln.Accept()
		go routeConnecitonRequests(conn)
	}
}

//MakeTransaction ...
func MakeTransaction() {
	var r1 = rand.New(rand.NewSource(time.Now().UnixNano()))
	flag := true
	var rndMinerPort string
	for flag == true {
		rndMinerPort = nodeObject.connectionNodeAddrss[r1.Intn(len(nodeObject.connectionNodeAddrss))]
		if rndMinerPort != nodeObject.nodePort {
			flag = false
		}
	}
	fmt.Println("Miner Port :: ", rndMinerPort)
	transToBeValidated := a1.Transaction{Sender: nodeObject.nodePort, Receiver: transParameters.Receiver, Amount: transParameters.Amount}
	conn, _ := net.Dial("tcp", "localhost:"+rndMinerPort)
	enc := gob.NewEncoder(conn)
	_ = enc.Encode("Validate_Transaction")
	_ = enc.Encode(transToBeValidated)
	_ = conn.Close()
}

func main() {
	nodeObject.nodeHostname = os.Args[3]
	nodeObject.nodePort = os.Args[2]
	nodeObject.serverAddress = os.Args[1]
	fmt.Println("Connecting to Stoshi Server .... ")
	conn, _ := net.Dial("tcp", "localhost:"+nodeObject.serverAddress)
	fmt.Println("Connected .... ")
	enc := gob.NewEncoder(conn)
	_ = enc.Encode("Address")
	_ = enc.Encode(nodeObject.nodeHostname + ":" + nodeObject.nodePort)
	fmt.Println("Node Listening at ", nodeObject.nodePort)
	ln, _ := net.Listen("tcp", ":"+nodeObject.nodePort)
	conn, _ = ln.Accept()
	routeConnecitonRequests(conn)
	_ = ln.Close()
	//reading an integer
	go listen()
	for {
		fmt.Println("------------ To Perfrom Transaction ------------")
		fmt.Println("Node Address : ", nodeObject.nodePort)
		fmt.Println("Peers Addresses : ", nodeObject.connectionNodeAddrss)
		fmt.Println("Enter Amount to Transfer : ")
		_, _ = fmt.Scan(&transParameters.Amount)
		fmt.Println("Enter Reciever Port : ")
		_, _ = fmt.Scan(&transParameters.Receiver)
		MakeTransaction()
	}
}
