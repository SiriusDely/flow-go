package main

import (
	"fmt"

	"github.com/dapperlabs/bamboo-node/internal/roles/consensus"
)

func main() {
	server, err := consensus.InitializeServer()
	if err != nil {
		panic(err)
	}
	fmt.Println("TEST")

	server.Start()
}
