package merkletree

import "fmt"

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
	fmtNode := fmt.Sprintf("Node : {Hash: %v, parent %v, IsLeaf:%t, IsDuplicate: %t}\n",
		n.Hash, n.Parent, n.IsLeaf, n.IsDuplicate)
	return fmtNode
}

/* func (n *Node) Draw() *drawer.Drawer {
	drawer.NewDrawer()
	n.Hash
} */
