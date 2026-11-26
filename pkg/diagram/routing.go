package diagram

const minCurve = 28

// routeEdges fills Points of every edge with the four points of an integer cubic
// Bezier curve: it leaves the left side of the dependent, with a horizontal
// tangent, and ends at the right side of the dependency node (or of the target
// container).
func routeEdges(s *Scene) {
	for i := range s.Edges {
		e := &s.Edges[i]
		sx, sy := e.From.X, e.From.Y+e.From.H/2
		var ex, ey int
		if e.To != nil {
			ex, ey = e.To.X+e.To.W, e.To.Y+e.To.H/2
		} else {
			ex, ey = e.ToContainer.X+e.ToContainer.W, e.ToContainer.Y+e.ToContainer.H/2
		}
		// Arrows normally point left. When the dependency lies to the right (across
		// containers) the curve keeps a short overshoot instead of looping far away.
		d := minCurve
		if ex <= sx && abs(sx-ex)/2 > d {
			d = abs(sx-ex) / 2
		}
		e.Points = []Point{{sx, sy}, {sx - d, sy}, {ex + d, ey}, {ex, ey}}
	}
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
