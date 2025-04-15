package core

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
)

// MerkleNode represents a node in the Merkle Tree
type MerkleNode struct {
	Hash  string
	Left  *MerkleNode
	Right *MerkleNode
	Data  string // Optional: Keep the original data in leaf nodes
}

// MerkleTree represents the Merkle Tree structure
type MerkleTree struct {
	Root     *MerkleNode
	LeafData []string
}

// NewMerkleTree builds a Merkle Tree from a list of data strings.  <-- **CAPITAL 'N' in NewMerkleTree is crucial for export**
func NewMerkleTree(dataList []string) (*MerkleTree, error) {
	if len(dataList) == 0 {
		return nil, errors.New("cannot build Merkle Tree from empty data list")
	}
	tree := &MerkleTree{LeafData: dataList}
	tree.Root = buildMerkleTreeNodes(convertDataListToNodes(dataList))
	return tree, nil
}

// convertDataListToNodes converts a list of data strings into a list of MerkleNodes (leaf nodes).
func convertDataListToNodes(dataList []string) []*MerkleNode {
	var nodes []*MerkleNode
	for _, data := range dataList {
		hash := calculateHash(data) // You'll need to implement calculateHash
		nodes = append(nodes, &MerkleNode{
			Hash: hash,
			Data: data, // Store the original data in leaf nodes
		})
	}
	return nodes
}

// buildMerkleTreeNodes recursively builds the Merkle Tree from a list of MerkleNodes.
func buildMerkleTreeNodes(nodes []*MerkleNode) *MerkleNode {
	if len(nodes) == 1 {
		return nodes[0] // Base case: only one node left, it's the root
	}

	var parentNodes []*MerkleNode
	for i := 0; i < len(nodes); i += 2 {
		leftNode := nodes[i]
		var rightNode *MerkleNode
		if i+1 < len(nodes) {
			rightNode = nodes[i+1]
		} else {
			// If odd number of nodes, duplicate the last node to make a pair
			rightNode = leftNode
		}

		// Calculate hash for parent node by combining hashes of children
		parentHash := calculateHash(leftNode.Hash + rightNode.Hash) // You'll need to implement calculateHash

		parentNodes = append(parentNodes, &MerkleNode{
			Hash:  parentHash,
			Left:  leftNode,
			Right: rightNode,
		})
	}

	return buildMerkleTreeNodes(parentNodes) // Recursive call with parent nodes
}

// calculateHash calculates the SHA-256 hash of a string and returns its hex representation.
func calculateHash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// GetMerkleRoot returns the root hash of the Merkle Tree.
func (mt *MerkleTree) GetMerkleRoot() string {
	if mt.Root != nil {
		return mt.Root.Hash
	}
	return ""
}

// GenerateMerkleProof generates a Merkle Proof for a given data string.// GenerateMerkleProof generates a Merkle Proof for a given data string.
func (mt *MerkleTree) GenerateMerkleProof(data string) ([]string, error) {
	var proof []string
	node := findLeafNode(mt.Root, data) // Helper function to find the leaf node
	if node == nil {
		return nil, errors.New("data not found in Merkle Tree")
	}
	fmt.Println("Generating proof for data:", data, ", leaf hash:", node.Hash) // Debug

	// Traverse up from the leaf node to the root, collecting sibling hashes
	current := node
	for current != mt.Root {
		parent := findParent(mt.Root, current) // Helper function to find parent
		if parent == nil {
			return nil, errors.New("could not find parent node") // Should not happen in a valid tree
		}

		var sibling *MerkleNode
		if parent.Left == current {
			sibling = parent.Right
			if sibling != nil {
				fmt.Println("  Current node is LEFT child, sibling (RIGHT) hash:", sibling.Hash) // Debug
				proof = append(proof, sibling.Hash)
			}
		} else {
			sibling = parent.Left
			if sibling != nil {
				fmt.Println("  Current node is RIGHT child, sibling (LEFT) hash:", sibling.Hash) // Debug
				proof = append(proof, sibling.Hash)
			}
		}
		current = parent
		fmt.Println("  Moving up to parent hash:", parent.Hash) // Debug
	}
	fmt.Println("Generated Merkle Proof:", proof) // Debug
	return proof, nil
}

// Helper function to find a leaf node by data (for proof generation)
func findLeafNode(root *MerkleNode, data string) *MerkleNode {
	targetHash := calculateHash(data) // Hash the data we are searching for

	var findNodeRecursive func(node *MerkleNode) *MerkleNode
	findNodeRecursive = func(node *MerkleNode) *MerkleNode {
		if node == nil {
			return nil
		}
		if node.Left == nil && node.Right == nil && node.Hash == targetHash { // Leaf node check based on hash
			return node
		}

		leftResult := findNodeRecursive(node.Left)
		if leftResult != nil {
			return leftResult
		}
		return findNodeRecursive(node.Right)
	}
	return findNodeRecursive(root)
}

// Helper function to find the parent of a node (for proof generation)
func findParent(root *MerkleNode, child *MerkleNode) *MerkleNode {
	if root == nil || (root.Left == nil && root.Right == nil) { // If root is nil or a leaf, no parent
		return nil
	}
	if root.Left == child || root.Right == child {
		return root
	}

	parent := findParent(root.Left, child)
	if parent != nil {
		return parent
	}
	return findParent(root.Right, child)
}

// VerifyMerkleProof verifies a Merkle Proof for a given data string against a root hash.
// VerifyMerkleProof verifies a Merkle Proof for a given data string against a root hash.
func (mt *MerkleTree) VerifyMerkleProof(data string, proof []string, rootHash string) bool {
	calculatedHash := calculateHash(data)                                                                          // Start with the hash of the data
	fmt.Println("Verifying proof for data:", data, ", initial hash:", calculatedHash, ", against root:", rootHash) // Debug
	fmt.Println("Proof received:", proof)                                                                          // Debug

	for _, proofHash := range proof {
		fmt.Println("  Current calculated hash:", calculatedHash, ", proof hash:", proofHash) // Debug
		calculatedHash = calculateHash(calculatedHash + proofHash)                            // **Correct order: data/current hash on the LEFT, proof hash (sibling) on the RIGHT**
		fmt.Println("  New calculated hash:", calculatedHash)                                 // Debug
	}

	fmt.Println("Final calculated hash:", calculatedHash) // Debug
	isValid := calculatedHash == rootHash
	fmt.Println("Is Valid?", isValid) // Debug
	return isValid
}

// String method for MerkleTree for easy printing
func (mt *MerkleTree) String() string {
	return fmt.Sprintf("MerkleTree{RootHash: %s, LeafData: %v}", mt.GetMerkleRoot(), mt.LeafData)
}
