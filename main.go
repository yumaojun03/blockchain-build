package main

import (
	"blockchain_go/cmd"
	"blockchain_go/protocal"
)

func main() {
	bc := protocal.NewBlockchain()
	defer bc.Close()

	cli := cmd.NewCLI(bc)
	cli.Run()
}
