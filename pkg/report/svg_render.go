package report

import (
	"encoding/xml"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/germainlefebvre4/tf2report/pkg/diagram"
)

const fontFamily = "ui-monospace, Menlo, Consolas, monospace"

// renderer accumulates the document in memory so that it is written once.
type renderer struct {
	b     strings.Builder
	scene *diagram.Scene
	icons *iconSet
	ox    int // offset of the scene inside the document
	oy    int
	used  map[string]icon // icons used by nodes, by family id
	chain map[*diagram.Node]string
}

// renderScene turns a laid-out Scene into an SVG document.
func renderScene(s *diagram.Scene, icons *iconSet) string {
	r := &renderer{scene: s, icons: icons, used: map[string]icon{}, chain: map[*diagram.Node]string{}}
	for _, n := range s.Nodes {
		if id, ic, ok := icons.resolve(n.Family); ok {
			r.used[id] = ic
			r.chain[n] = id
		}
	}
	width, height := r.documentSize()
	r.ox = docMargin + (width-2*docMargin-s.Width)/2
	r.oy = docMargin + legendHeight + docMargin
	r.writeHeader(width, height)
	r.writeDefs()
	r.writeLegend()
	if s.Empty {
		r.writeEmpty(width, height)
	} else {
		r.writeContainers()
		r.writeEdges()
		r.writeNodes()
	}
	r.b.WriteString("</svg>\n")
	return r.b.String()
}

func (r *renderer) documentSize() (int, int) {
	legendW := len(r.scene.Statuses)*legendEntryW + 2*docMargin
	if len(r.scene.Statuses) > 1 {
		legendW += (len(r.scene.Statuses) - 1) * legendGap
	}
	width := r.scene.Width + 2*docMargin
	for _, v := range []int{legendW, minDocWidth} {
		if v > width {
			width = v
		}
	}
	height := docMargin + legendHeight + docMargin + r.scene.Height + docMargin
	if r.scene.Empty {
		height += 3 * docMargin
	}
	return width, height
}

func (r *renderer) writeHeader(width, height int) {
	r.b.WriteString(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 ` + itoa(width) + ` ` + itoa(height) +
		`" width="` + itoa(width) + `" height="` + itoa(height) + `" role="img" font-family="` + esc(fontFamily) + `" font-size="11">` + "\n")
	r.b.WriteString("<title>Terraform plan diagram</title>\n")
	r.b.WriteString("<desc>" + esc(r.describe()) + "</desc>\n")
	r.b.WriteString(`<rect width="` + itoa(width) + `" height="` + itoa(height) + `" fill="` + colorBackdrop + `"/>` + "\n")
}

// describe summarizes the plan with counts per status.
func (r *renderer) describe() string {
	if r.scene.Empty {
		return "No resources to display"
	}
	counts := map[diagram.Status]int{}
	for _, n := range r.scene.Nodes {
		counts[n.Status]++
	}
	parts := make([]string, 0, len(r.scene.Statuses))
	for _, st := range r.scene.Statuses {
		parts = append(parts, statusStyles[st].label+" "+itoa(counts[st]))
	}
	return itoa(len(r.scene.Nodes)) + " resources: " + strings.Join(parts, ", ")
}

func (r *renderer) writeDefs() {
	r.b.WriteString("<defs>\n")
	r.b.WriteString(`<marker id="arrow" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="8" markerHeight="8" orient="auto">` +
		`<path d="M0 0L10 5L0 10z" fill="` + colorEdge + `"/></marker>` + "\n")
	ids := make(map[string]bool, len(r.used))
	for id := range r.used {
		ids[id] = true
	}
	for _, id := range sortedKeys(ids) {
		ic := r.used[id]
		r.b.WriteString(`<symbol id="icon-` + esc(id) + `" viewBox="` + esc(ic.viewBox) + `">` + ic.inner + "</symbol>\n")
	}
	r.b.WriteString("</defs>\n")
}

func (r *renderer) writeLegend() {
	r.b.WriteString(`<g class="legend">` + "\n")
	x := docMargin
	for _, st := range r.scene.Statuses {
		style := statusStyles[st]
		r.b.WriteString(`<g transform="translate(` + itoa(x) + ` ` + itoa(docMargin) + `)">`)
		r.b.WriteString(boxRect(0, 0, 40, 20, style))
		if style.symbol != "" {
			r.b.WriteString(`<text x="20" y="14" text-anchor="middle" font-weight="bold" fill="` + style.stroke + `">` + esc(style.symbol) + `</text>`)
		}
		r.b.WriteString(`<text x="48" y="14" fill="` + colorText + `">` + esc(style.label) + `</text></g>` + "\n")
		x += legendEntryW + legendGap
	}
	r.b.WriteString("</g>\n")
}

func (r *renderer) writeEmpty(width, height int) {
	r.b.WriteString(`<text x="` + itoa(width/2) + `" y="` + itoa(height-2*docMargin) + `" text-anchor="middle" fill="` + colorText +
		`" font-size="14">No resources to display</text>` + "\n")
}

func (r *renderer) writeContainers() {
	for _, c := range r.scene.Containers {
		r.writeContainer(c)
	}
}

func (r *renderer) writeContainer(c *diagram.Container) {
	fill := colorProvider
	if c.Kind == diagram.KindModule {
		fill = colorModule
	}
	label := truncate(c.Label, (c.W-16)/charWidth)
	r.b.WriteString(`<g class="container ` + string(c.Kind) + `">` +
		`<rect x="` + itoa(r.ox+c.X) + `" y="` + itoa(r.oy+c.Y) + `" width="` + itoa(c.W) + `" height="` + itoa(c.H) +
		`" rx="6" fill="` + fill + `" stroke="` + colorContainer + `" stroke-width="1"/>` +
		`<text x="` + itoa(r.ox+c.X+8) + `" y="` + itoa(r.oy+c.Y+16) + `" font-weight="bold" fill="` + colorText + `">` + esc(label) + `</text></g>` + "\n")
	for _, ch := range c.Children {
		r.writeContainer(ch)
	}
}

func (r *renderer) writeEdges() {
	r.b.WriteString(`<g class="edges">` + "\n")
	for _, e := range r.scene.Edges {
		if len(e.Points) != 4 {
			continue
		}
		p := e.Points
		d := "M" + pt(r, p[0]) + " C" + pt(r, p[1]) + " " + pt(r, p[2]) + " " + pt(r, p[3])
		r.b.WriteString(`<path d="` + d + `" fill="none" stroke="` + colorEdge + `" stroke-width="1.5" marker-end="url(#arrow)"/>` + "\n")
	}
	r.b.WriteString("</g>\n")
}

func pt(r *renderer, p diagram.Point) string {
	return itoa(r.ox+p.X) + " " + itoa(r.oy+p.Y)
}

func (r *renderer) writeNodes() {
	r.b.WriteString(`<g class="nodes">` + "\n")
	for _, n := range r.scene.Nodes {
		r.writeNode(n)
	}
	r.b.WriteString("</g>\n")
}

func (r *renderer) writeNode(n *diagram.Node) {
	style := statusStyles[n.Status]
	class := "node status-" + n.Status.String()
	symbol := style.symbol
	if n.Context {
		class += " context"
		symbol = ""
	}
	x, y := r.ox+n.X, r.oy+n.Y
	r.b.WriteString(`<g class="` + class + `" color="` + colorText + `">`)
	r.b.WriteString("<title>" + esc(n.Address) + "</title>")
	r.b.WriteString(boxRect(x, y, n.W, n.H, style))
	if id, ok := r.chain[n]; ok {
		r.b.WriteString(`<use href="#icon-` + esc(id) + `" x="` + itoa(x+8) + `" y="` + itoa(y+(n.H-iconSize)/2) +
			`" width="` + itoa(iconSize) + `" height="` + itoa(iconSize) + `"/>`)
	}
	nameMax := (n.W - nodeTextLeft - nodeTextRight - charWidth*utf8.RuneCountInString(symbol) - 4) / charWidth
	r.b.WriteString(`<text x="` + itoa(x+nodeTextLeft) + `" y="` + itoa(y+21) + `" font-weight="bold" fill="` + colorText + `">` +
		esc(truncate(n.Name, nameMax)) + `</text>`)
	if symbol != "" {
		r.b.WriteString(`<text x="` + itoa(x+n.W-nodeTextRight) + `" y="` + itoa(y+21) + `" text-anchor="end" font-weight="bold" fill="` +
			style.stroke + `">` + esc(symbol) + `</text>`)
	}
	typeMax := (n.W - nodeTextLeft - nodeTextRight) / charWidth
	r.b.WriteString(`<text x="` + itoa(x+nodeTextLeft) + `" y="` + itoa(y+38) + `" fill="` + colorText + `">` +
		esc(truncate(n.Type, typeMax)) + `</text></g>` + "\n")
}

// boxRect draws a rounded box with the fill, border style and opacity of a status.
func boxRect(x, y, w, h int, s statusStyle) string {
	out := `<rect x="` + itoa(x) + `" y="` + itoa(y) + `" width="` + itoa(w) + `" height="` + itoa(h) +
		`" rx="5" fill="` + s.fill + `" stroke="` + s.stroke + `" stroke-width="` + itoa(s.strokeW) + `"`
	if s.dash != "" {
		out += ` stroke-dasharray="` + s.dash + `"`
	}
	if s.opacity != "" {
		out += ` opacity="` + s.opacity + `"`
	}
	return out + `/>`
}

// truncate shortens s to at most max runes, ending with an ellipsis when cut.
func truncate(s string, max int) string {
	if max < 1 {
		max = 1
	}
	if utf8.RuneCountInString(s) <= max {
		return s
	}
	runes := []rune(s)
	return string(runes[:max-1]) + "…"
}

func itoa(v int) string { return strconv.Itoa(v) }

// esc escapes a text or attribute value, after removing characters that are
// illegal in XML 1.0.
func esc(s string) string {
	var clean strings.Builder
	for _, r := range s {
		if r == utf8.RuneError || !isXMLChar(r) {
			continue
		}
		clean.WriteRune(r)
	}
	var out strings.Builder
	_ = xml.EscapeText(&out, []byte(clean.String()))
	return out.String()
}

func isXMLChar(r rune) bool {
	return r == 0x9 || r == 0xA || r == 0xD ||
		(r >= 0x20 && r <= 0xD7FF) || (r >= 0xE000 && r <= 0xFFFD) || (r >= 0x10000 && r <= 0x10FFFF)
}
