package merkletree

import (
	"bytes"
	"fmt"
	"math"
	// "github.com/m1gwings/treedrawer/tree"
)

/*
Defines the hashfunction type that needs to be passed while building a merkle tree.
Takes byte array of data to be hashed as an argument
Returns
  - Hash of the data as a byte array
  - error in case of any errors.
*/
type HashFunction func(data []byte) ([]byte, error)

type MerkleTree struct {
	Root     *Node
	Leaves   []*Node
	HashFunc HashFunction
	Depth    int
}

type Data [][]byte

/*
Function verifies the security of the hash function passed based on the lenght of the hash.
If the lenght of the hash is lesser than 128 bits, an error is returned.
*/
func VerifyHashFuncMinSecurity(hashFunc HashFunction) error {
	//Limit security of hash
	testHash, err := hashFunc([]byte{1, 2, 3, 4, 5, 6, 7})
	if err != nil {
		return err
	}
	if len(testHash) < 16 {
		return fmt.Errorf("hash function is not secure enough, require a min 128 bit output to be generated")
	}
	return nil
}

/*
Builds a new merkle tree from list of data and using the hashfunction that is passed.
Note that hash output of the hashFunction should be a min of 128 bits, otherwise it is considered insecure.
Accepts
  - data list to be used for building the tree
  - Hash function to be used for hashing

Returns
  - Reference to the tree in case of no errors
  - error detailing cause of errror while building the tree
*/
func NewTree(data *Data, hashFunc HashFunction) (*MerkleTree, error) {
	if err := VerifyHashFuncMinSecurity(hashFunc); err != nil {
		return nil, err
	}
	var tree MerkleTree
	tree.HashFunc = hashFunc

	leafCount := len(*data)
	//fmt.Println("Number of leaves:", leafCount)
	if leafCount == 0 {
		return nil, fmt.Errorf("error: cannot build a merkle tree without any data")
	}
	err := populateLeaves(data, &tree)
	if err != nil {
		return nil, err
	}
	leafCount = len(tree.Leaves)
	//Calculate Tree depth based on number of leaves considering it is a binary hash tree
	tree.Depth = int(math.Log2(float64(leafCount)))
	tree.Root, err = buildIntermediateLevel(tree.Leaves, &tree)
	if err != nil {
		return nil, err
	}
	return &tree, nil
}

func populateLeaves(data *Data, tree *MerkleTree) error {
	leafCount := len(*data)

	//TODO : sort leaves and then build tree so that same root hash is generated even if data order is different.
	//Create Leaf nodes
	for _, val := range *data {
		node, err := buildLeafNode(val, nil, tree.HashFunc)
		if err != nil {
			return err
		}

		tree.Leaves = append(tree.Leaves, node)
	}

	if leafCount%2 == 1 {
		// Handle case of odd leaves. Create a duplicate leaf node
		node, err := buildLeafNode(tree.Leaves[leafCount-1].Data, tree.Leaves[leafCount-1].Hash, tree.HashFunc)
		node.IsDuplicate = true
		if err != nil {
			return err
		}
		tree.Leaves = append(tree.Leaves, node)
	}

	return nil
}

func buildLeafNode(data []byte, hash []byte, hashFunction HashFunction) (*Node, error) {
	var node Node
	var err error
	node.Data = data
	if hash == nil {
		node.Hash, err = hashFunction(node.Data)
	} else {
		node.Hash = hash
	}
	if err != nil {
		return nil, err
	}
	node.IsLeaf = true
	return &node, nil
}

func buildIntermediateLevel(nodes []*Node, tree *MerkleTree) (*Node, error) {
	//var levelNodes []*Node
	levelCount := len(nodes) / 2
	if len(nodes)%2 == 1 {
		levelCount = len(nodes)/2 + 1
	}
	levelNodes := make([]*Node, levelCount)
	levelIndex := 0

	//TODO:This recursion be optimized to run in parallel if there are too many leaves? Maybe breakdown chunks of subtrees build them.
	for j := 0; j < len(nodes); j = j + 2 {
		left := j
		right := j + 1
		//Possible vulnerability as explained here https://github.com/bitcoin/bitcoin/blob/master/src/consensus/merkle.cpp
		//This should not effect a tree where each data element is expected to be unique.
		//TODO: Can be optimized to take odd leaf to next level and not duplicate hash.
		if right == len(nodes) {
			right = j
		}
		node, err := createNonLeafNode(nodes[left], nodes[right], tree.HashFunc)
		if err != nil {
			return nil, err
		}
		levelNodes[levelIndex] = node
		levelIndex++
		//levelNodes = append(levelNodes, node)
		if len(nodes) == 2 {
			return node, nil
		}
	}
	return buildIntermediateLevel(levelNodes, tree)
}

func createNonLeafNode(left *Node, right *Node, hashFunction HashFunction) (*Node, error) {
	var node Node
	var err error
	node.Left = left
	node.Right = right
	node.Hash, err = hashFunction(append(left.Hash, right.Hash...))
	if err != nil {
		return nil, err
	}
	left.Parent = &node
	right.Parent = &node
	return &node, nil
}

func (t *MerkleTree) RootHash() []byte {
	return t.Root.Hash
}

func (t *MerkleTree) UpdateLeaf(oldValue *[]byte, newValue *[]byte) error {
	oldHash, err := t.HashFunc(*oldValue)
	if err != nil {
		return err
	}
	leafIndex := -1
	for i, node := range t.Leaves {
		if bytes.Equal(node.Hash, oldHash) {
			leafIndex = i
			break
		}
	}
	if leafIndex == -1 {
		return fmt.Errorf("could not find leaf node to update")
	}
	parent := t.Leaves[leafIndex].Parent
	newHash, err := t.HashFunc(*newValue)
	if err != nil {
		return err
	}
	t.Leaves[leafIndex].Hash = newHash
	t.Leaves[leafIndex].Data = *newValue
	lastLeafIndex := len(t.Leaves)
	if leafIndex == lastLeafIndex-1 && t.Leaves[lastLeafIndex].IsDuplicate {
		//Update if duplicate node is present.
		t.Leaves[lastLeafIndex].Hash = newHash
		t.Leaves[lastLeafIndex].Data = *newValue
		t.Leaves[lastLeafIndex].IsDuplicate = false
	}
	for parent != nil {
		if bytes.Equal(parent.Left.Hash, newHash) { //Left Node
			parent.Hash, err = t.HashFunc(append(newHash, parent.Right.Hash...))
		} else { //Right Node
			parent.Hash, err = t.HashFunc(append(parent.Left.Hash, newHash...))
		}
		if err != nil {
			return err
		}
		if parent.Parent == nil {
			break
		}
		parent = parent.Parent
		newHash = parent.Hash
	}

	return nil
}

/*
Generates merklePath for data if it is present in the tree.
MerklePath is a proof object containing list of hashes while traversing from the leaf to the root.
Accepts
  - data for which merkle path has to be generated

Returns
  - Reference to a Proof object
  - error in case any error occurs or data not found
*/
func (t *MerkleTree) GetMerklePath(data *[]byte) (*Proof, error) {
	leafFound := false
	dHash, err := t.HashFunc(*data)
	if err != nil {
		return nil, err
	}
	var proof Proof
	proof.Hashes = make([][]byte, t.Depth+1)
	proof.Indexes = make([]int, t.Depth+1)
	for _, node := range t.Leaves {
		if bytes.Equal(node.Hash, dHash) {
			leafFound = true
			curHash := dHash
			parent := node.Parent
			for index := 0; parent != nil; index++ {
				proof.Hashes[index] = curHash
				if !bytes.Equal(parent.Left.Hash, curHash) { //Right Node
					proof.Indexes[index] = 1
				}
				curHash = parent.Hash
				parent = parent.Parent
			}
			proof.Hashes[t.Depth] = t.RootHash()
			break
		}
	}
	if !leafFound {
		return nil, fmt.Errorf("data doesn't exist in the tree")
	}
	return &proof, nil
}

/*
Generates merkleProof for data if it is present in the tree.
Proof object containing list of sibling node hashes while traversing from the leaf to the root.
Accepts
  - data for which merkle proof has to be generated

Returns
  - Reference to a Proof object
  - error in case any error occurs or data not found
*/
func (t *MerkleTree) GenerateMerkleProof(data *[]byte) (*Proof, error) {
	leafFound := false

	dHash, err := t.HashFunc(*data)
	if err != nil {
		return nil, err
	}
	var proof Proof
	//TODO: Optimize
	//proof.Hashes = make([][]byte, t.Depth)
	//Traverse the tree upwards and get hashes and indexes required for the proof.
	for _, node := range t.Leaves {
		if bytes.Equal(node.Hash, dHash) {
			leafFound = true
			curHash := dHash
			parent := node.Parent
			for parent != nil {
				if bytes.Equal(parent.Left.Hash, curHash) { //Left Node
					proof.Hashes = append(proof.Hashes, parent.Right.Hash)
					proof.Indexes = append(proof.Indexes, 1)
				} else { //Right Node
					proof.Hashes = append(proof.Hashes, parent.Left.Hash)
					proof.Indexes = append(proof.Indexes, 0)
				}
				curHash = parent.Hash
				parent = parent.Parent
			}
			break
		}
	}
	if !leafFound {
		return nil, fmt.Errorf("data doesn't exist in the tree")
	}
	//fmt.Printf("Generated Proof for has %v is %v \n", dHash, proof)
	return &proof, nil
}

/*
Verifies if the proof for the data is valid or not
Accepts
  - the data for which proof is generated
  - Proof object indicating merkle proof

Returns
  - true if proof is valid
  - error in case of errors
*/
func (tree *MerkleTree) VerifyProof(data *[]byte, proof *Proof) (bool, error) {
	//fmt.Println("VerifyProof:To be implemented")
	dHash, err := tree.HashFunc(*data)
	if err != nil {
		return false, err
	}
	for i, val := range proof.Hashes {

		if proof.Indexes[i] == 0 {
			dHash, err = tree.HashFunc(append(val, dHash...))
			if err != nil {
				return false, err
			}
		} else {
			dHash, err = tree.HashFunc(append(dHash, val...))
			if err != nil {
				return false, err
			}
		}
	}
	// fmt.Println("generated rootHash is ", dHash)
	// fmt.Println("Tree's rootHash is ", t.RootHash())

	if !bytes.Equal(dHash, tree.RootHash()) {
		return false, fmt.Errorf("generated root hash not matches stored root")
	}
	return true, nil
}

func (m *MerkleTree) String() string {
	fmtTree := ""
	//TODO: Print the entire tree structure with all nodes.
	/* 	node := m.Root
	   	for node != nil {
	   		fmtTree += node.String()
	   		node = node.Left
	   	} */
	leaves := ""
	for i, val := range m.Leaves {
		leaves += fmt.Sprintf("%d,{%v}\n", i, val)
	}
	fmtTree = fmt.Sprintf("Depth: %d, \nRoot:%v,\n Leaves:%s", m.Depth, m.Root, leaves)
	return fmtTree
}

/* func (m *MerkleTree) PrettyPrint() {
	node := m.Root
	t := tree.NewTree(tree.NodeString(node.Hash))
	temp := t
	var left, right *Node
	for node.Left != nil {
		temp.AddChild(tree.NodeString(node.Left.Hash))
		temp.AddChild(tree.NodeString(node.Right.Hash))
	}
}
*/
