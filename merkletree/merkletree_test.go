package merkletree

import (
	"bytes"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"math/rand"
	"testing"
)

func HashFuncSHA256(data []byte) ([]byte, error) {
	hash := sha256.New()
	hash.Write(data)
	hashValue := hash.Sum(nil)
	//fmt.Println("Calculated hash len is :", len(hashValue))
	return hashValue, nil
}

func HashFuncSHA512(data []byte) ([]byte, error) {
	hash := sha512.New()
	hash.Write(data)
	hashValue := hash.Sum(nil)
	//fmt.Println("Calculated hash len is :", len(hashValue))
	return hashValue, nil
}

func TestMerkleBasics(t *testing.T) {
	var data Data
	expectedRootHash := []byte{95, 48, 204, 128, 19, 59, 147, 148, 21, 110, 36, 178, 51, 240, 196, 190, 50, 178, 78, 68, 187, 51, 129, 240, 44, 123, 165, 38, 25, 208, 254, 188}
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

	tree, err := NewTree(&data, HashFuncSHA256)
	if err != nil {
		fmt.Println("Failed to build Tree due to error ", err)
		t.FailNow()
	}
	fmt.Printf("Tree Root:  %v \n", tree.RootHash())
	if !bytes.Equal(tree.RootHash(), expectedRootHash) {
		fmt.Println("Root hash generated is not matching expected hash")
		t.FailNow()
	}
	/* 	err = tree.AddLeaf([]byte{'t', 'e', 's', 't'})
	   	if err != nil {
	   		fmt.Println("Failed to add a leaf Tree due to error ", err)
	   		t.FailNow()
	   	}
	   	fmt.Printf("New Tree Root:  %v \n", tree.RootHash()) */

	proof, err := tree.GenerateMerkleProof(&data[0])
	if err != nil {
		fmt.Println("Failed to generate Merkle Proof due to error ", err)
		t.FailNow()
	}
	fmt.Println("Generated Proof is ", *proof)
	expectedProof := Proof{
		Hashes: [][]byte{
			{54, 57, 239, 205, 8, 171, 178, 115, 177, 97, 158, 130, 231, 140, 41, 167, 223, 2, 193, 5, 27, 24, 32, 233, 159, 195, 149, 220, 170, 51, 38, 184},
			{103, 184, 144, 26, 195, 1, 53, 231, 77, 66, 3, 109, 250, 96, 67, 54, 225, 249, 120, 228, 158, 224, 214, 191, 72, 74, 70, 255, 39, 162, 174, 156}},
		Indexes: []int{1, 1},
	}
	match, err := proof.Equals(&expectedProof)
	if !match {
		fmt.Println("Generated Merkle Proof doesn't match expected proof due to error ", err)
		t.FailNow()
	}
	//fmt.Println("Proof generated is : \n", proof)

	verified, err := tree.VerifyProof(&data[0], proof)
	if err != nil {
		fmt.Println("Failed to verify Merkle Proof due to error ", err)
		t.FailNow()
	}
	if !verified {
		fmt.Println("Invalid proof for the data", err)
		t.FailNow()
	}
	expectedProof = Proof{
		Hashes: [][]byte{
			{88, 29, 67, 116, 87, 38, 224, 238, 98, 145, 17, 120, 191, 179, 136, 124, 63, 226, 149, 210, 158, 235, 116, 31, 14, 64, 249, 30, 138, 112, 144, 122},
			{103, 184, 144, 26, 195, 1, 53, 231, 77, 66, 3, 109, 250, 96, 67, 54, 225, 249, 120, 228, 158, 224, 214, 191, 72, 74, 70, 255, 39, 162, 174, 156},
			{95, 48, 204, 128, 19, 59, 147, 148, 21, 110, 36, 178, 51, 240, 196, 190, 50, 178, 78, 68, 187, 51, 129, 240, 44, 123, 165, 38, 25, 208, 254, 188}},
		Indexes: []int{0, 1, 0},
	}
	proof, err = tree.GetMerklePath(&data[2])
	if err != nil {
		fmt.Println("Failed to generate Merkle Path due to error ", err)
		t.FailNow()
	}
	match, err = proof.Equals(&expectedProof)
	if !match {
		fmt.Println("Generated Merkle Path doesn't match expected Path due to error ", err)
		t.FailNow()
	}
	//fmt.Println("Merkle Path is : \n", proof)

	err = tree.UpdateLeaf(&data[0], &data[1])
	if err != nil {
		fmt.Printf("Failed to Update Leaf node with oldValue %v due to error %v\n", data[0], err)
		t.FailNow()
	}
	expectedRootHash = []byte{65, 94, 211, 81, 47, 69, 136, 19, 206, 251, 153, 39, 235, 99, 159, 208, 220, 46, 32, 181, 213, 210, 117, 140, 11, 114, 70, 5, 49, 140, 135, 45}
	if !bytes.Equal(tree.RootHash(), expectedRootHash) {
		fmt.Println("Root hash generated is not matching expected hash after leaf is updated")
		t.FailNow()
	}
	fmt.Printf("Updated Tree Root: %v\n", tree.RootHash())

}

func TestMerkleTreeAdvanced(t *testing.T) {
	/*TODO:
	1. Test with different hash functions
	2. Verify all proofs and paths for a tree
	3. Verify tree construction, proof generation and verification with large data set e.g : 10,000 leaves
	4.
	*/
}

func TestMerkleTreeNegativeScenarios(t *testing.T) {
	/*TODO:
	1. Test for non existence of data
	2. Test with passing wrong proofs for data
	3. Duplicate trees (sam root hash) by passing data list with duplicates. Example: https://github.com/bitcoin/bitcoin/blob/master/src/consensus/merkle.cpp
	4. Test by using hash function of lower security
	5.
	*/
}

func TestMerkleTreeRareScenarios(t *testing.T) {
	/*TODO:
	1. Test to verify tree depth and identify duplications for large number of nodes
	2.
	*/
}

var tree *MerkleTree
var data Data

func merkleTreeProofGen(leafCount int, b *testing.B) {

	//b.ResetTimer()
	rand := rand.Intn(leafCount - 1)
	_, err := tree.GenerateMerkleProof(&data[rand])
	if err != nil {
		//fmt.Println("Failed to generate Merkle Proof due to error ", err)
		b.FailNow()
	}
}

func benchmarkMerkleTreeProofGeneration(leafCount int, b *testing.B) {
	treeCreate(leafCount)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		merkleTreeProofGen(100, b)
	}
}

func BenchmarkMerkleTreeProof100(b *testing.B) {
	benchmarkMerkleTreeProofGeneration(100, b)
}

func BenchmarkMerkleTreeProof1000(b *testing.B) {
	benchmarkMerkleTreeProofGeneration(1000, b)
}

func BenchmarkMerkleTreeProof5000(b *testing.B) {
	benchmarkMerkleTreeProofGeneration(5000, b)
}

func BenchmarkMerkleTreeBuild100(b *testing.B) {
	for i := 0; i < b.N; i++ {
		treeCreate(100)
	}
}

func BenchmarkMerkleTreeBuild1000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		treeCreate(1000)
	}
}

func BenchmarkMerkleTreeBuild5000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		treeCreate(5000)
	}
}

func treeCreate(count int) {
	var err error
	str := "Hello"
	for i := 0; i < count; i++ {
		rand := rand.Int()
		d := []byte(fmt.Sprintf("%s%d", str, rand))
		data = append(data, d)
	}

	tree, err = NewTree(&data, HashFuncSHA256)
	if err != nil {
		fmt.Println("Failed to build Tree due to error ", err)
	}
}
