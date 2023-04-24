package merkletree

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"testing"
)

func HashFuncSHA256(data []byte) ([]byte, error) {
	hash := sha256.New()
	hash.Write(data)
	hashValue := hash.Sum(nil)
	//fmt.Println("Calculated hash is :", hashValue)
	return hashValue, nil
}

func TestMerkleBasics(t *testing.T) {
	var data [][]byte
	expectedHash := []byte{95, 48, 204, 128, 19, 59, 147, 148, 21, 110, 36, 178, 51, 240, 196, 190, 50, 178, 78, 68, 187, 51, 129, 240, 44, 123, 165, 38, 25, 208, 254, 188}
	strings := []string{
		"Hello",
		"Hi",
		"Hey",
		"Hola",
	}
	for _, val := range strings {
		d := []byte(val)
		data = append(data, d)
	}

	tree, err := NewTree(data, HashFuncSHA256)
	if err != nil {
		fmt.Println("Failed to build Tree due to error ", err)
		t.FailNow()
	}
	//TODO: Add a verification for expected rootHash for input data.
	//fmt.Println("root hash is ", tree.RootHash())
	fmt.Printf("Tree : \n %v", tree)
	if !bytes.Equal(tree.RootHash(), expectedHash) {
		fmt.Println("Root hash generated is not matching expected hash")
		t.FailNow()
	}
	proof, err := tree.GenerateMerkleProof(data[0])
	if err != nil {
		fmt.Println("Failed to generate Merkle Proof due to error ", err)
		t.FailNow()
	}
	fmt.Println("Proof generated is : \n", proof)

	verified, err := tree.VerifyProof(data[0], proof)
	if err != nil {
		fmt.Println("Failed to verify Merkle Proof due to error ", err)
		t.FailNow()
	}
	if !verified {
		fmt.Println("Invalid proof for the data", err)
		t.FailNow()
	}
}
