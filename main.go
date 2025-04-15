package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"blockchain-project/core" // Import your "core" package
)

func main() {
	// 1. Create a new Blockchain
	blockchain := core.NewBlockchain()

	reader := bufio.NewReader(os.Stdin) // For reading user input

	for {
		fmt.Println("\nChoose an action:")
		fmt.Println("1: Mine a new block")
		fmt.Println("2: View Blockchain")
		fmt.Println("3: Tamper with a block")
		fmt.Println("4: Check Blockchain Validity")
		fmt.Println("5: Check Merkel tree")
		fmt.Println("6: Exit")
		fmt.Print("Enter your choice (1-6): ")

		choiceStr, _ := reader.ReadString('\n')
		choiceStr = strings.TrimSpace(choiceStr)
		choice, err := strconv.Atoi(choiceStr)
		if err != nil {
			fmt.Println("Invalid input. Please enter a number between 1 and 5.")
			continue
		}

		switch choice {
		case 1:
			mineNewBlock(blockchain, reader)
		case 2:
			printBlockchain(blockchain)
		case 3:
			tamperWithBlock(blockchain, reader)
		case 4:
			checkValidity(blockchain)
		case 5:
			testMerkleTree() // Function to test Merkle Tree
		case 6: // New option: Test Merkle Tree
			fmt.Println("Exiting program.")
			return
		default:
			fmt.Println("Invalid choice. Please enter a number between 1 and 6.")
		}
	}
}

func mineNewBlock(blockchain *core.Blockchain, reader *bufio.Reader) {
	lastBlock := blockchain.GetLastBlock()
	fmt.Print("Enter block data: ")
	blockData, _ := reader.ReadString('\n')
	blockData = strings.TrimSpace(blockData)

	fmt.Println("\nMining new block...")
	nonce := ProofOfWork(lastBlock, blockData)
	newBlock := core.NewBlock(lastBlock.BlockNumber+1, lastBlock.Hash, blockData, nonce)

	if err := blockchain.AddBlock(newBlock); err != nil {
		fmt.Println("Error adding block:", err)
	} else {
		fmt.Println("Block mined and added to blockchain. Hash:", newBlock.Hash)
	}
}

func printBlockchain(blockchain *core.Blockchain) {
	fmt.Println("\n----- Blockchain -----")
	for _, block := range blockchain.Blocks {
		fmt.Printf("Block Number: %d\n", block.BlockNumber)
		fmt.Printf("Timestamp: %d\n", block.Timestamp)
		fmt.Printf("Previous Block Hash: %s\n", block.PrevBlockHash)
		fmt.Printf("Transactions: %s\n", block.Transactions)
		fmt.Printf("Nonce: %d\n", block.Nonce)
		fmt.Printf("Block Hash: %s\n", block.Hash)
		fmt.Println("-----------------------")
	}
}

func tamperWithBlock(blockchain *core.Blockchain, reader *bufio.Reader) {
	if len(blockchain.Blocks) <= 1 {
		fmt.Println("Not enough blocks to tamper with (need at least Block 1).")
		return
	}

	fmt.Print("Enter block number to tamper with (starting from 1): ")
	blockNumStr, _ := reader.ReadString('\n')
	blockNumStr = strings.TrimSpace(blockNumStr)
	blockNumber, err := strconv.Atoi(blockNumStr)
	if err != nil || blockNumber <= 0 || blockNumber >= len(blockchain.Blocks) {
		fmt.Println("Invalid block number.")
		return
	}

	fmt.Print("Enter new data for block: ")
	newData, _ := reader.ReadString('\n')
	newData = strings.TrimSpace(newData)

	tamperedBlock := blockchain.Blocks[blockNumber]
	originalHash := tamperedBlock.Hash
	tamperedBlock.Transactions = newData
	tamperedBlock.Hash = tamperedBlock.CalculateBlockHash()

	fmt.Println("\n--- Tampering with Block", blockNumber, "---")
	fmt.Println("Block", blockNumber, "Hash before tampering:", originalHash)
	fmt.Println("Block", blockNumber, "Hash after tampering:", tamperedBlock.Hash)
	fmt.Println("Blockchain validity after tampering:", core.IsChainValid(blockchain.Blocks))
}

func checkValidity(blockchain *core.Blockchain) {
	isValid := core.IsChainValid(blockchain.Blocks)
	fmt.Println("\nIs Blockchain valid?", isValid)
	if !isValid {
		fmt.Println("Blockchain is INVALID! Tampering detected.")
	}
}

// --- Very Basic Proof-of-Work (PoW) Example ---
// In a real PoW, you'd adjust difficulty and have more complex logic.
func ProofOfWork(lastBlock *core.Block, data string) int {
	nonce := 0
	targetPrefix := "0000" // Target prefix: Hashes should start with "0000" (adjust difficulty by changing prefix length)

	fmt.Println("Mining started...")
	startTime := time.Now()

	for {
		block := core.NewBlock(lastBlock.BlockNumber+1, lastBlock.Hash, data, nonce)
		hash := block.Hash

		if len(hash) >= len(targetPrefix) && hash[:len(targetPrefix)] == targetPrefix {
			elapsed := time.Since(startTime)
			fmt.Printf("Mining finished in %s. Hash: %s, Nonce: %d\n", elapsed, hash, nonce)
			return nonce
		}
		nonce++
		// Basic rate limiting to avoid excessive CPU usage in this example:
		if nonce%100000 == 0 {
			time.Sleep(10 * time.Millisecond) // Sleep briefly every 100,000 iterations
		}
	}
}

func testMerkleTree() {
	fmt.Println("\n--- Merkle Tree Test ---")

	dataList := []string{"data1", "data2", "data3", "data4", "data5"} // Example data
	fmt.Println("Data list for Merkle Tree:", dataList)

	merkleTree, err := core.NewMerkleTree(dataList)
	if err != nil {
		fmt.Println("Error creating Merkle Tree:", err)
		return
	}
	fmt.Println("Merkle Tree created:", merkleTree) // Use String() method for printing

	rootHash := merkleTree.GetMerkleRoot()
	fmt.Println("Merkle Root:", rootHash)

	// Test proof generation and verification for "data3"
	testData := "data3"
	proof, err := merkleTree.GenerateMerkleProof(testData)
	if err != nil {
		fmt.Println("Error generating Merkle Proof:", err)
		return
	}
	fmt.Println("Merkle Proof for", testData, ":", proof)

	isValidProof := merkleTree.VerifyMerkleProof(testData, proof, rootHash)
	fmt.Println("Is Merkle Proof valid for", testData, "?", isValidProof)

	// Test verification for incorrect data (should fail)
	invalidData := "wrong_data"
	isInvalidProofValid := merkleTree.VerifyMerkleProof(invalidData, proof, rootHash)               // Using proof for "data3"
	fmt.Println("Is Merkle Proof valid for", invalidData, "(incorrect data)?", isInvalidProofValid) // Should be false

	// Test verification with tampered proof (should fail)
	if len(proof) > 0 {
		tamperedProof := append([]string{}, proof[1:]...) // Remove first element of proof
		isTamperedProofValid := merkleTree.VerifyMerkleProof(testData, tamperedProof, rootHash)
		fmt.Println("Is Merkle Proof valid with tampered proof?", isTamperedProofValid) // Should be false
	}
}
