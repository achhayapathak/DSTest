// 	Problem statement:

// Implement a data structure for Merkle tree in Go where each leaf node represents a string value.
// 1. Implement a method that takes an array of strings as input and returns an instance of the trie (to get Merkle root of the resulting trie)
// 2. Implement a method to return a Merkle proof data structure that proves the existence of a string within the trie (can store references of leaf nodes in the Merkle trie data structure for the same)
// 3. Use the generated proof to verify that the value actually exists in the said Merkle trie

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type MerkleNode struct {
	Hashvalue  []byte
	LeftChild  *MerkleNode
	RightChild *MerkleNode
	Parent     *MerkleNode
}

func ConstructTree(values []string, LeafMap map[string]*MerkleNode) *MerkleNode {
	hashValues := []*MerkleNode{}

	for _, val := range values {

		sum := sha256.Sum256([]byte(val))
		node := MerkleNode{sum[:], nil, nil, nil}

		hashValues = append(hashValues, &node)

		LeafMap[val] = &node
	}

	for len(hashValues) > 1 {
		tempHashes := []*MerkleNode{}

		for i := 0; i < len(hashValues); i += 2 {
			tempArray := []byte{}

			tempArray = append(tempArray[:], hashValues[i].Hashvalue...)
			if i+1 == len(hashValues) {
				tempArray = append(tempArray[:], hashValues[i].Hashvalue...)
			} else {
				tempArray = append(tempArray[:], hashValues[i+1].Hashvalue...)
			}

			sum := sha256.Sum256([]byte(tempArray))
			var node MerkleNode

			if i+1 == len(hashValues) {
				node = MerkleNode{sum[:], hashValues[i], hashValues[i], nil}
				hashValues[i].Parent = &node
			} else {
				node = MerkleNode{sum[:], hashValues[i], hashValues[i+1], nil}
				hashValues[i].Parent = &node
				hashValues[i+1].Parent = &node
			}

			tempHashes = append(tempHashes, &node)
		}

		hashValues = tempHashes

	}

	return hashValues[0]
}

func MerkleProof(leafNode *MerkleNode) ([][]byte, []bool) {
	proof := [][]byte{}
	isLeft := []bool{}

	for leafNode.Parent != nil {
		if leafNode.Parent.LeftChild == leafNode {
			proof = append(proof, leafNode.Parent.RightChild.Hashvalue)
			isLeft = append(isLeft, false)
		} else {
			proof = append(proof, leafNode.Parent.LeftChild.Hashvalue)
			isLeft = append(isLeft, true)
		}

		leafNode = leafNode.Parent
	}

	return proof, isLeft
}

func VerifyMerkleProof(result *MerkleNode, node string, proof [][]byte, isLeft []bool) bool {
	nodeHash := sha256.Sum256([]byte(node))

	for i := 0; i < len(proof); i++ {
		tempArray := []byte{}
		if isLeft[i] {
			tempArray = append(tempArray[:], proof[i][:]...)
			tempArray = append(tempArray[:], nodeHash[:]...)
		} else {
			tempArray = append(tempArray[:], nodeHash[:]...)
			tempArray = append(tempArray[:], proof[i][:]...)
		}

		nodeHash = sha256.Sum256([]byte(tempArray))
	}

	return nodeHash == [32]byte(result.Hashvalue)
}

func main() {
	// Map to store LeafNode references.
	LeafMap := make(map[string]*MerkleNode)

	Input := []string{"alice", "bob", "charlie", "david", "erin", "fiona", "george", "hannah"}

	// result contains reference to the Merkle root of the tree
	result := ConstructTree(Input, LeafMap)

	// converting root hash to hexadecimal human readable form
	str := hex.EncodeToString(result.Hashvalue)

	fmt.Println("The Root Hash is:", str)

	// node value for which you want to verify the Merkle tree
	verificationNode := "david"

	// proof array contains the hashes and isLeft marks whether hash correspond to left child or right child
	proof, isLeft := MerkleProof(LeafMap[verificationNode])

	// fmt.Println("Merkle proof for david", proof, isLeft)

	verificationValue := VerifyMerkleProof(result, verificationNode, proof, isLeft)

	fmt.Println("Is tree correctly verified for", verificationNode, ":", verificationValue)
}
