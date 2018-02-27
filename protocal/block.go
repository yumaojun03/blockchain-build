package protocal

import (
	"bytes"
	"crypto/sha256"
	"strconv"
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
}

// BlockData use to contain any data, eg transaction
type BlockData struct {
	Data [][]byte `json:"data,omitempty"`
}

// BlockMetadata use to save the data's meta information
type BlockMetadata struct {
	Metadata [][]byte `json:"metadata,omitempty"`
}

// SetHash use to set block hash
func (b *Block) SetHash() {
	ts := []byte(strconv.FormatInt(b.Header.Timestamp, 10))
	ph := b.Header.PreviousHash
	data := bytes.Join(b.Data.Data, []byte{})

	hashData := bytes.Join([][]byte{ph, data, ts}, []byte{})
	hash := sha256.Sum256(hashData)
	b.Header.DataHash = hash[:]
}

// NewBlock construct a block with no data and no metadata.
func NewBlock(data string, previousHash []byte) *Block {
	header := &BlockHeader{1, time.Now().Unix(), previousHash, []byte{}}
	dataB := &BlockData{[][]byte{[]byte(data)}}

	block := &Block{header, dataB, nil}
	block.SetHash()

	return block
}

// NewGenesisBlock creates and returns genesis Block
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}
