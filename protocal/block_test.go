package protocal_test

import (
	"fmt"
	"testing"

	proto "blockchain_go/protocal"
)

func Test_NewBlock(t *testing.T) {
	bc := proto.NewBlockchain()
	bc.AddBlock("Send 1 BTC to Ivan")
	bc.AddBlock("Send 2 more BTC to Ivan")

	for _, block := range bc.Blocks {
		fmt.Printf("Prev. hash: %x\n", block.Header.PreviousHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Header.DataHash)
		fmt.Println()
	}
}
