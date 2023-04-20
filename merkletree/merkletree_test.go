package merkletree

import (
	"crypto/sha256"
	"fmt"
	"testing"
)

func HashFuncSHA256(data []byte) ([]byte, error) {
	hash := sha256.New()
	hash.Write(data)
	hashValue := hash.Sum(nil)
	fmt.Println("Calculated hash is :", hashValue)
	return hashValue, nil
}

func TestMerkleBasics(t *testing.T) {
	var data [][]byte
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
	fmt.Println("root hash is ", tree.RootHash())
}
