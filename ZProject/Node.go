package assignment02IBC_master

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
)

type Node struct {
	Address    string
	BlockChain *Block
	Nodes      map[string]bool
	Votes      map[string]string
}

func (node *Node) SendBlockChain(address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println(err)
		return
	}
	encoder := gob.NewEncoder(conn)
	err = encoder.Encode("ReceiveBlockChain")

	err = encoder.Encode(node.BlockChain)
	err = conn.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (node *Node) ReceiveBlockChain(/*conn net.Conn*/decoder *gob.Decoder) {
	fmt.Println("decode started in receive blockchain")
	var blockChain *Block
	//decoder := gob.NewDecoder(conn)
	err := decoder.Decode(&blockChain)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("decode successful in client receive blockchain")
	if VerifyChain(blockChain) {
		node.BlockChain = blockChain
		fmt.Println("blockchain received")
		ListBlocks(node.BlockChain)
		return
	}
	fmt.Println("received blockchain not verified")

}

func (node *Node) MineBlock(transaction Transaction) {
	if (transaction.Sender != "early joiners reward") && (CalculateAmount(transaction.Sender, *node.BlockChain) > transaction.Amount) {
		fmt.Println("not enough balance")
		return
	}
	if transaction.IsEmpty() {
		fmt.Println("transaction is empty")
		return
	}
	t2 := Transaction{
		Amount:   1,
		Sender:   "mining reward",
		Receiver: node.Address,
	}

	var trns []Transaction
	trns = append(trns, transaction)
	trns = append(trns, t2)
	node.BlockChain = InsertBlock(trns, node.BlockChain)
}
func (node *Node) VerifyBlock() {

}
func (node *Node) InitiateTransaction() {

}
func (node *Node) ReceiveTransaction() {

}
func (node *Node) SendVote() {

}
func (node *Node) ReceiveVote() {

}

func (node *Node) VerifyVote() {

}

func (node *Node) DetermineNextMiner() {

}
func (node *Node) SendNodes(address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println(err)
		return
	}
	encoder := gob.NewEncoder(conn)

	err = encoder.Encode("ReceiveNodes")
	err = encoder.Encode(node.Nodes)
	err = conn.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (node *Node) ReceiveNodes(/*conn net.Conn*/decoder *gob.Decoder) {
	var nodes map[string]bool
	fmt.Println("decode started in receive nodes")
	//decoder := gob.NewDecoder(conn)
	err := decoder.Decode(&nodes)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("decode successful in client receive nodes")
	for n := range nodes {
		if n == node.Address {
			continue
		}
		node.AddNode(n)
	}
}
func (node *Node) AddNode(address string) {
	if node.Nodes == nil {
		node.Nodes = make(map[string]bool)
	}
	node.Nodes[address] = true
	fmt.Println("no of nodes inside function: ")
	fmt.Println(len(node.Nodes))

}
func (node *Node) FloodNodes() {
	for currentNode := range node.Nodes {
		node.SendNodes(currentNode)
	}
}
func (node *Node) FloodBlockChain() {
	for currentNode := range node.Nodes {
		node.SendBlockChain(currentNode)
	}
}
func (node *Node) HandleConnections(conn net.Conn) {
	decoder := gob.NewDecoder(conn)
	fmt.Println("decode started in client handle connection")
	var message string
	err := decoder.Decode(&message)
	if err != nil {
		fmt.Println("error from client handle connection")
		fmt.Println(err)
		return
	}
	fmt.Println("decode successful in client handle connection. message: "+message)
	switch message {
	case "ReceiveNodes":
		node.ReceiveNodes(decoder)
	case "ReceiveBlockChain":
		node.ReceiveBlockChain(decoder)
	}

}

func (node *Node) ListenConnections() {
	ln, err := net.Listen("tcp", node.Address)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("error from client listen connection")
			log.Println(err)
			return
		}
		fmt.Println("connection accepted successfully")
		go node.HandleConnections(conn)

	}

}

func (node *Node) HandleConnectionsSatoshi(conn net.Conn, channel chan bool) {
	decoder := gob.NewDecoder(conn)
	var address string
	err := decoder.Decode(&address)
	if err != nil {
		fmt.Println(err)
		channel <- false
		return
	}
	fmt.Println("connected to: " + address)
	node.AddNode(address)
	node.MineBlock(Transaction{
		Amount:   10,
		Sender:   "early joiners reward",
		Receiver: address,
	})
	//fmt.Println("no of nodes: ")
	//fmt.Println(len(node.Nodes))
	if len(node.Nodes) < 4 {
		channel <- false
	} else {
		channel <- true
	}
	return
}

func (node *Node) ListenConnectionsSatoshi() {

	ln, err := net.Listen("tcp", node.Address)
	if err != nil {
		log.Fatal(err)
	}
	channel := make(chan bool)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			return
		}
		go node.HandleConnectionsSatoshi(conn, channel)
		if <-channel {
			break
		}
	}
	_ = ln.Close()

}
