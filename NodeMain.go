package main

import (
	Project "./ZProject"
	"encoding/gob"
	"fmt"
	"net"
	"os"
)

func main() {

	conn, err := net.Dial("tcp", "localhost:3001")
	if err != nil {
		fmt.Println(err)
		return
	}
	encoder := gob.NewEncoder(conn)
	err = encoder.Encode("localhost:" + os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	node := Project.Node{
		Address:    ":" + os.Args[1],
		BlockChain: nil,
		Nodes:      nil,
	}
	node.ListenConnections()

}
