package diagram

import "sort"

// layerNodes assigns the direct nodes of a container to columns. Only edges whose
// two ends are both in nodes count. A dependency sits in a lower column than its
// dependents; cycles are broken by dropping back edges from layering only. The
// result is ordered alphabetically, then by two barycenter sweeps.
func layerNodes(nodes []*Node, edges []Edge) [][]*Node {
	sorted := append([]*Node(nil), nodes...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Address < sorted[j].Address })
	in := make(map[*Node]bool, len(sorted))
	for _, n := range sorted {
		in[n] = true
	}
	deps := make(map[*Node][]*Node) // dependent -> dependencies
	for _, e := range edges {
		if e.To != nil && in[e.From] && in[e.To] && e.From != e.To {
			deps[e.From] = append(deps[e.From], e.To)
		}
	}
	for _, list := range deps {
		sort.Slice(list, func(i, j int) bool { return list[i].Address < list[j].Address })
	}
	kept := breakCycles(sorted, deps)
	layer := make(map[*Node]int, len(sorted))
	var depth func(n *Node) int
	depth = func(n *Node) int {
		if l, ok := layer[n]; ok {
			return l
		}
		l := 0
		for _, d := range kept[n] {
			if dl := depth(d) + 1; dl > l {
				l = dl
			}
		}
		layer[n] = l
		return l
	}
	maxLayer := 0
	for _, n := range sorted {
		if l := depth(n); l > maxLayer {
			maxLayer = l
		}
	}
	cols := make([][]*Node, maxLayer+1)
	for _, n := range sorted {
		cols[layer[n]] = append(cols[layer[n]], n)
		n.Layer = layer[n]
	}
	orderLayers(cols, deps)
	return cols
}

// breakCycles returns the dependencies without the back edges found by a DFS over
// the nodes sorted by address.
func breakCycles(sorted []*Node, deps map[*Node][]*Node) map[*Node][]*Node {
	const (
		unseen = iota
		onStack
		done
	)
	state := make(map[*Node]int, len(sorted))
	kept := make(map[*Node][]*Node, len(sorted))
	var visit func(n *Node)
	visit = func(n *Node) {
		state[n] = onStack
		for _, d := range deps[n] {
			switch state[d] {
			case onStack:
				continue // back edge: not used for layering
			case unseen:
				kept[n] = append(kept[n], d)
				visit(d)
			default:
				kept[n] = append(kept[n], d)
			}
		}
		state[n] = done
	}
	for _, n := range sorted {
		if state[n] == unseen {
			visit(n)
		}
	}
	return kept
}

// orderLayers runs two barycenter sweeps (forward, then backward) with a stable
// sort and an address tie-break. Barycenters are compared as integer fractions.
func orderLayers(cols [][]*Node, deps map[*Node][]*Node) {
	adj := make(map[*Node][]*Node)
	for from, list := range deps {
		for _, to := range list {
			adj[from] = append(adj[from], to)
			adj[to] = append(adj[to], from)
		}
	}
	for i := 1; i < len(cols); i++ {
		sweep(cols[i], cols[i-1], adj)
	}
	for i := len(cols) - 2; i >= 0; i-- {
		sweep(cols[i], cols[i+1], adj)
	}
}

type bary struct{ sum, count int }

// sweep reorders one layer by the average position of its neighbors in the
// reference layer; nodes without neighbors keep their current position.
func sweep(layer, ref []*Node, adj map[*Node][]*Node) {
	pos := make(map[*Node]int, len(ref))
	for i, n := range ref {
		pos[n] = i
	}
	keys := make(map[*Node]bary, len(layer))
	for i, n := range layer {
		b := bary{}
		for _, m := range adj[n] {
			if p, ok := pos[m]; ok {
				b.sum += p
				b.count++
			}
		}
		if b.count == 0 {
			b = bary{sum: i, count: 1}
		}
		keys[n] = b
	}
	sort.SliceStable(layer, func(i, j int) bool {
		a, b := keys[layer[i]], keys[layer[j]]
		if l, r := a.sum*b.count, b.sum*a.count; l != r {
			return l < r
		}
		return layer[i].Address < layer[j].Address
	})
}
