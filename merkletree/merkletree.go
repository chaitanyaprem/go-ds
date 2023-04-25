package merkletree

import (
	"bytes"
	"fmt"
	"math"
	// "github.com/m1gwings/treedrawer/tree"
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

type Proof struct {
	Hashes  [][]byte
	Indexes []int
}

func (p *Proof) Equals(p1 *Proof) (bool, error) {
	for i, _ := range p.Hashes {
		if !bytes.Equal(p.Hashes[i], p1.Hashes[i]) {
			return false, fmt.Errorf("proof mismatch at %d for hash", i)
		}
		if p.Indexes[i] != p1.Indexes[i] {
			return false, fmt.Errorf("proof mismatch at %d for the index", i)
		}
	}
	return true, nil
}

func (n *Node) String() string {
	fmtNode := fmt.Sprintf("Node : {Hash: %v, parent %v, IsLeaf:%t, IsDuplicate: %t}\n",
		n.Hash, n.Parent, n.IsLeaf, n.IsDuplicate)
	return fmtNode
}

type Data [][]byte

func NewTree(data Data, hashFunc HashFunction) (*MerkleTree, error) {
	var tree MerkleTree
	tree.HashFunc = hashFunc

	var err error
	leafCount := len(data)
	//fmt.Println("Number of leaves:", leafCount)
	if leafCount == 0 {
		return nil, fmt.Errorf("error: cannot build a merkle tree without any data")
	}
	err = populateLeaves(data, &tree)
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

func populateLeaves(data Data, tree *MerkleTree) error {
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
	node.Hash, err = tree.HashFunc(node.Data)
	if err != nil {
		return nil, err
	}
	node.IsLeaf = true
	return &node, nil
}

func buildIntermediateLevel(nodes []*Node, tree *MerkleTree) (*Node, error) {
	//var levelNodes []*Node
	levelNodes := make([]*Node, len(nodes)/2)
	levelIndex := 0
	for j := 0; j < len(nodes); j = j + 2 {
		node, err := createNonLeafNode(nodes[j], nodes[j+1], tree.HashFunc)
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

func (t *MerkleTree) UpdateLeaf(oldValue []byte, newValue []byte) error {
	oldHash, err := t.HashFunc(oldValue)
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
	newHash, err := t.HashFunc(newValue)
	if err != nil {
		return err
	}
	t.Leaves[leafIndex].Hash = newHash
	t.Leaves[leafIndex].Data = newValue
	lastLeafIndex := len(t.Leaves)
	if leafIndex == lastLeafIndex-1 && t.Leaves[lastLeafIndex].IsDuplicate {
		//Update if duplicate node is present.
		t.Leaves[lastLeafIndex].Hash = newHash
		t.Leaves[lastLeafIndex].Data = newValue
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

func (t *MerkleTree) GetMerklePath(data []byte) (*Proof, error) {
	dHash, err := t.HashFunc(data)
	if err != nil {
		return nil, err
	}
	var proof Proof
	proof.Hashes = make([][]byte, t.Depth+1)
	proof.Indexes = make([]int, t.Depth+1)
	for _, node := range t.Leaves {
		if bytes.Equal(node.Hash, dHash) {
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
	return &proof, nil
}

func (t *MerkleTree) GenerateMerkleProof(data []byte) (*Proof, error) {
	dHash, err := t.HashFunc(data)
	if err != nil {
		return nil, err
	}
	var proof Proof
	//TODO: Optimize
	//proof.Hashes = make([][]byte, t.Depth)
	//Traverse the tree upwards and get hashes and indexes required for the proof.
	for _, node := range t.Leaves {
		if bytes.Equal(node.Hash, dHash) {
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
	//fmt.Printf("Generated Proof for has %v is %v \n", dHash, proof)
	return &proof, nil
}

func (t *MerkleTree) VerifyProof(data []byte, proof *Proof) (bool, error) {
	//fmt.Println("VerifyProof:To be implemented")
	dHash, err := t.HashFunc(data)
	if err != nil {
		return false, err
	}
	for i, val := range proof.Hashes {

		if proof.Indexes[i] == 0 {
			dHash, err = t.HashFunc(append(val, dHash...))
			if err != nil {
				return false, err
			}
		} else {
			dHash, err = t.HashFunc(append(dHash, val...))
			if err != nil {
				return false, err
			}
		}
	}
	// fmt.Println("generated rootHash is ", dHash)
	// fmt.Println("Tree's rootHash is ", t.RootHash())

	if !bytes.Equal(dHash, t.RootHash()) {
		return false, fmt.Errorf("proof verification failed due to mismatch in generated root hash")
	}
	return true, nil
}

func (t *MerkleTree) VerifyTree() (bool, error) {
	fmt.Println("VerifyTree:Verifying proofs of all leaves")
	//TODO: Generate rootHash and verifying it against stored root.
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

/* func (n *Node) Draw() *drawer.Drawer {
	drawer.NewDrawer()
	n.Hash
} */

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
