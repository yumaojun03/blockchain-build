package protocal

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"

// BlockChain keeps a sequence of Blocks
type BlockChain struct {
	tip []byte
	db  *bolt.DB
}

// BlockChainIterator use to iterator blocks
type BlockChainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

// Close use colse conn
func (bc *BlockChain) Close() {
	bc.db.Close()
}

// AddBlock saves provided data as a block in the blockchain
func (bc *BlockChain) AddBlock(data string) {
	var lastHash []byte
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(data, lastHash)

	bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Header.DataHash, SerializeBlock(newBlock))
		if err != nil {
			log.Panic(err)
		}
		err = b.Put([]byte("l"), newBlock.Header.DataHash)
		if err != nil {
			log.Panic(err)
		}
		bc.tip = newBlock.Header.DataHash

		return nil
	})
}

// Iterator use to new an bci
func (bc *BlockChain) Iterator() *BlockChainIterator {
	bci := &BlockChainIterator{bc.tip, bc.db}

	return bci
}

// NewBlockchain creates a new Blockchain with genesis Block
func NewBlockchain() *BlockChain {
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			fmt.Println("No existing blockchain found. Creating a new one...")
			genesis := NewGenesisBlock()

			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Panic(err)
			}

			err = b.Put(genesis.Header.DataHash, SerializeBlock(genesis))
			if err != nil {
				log.Panic(err)
			}

			err = b.Put([]byte("l"), genesis.Header.DataHash)
			if err != nil {
				log.Panic(err)
			}
			tip = genesis.Header.DataHash
		} else {
			tip = b.Get([]byte("l"))
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := &BlockChain{tip, db}

	return bc
}

// Next use to iter next block
func (i *BlockChainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block = DeSerializeBlock(encodedBlock)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	i.currentHash = block.Header.PreviousHash

	return block
}
