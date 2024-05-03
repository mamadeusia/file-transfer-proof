package merkletree

func GetProofIndexWithLeafHashes(j int, leafHashes []string) []Hashable {

	if len(leafHashes) == 1 { //for the case on one leaf
		return []Hashable{}
	}

	var output []Hashable
	var nodes []Node
	for i := 0; i < len(leafHashes); i += 2 {
		if i+1 < len(leafHashes) {
			node := NewNode(leafHash(leafHashes[i]), leafHash(leafHashes[i+1]))
			if j == i {
				output = append(output, leafHash(leafHashes[i+1]))
				node.inProofTree = true
			} else if j == i+1 {
				output = append(output, leafHash(leafHashes[i]))
				node.inProofTree = true
			}
			nodes = append(nodes, node)
		} else {
			node := NewNode(leafHash(leafHashes[i]), leafHash(leafHashes[i]))
			if j == i {
				output = append(output, leafHash(leafHashes[i]))
				node.inProofTree = true
			}
			nodes = append(nodes, node)
		}
	}
	if len(nodes) == 1 {
		return []Hashable{output[0]}
	}
	return appendProof(nodes, output)

}

func appendProof(parts []Node, input []Hashable) []Hashable {

	output := input
	var newParts []Node
	for i := 0; i < len(parts); i += 2 {
		if i+1 < len(parts) {
			node := NewNode(parts[i], parts[i+1])
			if parts[i].inProofTree {
				output = append(output, parts[i+1])
				node.inProofTree = true
			} else if parts[i+1].inProofTree {
				output = append(output, parts[i])
				node.inProofTree = true
			}
			newParts = append(newParts, node)
		} else {
			node := NewNode(parts[i], parts[i])
			if parts[i].inProofTree {
				output = append(output, parts[i])
				node.inProofTree = true
			}
			newParts = append(newParts, node)
		}
	}
	if len(newParts) == 1 {
		return output
	} else if len(newParts) > 1 {
		return appendProof(newParts, output)
	} else {
		panic("huh?!")
	}
}

func CheckProofWithLeafHash(rootHash string, leafHash string, proofs []string) bool {
	lastNode, err := hashableLeafHashFromString(leafHash)
	if err != nil {
		return false
	}
	var node Node
	for _, proof := range proofs {
		n, err := hashableLeafHashFromString(proof)
		if err != nil {
			return false
		}
		node = NewNode(lastNode, n)
		lastNode = node
	}

	return RootToString(node) == rootHash

}
