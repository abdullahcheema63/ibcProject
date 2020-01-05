package main

import (
	Project "./ZProject"
	"fmt"
)

func main() {

	node := Project.Node{
		Address:    ":3001",
		BlockChain: nil,
		Nodes:      nil,
	}
	node.ListenConnectionsSatoshi()
	node.FloodNodes()
	node.FloodBlockChain()
	_, _ = fmt.Scanln()
}
