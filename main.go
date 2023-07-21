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
	// Parent *MerkleNode
}

func ConstructTree(values []string) *MerkleNode {
	hashValues := []MerkleNode{}
	for _, val := range values {
		// hash values
		sum := sha256.Sum256([]byte(val))
		node := MerkleNode{sum[:], nil, nil}
		hashValues = append(hashValues, node)
	}
	
	for len(hashValues) > 1 {
		tempHashes := []MerkleNode{}
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
				node = MerkleNode{sum[:], &hashValues[i], &hashValues[i]}
			} else {
				node = MerkleNode{sum[:], &hashValues[i], &hashValues[i+1]}
			}
			tempHashes = append(tempHashes, node)
		}

		hashValues = tempHashes

	}

	return &hashValues[0]
}

func main() {

	// 	Problem statement:

	// Implement a data structure for Merkle tree in Go where each leaf node represents a string value.
	// 1. Implement a method that takes an array of strings as input and returns an instance of the trie (to get Merkle root of the resulting trie)
	// 2. Implement a method to return a Merkle proof data structure that proves the existence of a string within the trie (can store references of leaf nodes in the Merkle trie data structure for the same)
	// 3. Use the generated proof to verify that the value actually exists in the said Merkle trie

	Input := []string{"alice", "bob", "charlie", "david", "erin", "fiona", "george", "hannah"}

	result := ConstructTree(Input)

	str := hex.EncodeToString(result.Hashvalue)

	fmt.Println("The Root Hash is:", str)
}
