package core

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

// Block represents a block in the blockchain
type Block struct {
	BlockNumber   int    // The block number or height in the chain
	Timestamp     int64  // Timestamp of when the block was created
	PrevBlockHash string // Hash of the previous block in the chain
	Transactions  string // For now, let's keep transactions simple as a string (we'll refine later)
	Nonce         int    // Nonce for Proof-of-Work
	Hash          string // The hash of the current block
}

// NewBlock creates a new Block.  Note: Hash is calculated later.
func NewBlock(blockNumber int, prevBlockHash string, transactions string, nonce int) *Block {
	block := &Block{
		BlockNumber:   blockNumber,
		Timestamp:     time.Now().UnixNano(),
		PrevBlockHash: prevBlockHash,
		Transactions:  transactions,
		Nonce:         nonce,
	}
	block.Hash = block.CalculateBlockHash() // Calculate hash immediately after creation
	return block
}

// CalculateBlockHash calculates the hash of the block.
func (b *Block) CalculateBlockHash() string {
	record := string(b.BlockNumber) + string(b.Timestamp) + b.PrevBlockHash + b.Transactions + string(b.Nonce)
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}
