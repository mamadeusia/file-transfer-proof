package merkletree

func BuildTreeWithLeafHashes(leafHashse []string) Node {
	var nodes []Hashable
	for i := 0; i < len(leafHashse); i += 2 {
		if i+1 < len(leafHashse) {
			nodes = append(nodes, NewNode(leafHash(leafHashse[i]), leafHash(leafHashse[i+1])))
		} else {
			nodes = append(nodes, NewNode(leafHash(leafHashse[i]), leafHash(leafHashse[i])))
		}
	}
	var output []Hashable
	if len(nodes) == 1 {
		return nodes[0].(Node)
	} else {
		output = buildTree(nodes)
	}
	return output[0].(Node)

}

func buildTree(parts []Hashable) []Hashable {

	var nodes []Hashable
	var i int
	for i = 0; i < len(parts); i += 2 {
		if i+1 < len(parts) {
			nodes = append(nodes, NewNode(parts[i], parts[i+1]))
		} else {
			nodes = append(nodes, NewNode(parts[i], parts[i]))
		}
	}
	if len(nodes) == 1 {
		return nodes
	} else if len(nodes) > 1 {
		return buildTree(nodes)
	} else {
		panic("huh?!")
	}
}
