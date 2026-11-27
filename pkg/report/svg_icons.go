package report

import (
	"bytes"
	"embed"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// builtinFS holds the original generic glyphs; vendorFS holds the cloud provider
// icons (selection of official artwork). They live apart because the vendor icons
// are regenerated over the years by scripts/build-icons.py.
//
//go:embed icons/*.svg
var builtinFS embed.FS

//go:embed icons/vendor/*/*.svg
var vendorFS embed.FS

// maxIconSize is the largest icon file accepted from a pack.
const maxIconSize = 256 * 1024

const svgNS = "http://www.w3.org/2000/svg"

// allowedElements is the allowlist of the icon sanitizer.
var allowedElements = map[string]bool{
	"svg": true, "g": true, "path": true, "rect": true, "circle": true, "ellipse": true,
	"line": true, "polyline": true, "polygon": true, "defs": true, "linearGradient": true,
	"radialGradient": true, "stop": true, "clipPath": true, "title": true, "desc": true,
}

// allowedAttrs is the presentation allowlist; every other attribute rejects the file.
var allowedAttrs = map[string]bool{
	"id": true, "viewBox": true, "width": true, "height": true, "x": true, "y": true, "cx": true,
	"cy": true, "r": true, "rx": true, "ry": true, "x1": true, "y1": true, "x2": true, "y2": true,
	"points": true, "d": true, "fill": true, "fill-opacity": true, "fill-rule": true, "clip-rule": true,
	"stroke": true, "stroke-width": true, "stroke-linecap": true, "stroke-linejoin": true,
	"stroke-dasharray": true, "stroke-dashoffset": true, "stroke-miterlimit": true, "stroke-opacity": true,
	"opacity": true, "transform": true, "style": true, "offset": true, "stop-color": true,
	"stop-opacity": true, "gradientUnits": true, "gradientTransform": true, "spreadMethod": true,
	"fx": true, "fy": true, "clip-path": true, "clipPathUnits": true, "href": true,
	"preserveAspectRatio": true, "version": true,
}

var urlRef = regexp.MustCompile(`url\(\s*#([^)\s]+)\s*\)`)

// icon is a sanitized icon, ready to be emitted as a <symbol>.
type icon struct {
	viewBox string
	inner   string
}

// iconPack holds sanitized icons by family id.
type iconPack struct {
	icons map[string]icon
}

// iconSet resolves a family chain at each level against the user pack first, then
// the embedded vendor icons, then the built-in glyphs.
type iconSet struct {
	pack    *iconPack
	vendor  *iconPack
	builtin *iconPack
}

// newIconSet loads the built-in glyphs and, when dir is not empty, the user pack.
// An unreadable or missing directory is an error that names it.
func newIconSet(dir string, onSkip func(file, reason string)) (*iconSet, error) {
	builtin, err := loadBuiltin()
	if err != nil {
		return nil, err
	}
	vendor, err := loadVendor()
	if err != nil {
		return nil, err
	}
	set := &iconSet{builtin: builtin, vendor: vendor, pack: &iconPack{icons: map[string]icon{}}}
	if dir != "" {
		pack, err := loadPack(dir, onSkip)
		if err != nil {
			return nil, err
		}
		set.pack = pack
	}
	return set, nil
}

func loadBuiltin() (*iconPack, error) {
	entries, err := builtinFS.ReadDir("icons")
	if err != nil {
		return nil, fmt.Errorf("read built-in icons: %w", err)
	}
	pack := &iconPack{icons: make(map[string]icon, len(entries))}
	for _, e := range entries {
		data, err := builtinFS.ReadFile("icons/" + e.Name())
		if err != nil {
			return nil, fmt.Errorf("read built-in icon %s: %w", e.Name(), err)
		}
		id := strings.TrimSuffix(e.Name(), ".svg")
		ic, err := sanitizeIcon(data, id)
		if err != nil {
			return nil, fmt.Errorf("built-in icon %s: %w", e.Name(), err)
		}
		pack.icons[id] = ic
	}
	return pack, nil
}

// loadVendor sanitizes the embedded cloud provider icons, keyed by family id.
func loadVendor() (*iconPack, error) {
	pack := &iconPack{icons: map[string]icon{}}
	err := fs.WalkDir(vendorFS, "icons/vendor", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".svg") {
			return err
		}
		data, err := vendorFS.ReadFile(path)
		if err != nil {
			return err
		}
		id := strings.TrimSuffix(d.Name(), ".svg")
		ic, err := sanitizeIcon(data, id)
		if err != nil {
			return fmt.Errorf("vendor icon %s: %w", path, err)
		}
		pack.icons[id] = ic
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("read vendor icons: %w", err)
	}
	return pack, nil
}

// loadPack reads the *.svg files of dir, non-recursively. A file that fails the
// sanitizer is reported through onSkip and ignored; a directory with no usable
// file is not an error.
func loadPack(dir string, onSkip func(file, reason string)) (*iconPack, error) {
	info, err := os.Stat(dir)
	if err != nil {
		return nil, fmt.Errorf("icons directory %s: %w", dir, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("icons directory %s: not a directory", dir)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("icons directory %s: %w", dir, err)
	}
	pack := &iconPack{icons: make(map[string]icon)}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".svg") {
			continue
		}
		skip := func(reason string) {
			if onSkip != nil {
				onSkip(e.Name(), reason)
			}
		}
		data, err := readLimited(filepath.Join(dir, e.Name()))
		if err != nil {
			skip(err.Error())
			continue
		}
		id := strings.TrimSuffix(e.Name(), ".svg")
		ic, err := sanitizeIcon(data, id)
		if err != nil {
			skip(err.Error())
			continue
		}
		pack.icons[id] = ic
	}
	return pack, nil
}

// readLimited reads a file and rejects it when it exceeds maxIconSize.
func readLimited(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close() //nolint:errcheck // read-only file, a close failure has no consequence
	data, err := io.ReadAll(io.LimitReader(f, maxIconSize+1))
	if err != nil {
		return nil, err
	}
	if len(data) > maxIconSize {
		return nil, fmt.Errorf("file larger than %d KiB", maxIconSize/1024)
	}
	return data, nil
}

// resolve returns the family id and icon of the first level of the chain that the
// pack or the built-in set can serve.
func (s *iconSet) resolve(chain []string) (string, icon, bool) {
	for _, id := range chain {
		if ic, ok := s.pack.icons[id]; ok {
			return id, ic, true
		}
		if ic, ok := s.vendor.icons[id]; ok {
			return id, ic, true
		}
		if ic, ok := s.builtin.icons[id]; ok {
			return id, ic, true
		}
	}
	return "", icon{}, false
}

// sanitizeIcon validates an SVG file against the allowlist and re-serializes it
// with ids prefixed by the family id. Nothing of the original bytes is copied.
func sanitizeIcon(data []byte, family string) (icon, error) {
	sz := &sanitizer{dec: xml.NewDecoder(bytes.NewReader(data)), family: family}
	sz.dec.Strict = true
	for {
		tok, err := sz.dec.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return icon{}, fmt.Errorf("malformed XML: %w", err)
		}
		if err := sz.token(tok); err != nil {
			return icon{}, err
		}
	}
	if !sz.seenRoot {
		return icon{}, errors.New("no root element")
	}
	if sz.viewBox == "" {
		sz.viewBox = "0 0 24 24"
	}
	return icon{viewBox: sz.viewBox, inner: sz.out.String()}, nil
}

// sanitizer holds the state of one sanitizeIcon run.
type sanitizer struct {
	dec      *xml.Decoder
	family   string
	out      strings.Builder
	stack    []string
	viewBox  string
	seenRoot bool
}

func (z *sanitizer) token(tok xml.Token) error {
	switch t := tok.(type) {
	case xml.Directive:
		return errors.New("DOCTYPE or entity declaration not allowed")
	case xml.ProcInst:
		if t.Target != "xml" {
			return errors.New("processing instruction not allowed")
		}
	case xml.StartElement:
		return z.start(t)
	case xml.EndElement:
		z.stack = z.stack[:len(z.stack)-1]
		if len(z.stack) > 0 {
			z.out.WriteString("</" + t.Name.Local + ">")
		}
	case xml.CharData:
		if n := len(z.stack); n > 0 && (z.stack[n-1] == "title" || z.stack[n-1] == "desc") {
			_ = xml.EscapeText(&z.out, t)
		}
	}
	return nil
}

func (z *sanitizer) start(t xml.StartElement) error {
	root := len(z.stack) == 0
	if root {
		if z.seenRoot {
			return errors.New("more than one root element")
		}
		z.seenRoot = true
		if t.Name.Local != "svg" {
			return fmt.Errorf("root element is %q, not svg", t.Name.Local)
		}
	}
	vb, err := writeStart(&z.out, t, z.family, root)
	if err != nil {
		return err
	}
	if root {
		z.viewBox = vb
	}
	z.stack = append(z.stack, t.Name.Local)
	return nil
}

// writeStart validates one start tag and writes its sanitized form. The root svg
// element is not written (it becomes the <symbol>); its viewBox is returned.
func writeStart(out *strings.Builder, t xml.StartElement, family string, root bool) (string, error) {
	if t.Name.Space != "" && t.Name.Space != svgNS {
		return "", fmt.Errorf("element in foreign namespace %q", t.Name.Space)
	}
	if !allowedElements[t.Name.Local] {
		return "", fmt.Errorf("element <%s> not allowed", t.Name.Local)
	}
	viewBox := ""
	var attrs strings.Builder
	for _, a := range t.Attr {
		name, value, keep, err := checkAttr(a, family)
		if err != nil {
			return "", err
		}
		if !keep {
			continue
		}
		if root {
			if name == "viewBox" {
				viewBox = value
			}
			continue
		}
		attrs.WriteString(" " + name + `="`)
		_ = xml.EscapeText(&attrs, []byte(value))
		attrs.WriteString(`"`)
	}
	if !root {
		out.WriteString("<" + t.Name.Local + attrs.String() + ">")
	}
	return viewBox, nil
}

// checkAttr validates one attribute. Namespace declarations are dropped; every
// other attribute must be on the allowlist, and references must stay local.
func checkAttr(a xml.Attr, family string) (name, value string, keep bool, err error) {
	if a.Name.Space == "xmlns" || (a.Name.Space == "" && a.Name.Local == "xmlns") {
		return "", "", false, nil
	}
	name = a.Name.Local
	if strings.HasPrefix(strings.ToLower(name), "on") {
		return "", "", false, fmt.Errorf("event attribute %q not allowed", name)
	}
	if a.Name.Space != "" && (name != "href" || a.Name.Space != "http://www.w3.org/1999/xlink") {
		return "", "", false, fmt.Errorf("attribute %s:%s not allowed", a.Name.Space, name)
	}
	if !allowedAttrs[name] {
		return "", "", false, fmt.Errorf("attribute %q not allowed", name)
	}
	value, err = checkValue(name, a.Value, family)
	if err != nil {
		return "", "", false, err
	}
	return name, value, true, nil
}

// checkValue rejects external references and prefixes local ids and references.
func checkValue(name, value, family string) (string, error) {
	lower := strings.ToLower(value)
	switch {
	case name == "href":
		if !strings.HasPrefix(value, "#") {
			return "", errors.New("href that is not a local fragment")
		}
		return "#" + family + "-" + value[1:], nil
	case name == "id":
		return family + "-" + value, nil
	case name == "style" && strings.Contains(lower, "url("):
		return "", errors.New("url() in style attribute")
	case strings.Contains(lower, "@import"), strings.Contains(lower, "javascript:"), strings.Contains(lower, "expression("):
		return "", fmt.Errorf("forbidden content in %q", name)
	}
	if strings.Contains(lower, "url(") {
		rest := urlRef.ReplaceAllString(value, "")
		if strings.Contains(strings.ToLower(rest), "url(") {
			return "", fmt.Errorf("external url() in %q", name)
		}
		return urlRef.ReplaceAllString(value, "url(#"+family+"-$1)"), nil
	}
	return value, nil
}

// sortedKeys returns the keys of a set in lexical order.
func sortedKeys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
