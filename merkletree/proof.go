package merkletree

import (
	"bytes"
	"fmt"
)

type Proof struct {
	Hashes  [][]byte
	Indexes []int
}

func (p *Proof) Equals(p1 *Proof) (bool, error) {
	for i := range p.Hashes {
		if !bytes.Equal(p.Hashes[i], p1.Hashes[i]) {
			return false, fmt.Errorf("proof mismatch at %d for hash", i)
		}
		if p.Indexes[i] != p1.Indexes[i] {
			return false, fmt.Errorf("proof mismatch at %d for the index", i)
		}
	}
	return true, nil
}
