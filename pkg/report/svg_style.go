package report

import "github.com/germainlefebvre4/tf2report/pkg/diagram"

// statusStyle holds every visual cue of one status: color, symbol and border.
// Color is never the only cue.
type statusStyle struct {
	label   string
	fill    string // tint behind the label
	stroke  string // border and symbol color
	symbol  string // empty for no-op
	dash    string // stroke-dasharray, empty for solid
	strokeW int
	opacity string // empty for fully opaque
}

const (
	colorText      = "#1f2328"
	colorEdge      = "#57606a"
	colorBackdrop  = "#ffffff"
	colorProvider  = "#f6f8fa"
	colorModule    = "#fbfcfd"
	colorContainer = "#afb8c1"
)

// statusStyles is the single status table, indexed by diagram.Status. The legend
// reads the same table. Text on every tint keeps a contrast ratio above 4.5:1.
var statusStyles = [...]statusStyle{
	diagram.StatusCreate:  {"create", "#dafbe1", "#116329", "+", "", 2, ""},
	diagram.StatusUpdate:  {"update", "#ffe8cc", "#953800", "~", "6 3", 2, ""},
	diagram.StatusDelete:  {"delete", "#ffdcd7", "#a40e26", "-", "2 3", 2, ""},
	diagram.StatusReplace: {"replace", "#eddcff", "#6639ba", "-/+", "8 3 2 3", 2, ""},
	diagram.StatusNoop:    {"no-op", "#eaeef2", "#6e7781", "", "", 1, "0.6"},
}

// Layout metrics of the document around the Scene.
const (
	docMargin     = 24
	legendHeight  = 28
	legendEntryW  = 150
	legendGap     = 20
	minDocWidth   = 320
	charWidth     = 7 // estimated width of one character at 11 px, monospace
	iconSize      = 24
	nodeTextLeft  = 40
	nodeTextRight = 8
)
