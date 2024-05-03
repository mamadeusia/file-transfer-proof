package merkletree

import (
	"fmt"
	"strings"
)

func PrintTree(node Node) {
	if node.right == node.left { //for the case that size of addresses==1
		PrintRoot(node)
	} else {
		printNode(node, 0)
	}
}

func PrintRoot(node Node) {
	if node.right == node.left { //for the case that size of addresses==1
		fmt.Printf("0x%s \n", node.right.hash())
	} else {
		fmt.Printf("0x%s \n", node.hash())
	}
}

func RootToString(node Node) string {
	if node.right == node.left { //for the case that size of addresses==1
		return node.right.hash().String()
	} else {
		return node.hash().String()
	}
}

func PrintfProof(proofs []Hashable) {
	fmt.Printf("[")
	for i, proof := range proofs {
		fmt.Printf("\"0x%s\"", proof.hash())
		if i != len(proofs)-1 {
			fmt.Printf(",")
		}
	}
	fmt.Printf("]\n")

}

func ProofToStringSlice(proofs []Hashable) []string {
	var output []string
	for _, proof := range proofs {
		output = append(output, proof.hash().String())
	}
	return output
}

func printNode(node Node, level int) {
	fmt.Printf("(%d) %s %s\n", level, strings.Repeat("-", level), node.hash())
	if l, ok := node.left.(Node); ok {
		printNode(l, level+1)
	} else if l, ok := node.left.(leaf); ok {
		fmt.Printf("(%d) %s %s (data: %s)\n", level+1, strings.Repeat("-", level+1), l.hash(), l)
	}
	if r, ok := node.right.(Node); ok {
		printNode(r, level+1)
	} else if r, ok := node.right.(leaf); ok {
		fmt.Printf("(%d) %s %s (data: %s)\n", level+1, strings.Repeat("-", level+1), r.hash(), r)
	}
}
