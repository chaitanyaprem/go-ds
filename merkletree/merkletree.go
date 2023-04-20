package merkletree

import (
	"fmt"
	"math"
)

type HashFunction func(data []byte) ([]byte, error)

type MerkleTree struct {
	Root     *Node
	Leaves   []*Node
	HashFunc HashFunction
	Depth    int
}

type Node struct {
	Hash        []byte
	Parent      *Node
	Left        *Node
	Right       *Node
	IsLeaf      bool
	IsDuplicate bool
	Data        []byte
	//TODO: Add Path metadata
}

func (n *Node) String() string {
	fmtNode := ""

	return fmtNode
}

func NewTree(data [][]byte, hashFunc HashFunction) (*MerkleTree, error) {
	var tree MerkleTree
	tree.HashFunc = hashFunc
	var err error
	leafCount := len(data)
	fmt.Println(leafCount, " is leaf node count")
	if leafCount == 0 {
		return nil, fmt.Errorf("error: cannot build a merkle tree without any data")
	}
	err = populateLeaves(data, &tree)
	if err != nil {
		return nil, err
	}
	//Calculate Tree depth based on number of leaves considering it is a binary hash tree
	tree.Depth = int(math.Log2(float64(leafCount)))
	tree.Root, err = buildIntermediateLevel(tree.Leaves, &tree)
	if err != nil {
		return nil, err
	}
	return &tree, nil
}

func populateLeaves(data [][]byte, tree *MerkleTree) error {
	leafCount := len(data)

	//TODO : sort leaves and then build tree so that same root hash is generated even if data order is different.
	//Create Leaf nodes
	for _, val := range data {
		node, err := buildLeafNode(val, tree)
		if err != nil {
			return err
		}
		tree.Leaves = append(tree.Leaves, node)
	}

	if leafCount%2 == 1 {
		// Handle case of odd leaves. Create a null leaf node
		node, err := buildLeafNode(tree.Leaves[leafCount-1].Data, tree)
		node.IsDuplicate = true
		if err != nil {
			return err
		}
		tree.Leaves = append(tree.Leaves, node)
	}
	return nil
}

func buildLeafNode(data []byte, tree *MerkleTree) (*Node, error) {
	var node Node
	var err error
	node.Data = data
	//TODO: What if hashfunction returns error.
	node.Hash, err = tree.HashFunc(node.Data)
	if err != nil {
		return nil, err
	}
	node.IsLeaf = true
	return &node, nil
}

func buildIntermediateLevel(nodes []*Node, tree *MerkleTree) (*Node, error) {
	var levelNodes []*Node
	for j := 0; j < len(nodes); j = j + 2 {
		node, err := createNonLeafNode(tree.Leaves[j], tree.Leaves[j+1], tree.HashFunc)
		if err != nil {
			return nil, err
		}
		levelNodes = append(levelNodes, node)
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

func (t *MerkleTree) GenerateMerkleProof() {
	fmt.Println("GenerateMerkleProof:To be implemented")

}

func (t *MerkleTree) VerifyData(data []byte) (bool, error) {
	fmt.Println("VerifyData:To be implemented")
	return true, nil
}

func (t *MerkleTree) VerifyTree() (bool, error) {
	fmt.Println("VerifyTree:To be implemented")

	return true, nil
}

func (m *MerkleTree) String() string {
	fmtTree := ""
	//TODO: Traverse the tree and invoke string on each node.
	node := m.Root
	for node != nil {
		fmtTree += node.String()
		node = node.Left
	}
	return fmtTree
}
