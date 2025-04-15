package core

import (
	"fmt"
)

// Blockchain represents the entire blockchain
type Blockchain struct {
	Blocks []*Block // A slice of Blocks, forming the chain
}

// NewBlockchain creates a new Blockchain and initializes it with the Genesis Block.
func NewBlockchain() *Blockchain {
	genesisBlock := createGenesisBlock() // Create the very first block
	return &Blockchain{
		Blocks: []*Block{genesisBlock}, // Initialize the blockchain with the genesis block
	}
}

// createGenesisBlock is a private function to create the very first block of the blockchain (Block #0).
func createGenesisBlock() *Block {
	// Genesis block has BlockNumber 0, no previous block hash, and some initial data (can be empty or predefined).
	genesisBlock := NewBlock(0, "", "Genesis Block Transactions", 0) // No previous hash for Genesis
	fmt.Println("Genesis block created:", genesisBlock.Hash)         // Informative message
	return genesisBlock
}

// AddBlock adds a new block to the blockchain, after performing basic validation.
func (bc *Blockchain) AddBlock(newBlock *Block) error {
	lastBlock := bc.GetLastBlock() // Get the last block in the chain

	// Basic validation:
	if newBlock.BlockNumber != lastBlock.BlockNumber+1 {
		return fmt.Errorf("invalid block number. Expected: %d, got: %d", lastBlock.BlockNumber+1, newBlock.BlockNumber)
	}
	if newBlock.PrevBlockHash != lastBlock.Hash {
		return fmt.Errorf("invalid previous block hash. Expected: %s, got: %s", lastBlock.Hash, newBlock.PrevBlockHash)
	}
	if newBlock.CalculateBlockHash() != newBlock.Hash { // Re-calculate hash and compare
		return fmt.Errorf("invalid block hash. Calculated hash doesn't match block's hash.")
	}

	bc.Blocks = append(bc.Blocks, newBlock) // Add the new valid block to the chain
	return nil
}

// GetLastBlock returns the latest block in the blockchain
func (bc *Blockchain) GetLastBlock() *Block {
	if len(bc.Blocks) == 0 {
		return nil // Should not happen in a properly initialized blockchain, but handle for robustness
	}
	return bc.Blocks[len(bc.Blocks)-1] // Return the last block in the slice
}

// IsBlockValid checks if a given block is valid in the context of the blockchain (basic checks).
func (bc *Blockchain) IsBlockValid(block *Block) bool {
	lastBlock := bc.GetLastBlock()

	if block.BlockNumber != lastBlock.BlockNumber+1 {
		return false
	}
	if block.PrevBlockHash != lastBlock.Hash {
		return false
	}
	if block.CalculateBlockHash() != block.Hash {
		return false
	}
	return true
}

// ReplaceChain replaces the current blockchain with a new chain.
// This is a very basic chain replacement - in a real system, you'd need more sophisticated chain selection logic.
func (bc *Blockchain) ReplaceChain(newBlocks []*Block) {
	if len(newBlocks) > len(bc.Blocks) && IsChainValid(newBlocks) { // Basic: New chain is longer and valid
		fmt.Println("Replacing current chain with a longer valid chain.")
		bc.Blocks = newBlocks
	} else {
		fmt.Println("Received chain is not longer or invalid. Ignoring.")
	}
}

// IsChainValid checks if a given blockchain is valid.
func IsChainValid(chain []*Block) bool {
	if len(chain) == 0 {
		return true // Empty chain is considered valid (for now, could be changed)
	}
	// Check Genesis block (block #0) - you might want to add specific genesis block validation
	if chain[0].BlockNumber != 0 {
		return false
	}
	if chain[0].PrevBlockHash != "" {
		return false
	}
	if chain[0].CalculateBlockHash() != chain[0].Hash {
		return false
	}

	// Check subsequent blocks
	for i := 1; i < len(chain); i++ {
		currentBlock := chain[i]
		prevBlock := chain[i-1]

		if currentBlock.BlockNumber != prevBlock.BlockNumber+1 {
			return false
		}
		if currentBlock.PrevBlockHash != prevBlock.Hash {
			return false
		}
		if currentBlock.CalculateBlockHash() != currentBlock.Hash {
			return false
		}
	}
	return true // If all checks pass, the chain is valid
}
