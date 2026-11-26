package diagram

// Layout metrics, in the unit of the Scene coordinates.
const (
	NodeW        = 184
	NodeH        = 52
	ColGap       = 56
	RowGap       = 16
	ContainerPad = 16
	HeaderH      = 24
)

// Layout fills the integer coordinates of every container, node and edge, and the
// bounding box of the Scene. Provider containers are stacked alphabetically at
// the origin. Layout is deterministic.
func Layout(s *Scene) {
	y := 0
	width := 0
	for i, c := range s.Containers {
		if i > 0 {
			y += RowGap
		}
		sizeContainer(c, s.Edges)
		placeContainer(c, 0, y)
		y += c.H
		if c.W > width {
			width = c.W
		}
	}
	s.Width, s.Height = width, y
	routeEdges(s)
}

// sizeContainer computes W and H from the content and stores node offsets relative
// to the container origin: nodes first, then child containers stacked beneath.
func sizeContainer(c *Container, edges []Edge) {
	nodesW, nodesH := arrangeNodes(c, edges)
	contentW, contentH := nodesW, nodesH
	for _, ch := range c.Children {
		sizeContainer(ch, edges)
		if contentH > 0 {
			contentH += RowGap
		}
		contentH += ch.H
		if ch.W > contentW {
			contentW = ch.W
		}
	}
	c.W = contentW + 2*ContainerPad
	c.H = HeaderH + contentH + 2*ContainerPad
	if len(c.Nodes) == 0 && len(c.Children) == 0 {
		c.H = HeaderH + 2*ContainerPad
	}
}

// arrangeNodes positions the direct nodes relative to the container origin and
// returns the size of the block they occupy. A container with an edge between two
// of its direct nodes uses layers; otherwise a regular grid.
func arrangeNodes(c *Container, edges []Edge) (int, int) {
	if len(c.Nodes) == 0 {
		return 0, 0
	}
	if hasIntraEdge(c, edges) {
		cols := layerNodes(c.Nodes, edges)
		rows := 0
		for ci, col := range cols {
			for ri, n := range col {
				setBox(n, ContainerPad+ci*(NodeW+ColGap), HeaderH+ContainerPad+ri*(NodeH+RowGap))
			}
			if len(col) > rows {
				rows = len(col)
			}
		}
		return len(cols)*NodeW + (len(cols)-1)*ColGap, rows*NodeH + (rows-1)*RowGap
	}
	n := len(c.Nodes)
	cols := ceilSqrt(n)
	rows := (n + cols - 1) / cols
	for i, node := range c.Nodes {
		setBox(node, ContainerPad+(i%cols)*(NodeW+ColGap), HeaderH+ContainerPad+(i/cols)*(NodeH+RowGap))
	}
	return cols*NodeW + (cols-1)*ColGap, rows*NodeH + (rows-1)*RowGap
}

func setBox(n *Node, x, y int) {
	n.X, n.Y, n.W, n.H = x, y, NodeW, NodeH
}

func hasIntraEdge(c *Container, edges []Edge) bool {
	in := make(map[*Node]bool, len(c.Nodes))
	for _, n := range c.Nodes {
		in[n] = true
	}
	for _, e := range edges {
		if e.To != nil && e.From != e.To && in[e.From] && in[e.To] {
			return true
		}
	}
	return false
}

// ceilSqrt returns ceil(sqrt(n)) using integers only.
func ceilSqrt(n int) int {
	r := 0
	for r*r < n {
		r++
	}
	return r
}

// placeContainer turns relative offsets into absolute coordinates.
func placeContainer(c *Container, x, y int) {
	c.X, c.Y = x, y
	for _, n := range c.Nodes {
		n.X += x
		n.Y += y
	}
	cy := y + HeaderH + ContainerPad
	if len(c.Nodes) > 0 {
		cy += nodesBlockHeight(c) + RowGap
	}
	for _, ch := range c.Children {
		placeContainer(ch, x+ContainerPad, cy)
		cy += ch.H + RowGap
	}
}

// nodesBlockHeight returns the height of the block of direct nodes, from the
// offsets stored by arrangeNodes (still relative when it is called).
func nodesBlockHeight(c *Container) int {
	bottom := 0
	for _, n := range c.Nodes {
		if b := n.Y - c.Y + n.H; b > bottom {
			bottom = b
		}
	}
	return bottom - HeaderH - ContainerPad
}
