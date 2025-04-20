package graph

type dfsStates int

const (
	UNVISITED dfsStates = iota
	VISITING
	VISITED
)

func dfsVisit[T any](
	n *Node[T],
	state map[*Node[T]]dfsStates,
	parent map[*Node[T]]*Node[T],
	cycle *[]*Node[T],
) bool {
	state[n] = VISITING // Visiting

	for _, neighbour := range n.GetOutgoingNodes() {
		switch state[neighbour] {
		case UNVISITED: // Unvisited
			parent[neighbour] = n
			if dfsVisit(neighbour, state, parent, cycle) {
				return true
			}
		case VISITING: // Visiting => ciclo detectado
			curr := n
			*cycle = []*Node[T]{neighbour}
			for curr != neighbour {
				*cycle = append(*cycle, curr)
				curr = parent[curr]
			}
			*cycle = append(*cycle, neighbour) // cerrar ciclo
			return true
		}
	}

	state[n] = VISITED // Visited
	return false
}
