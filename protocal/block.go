package protocal

import (
	"bytes"
	"encoding/gob"
	"time"
)

// Block is finalized block structure to be shared among the orderer and peer
// Note that the BlockHeader chains to the previous BlockHeader, and the BlockData hash is embedded
// in the BlockHeader.  This makes it natural and obvious that the Data is included in the hash, but
// the Metadata is not.
type Block struct {
	Header   *BlockHeader   `json:"header,omitempty"`
	Data     *BlockData     `json:"data,omitempty"`
	MetaData *BlockMetadata `json:"meta,omitempty"`
}

// BlockHeader is the element of the block which forms the block chain
// The block header is hashed using the configured chain hashing algorithm
// over the ASN.1 encoding of the BlockHeader
type BlockHeader struct {
	Version      uint64 `json:"version,omitempty"`
	Timestamp    int64  `json:"timestamp,omitempty"`
	PreviousHash []byte `json:"previous_hash,omitempty"`
	DataHash     []byte `json:"data_hash,omitempty"`
	Nonce        int    `json:"nonce,omitempty"`
}

// BlockData use to contain any data, eg transaction
type BlockData struct {
	Data [][]byte `json:"data,omitempty"`
}

// BlockMetadata use to save the data's meta information
type BlockMetadata struct {
	Metadata [][]byte `json:"metadata,omitempty"`
}

// NewBlock construct a block with no data and no metadata.
func NewBlock(data string, previousHash []byte) *Block {
	header := &BlockHeader{1, time.Now().Unix(), previousHash, []byte{}, 0}
	dataB := &BlockData{[][]byte{[]byte(data)}}

	block := &Block{header, dataB, nil}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Header.Nonce = nonce
	block.Header.DataHash = hash[:]

	return block
}

// NewGenesisBlock creates and returns genesis Block
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

func SerializeBlock(b *Block) []byte {
	var buf bytes.Buffer

	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(b); err != nil {
		panic(err)
	}

	return buf.Bytes()
}

func DeSerializeBlock(d []byte) *Block {
	block := new(Block)

	decoder := gob.NewDecoder(bytes.NewReader(d))
	if err := decoder.Decode(block); err != nil {
		panic(err)
	}

	return block
}
