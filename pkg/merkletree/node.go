package merkletree

import (
	"bytes"
	"encoding/hex"
	"errors"
)

type leaf string

func hashableLeafHashFromString(s string) (Hashable, error) {
	data, err := hex.DecodeString(s)
	if err != nil || len(data) != 32 {
		return nil, errors.New("wrong format, string is not 32 byte hash")
	}
	return leafHash(s), nil
}

type leafHash string

type Node struct {
	left        Hashable
	right       Hashable
	inProofTree bool
}

// TODO :: change the structure of this function , Note : because input is forced to be 32 byte it's true
func (b leaf) hash() Hash {
	data, _ := hex.DecodeString(string(b))
	return hash([]byte(data)[:])
}

// TODO :: change the output, add error, prevent for other than 32byte hex string.
func (b leafHash) hash() Hash {
	data, err := hex.DecodeString(string(b))
	if err != nil {
		return [32]byte{}
	}
	var output Hash
	copy(output[:], data[0:32])
	return output
}

func (n Node) hash() Hash {
	var l, r [32]byte
	l = n.left.hash()
	r = n.right.hash()
	return hash(append(l[:], r[:]...))
}

func NewNode(left Hashable, right Hashable) Node {
	leftHash := left.hash()
	rightHash := right.hash()
	g := bytes.Compare(leftHash[:], rightHash[:])
	if g <= 0 {
		return Node{left: left, right: right}
	} else {
		return Node{left: right, right: left}
	}
}

// type LeafSorted []Leaf

// func (l LeafSorted) Len() int { return len(l) }

// func (l LeafSorted) Less(i, j int) bool {
// 	left := l[i].hash()
// 	right := l[j].hash()
// 	g := bytes.Compare(left[:], right[:])
// 	if g <= 0 {
// 		return true
// 	} else {
// 		return false
// 	}
// }

// func (l LeafSorted) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
