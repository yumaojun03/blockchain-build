package cmd

import (
	"blockchain_go/protocal"
	"flag"
	"fmt"
	"log"
	"os"
)

type CLI struct {
	bc *protocal.BlockChain
}

func NewCLI(bc *protocal.BlockChain) *CLI {
	return &CLI{bc}
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) Run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	addBlockData := addBlockCmd.String("data", "", "Block data")

	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	// parse cli arguments
	switch os.Args[1] {
	case "addblock":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	//
	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cli.addBlock(*addBlockData)

	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  addblock -data BLOCK_DATA - add a block to the blockchain")
	fmt.Println("  printchain - print all the blocks of the blockchain")
}

func (cli *CLI) addBlock(data string) {
	cli.bc.AddBlock(data)
	fmt.Println("Success!")
}

func (cli *CLI) printChain() {
	bci := cli.bc.Iterator()

	for {
		block := bci.Next()

		fmt.Printf("Prev. hash: %x\n", block.Header.PreviousHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Header.DataHash)
		fmt.Printf("Nonce: %d\n", block.Header.Nonce)
		pow := protocal.NewProofOfWork(block)
		fmt.Printf("POW: %t\n", pow.Validate())
		fmt.Println()

		if len(block.Header.PreviousHash) == 0 {
			break
		}
	}
}
