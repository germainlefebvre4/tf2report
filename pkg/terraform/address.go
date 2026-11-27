package terraform

import "strings"

// modeData is the mode of a data source address.
const (
	modeData      = "data"
	moduleKeyword = "module"
)

// ModuleStep is one `module.<name>[<key>]` step of a resource address.
type ModuleStep struct {
	Name string
	Key  string
}

// Address is a parsed resource address such as
// `module.app[0].aws_instance.web["a"]`.
type Address struct {
	Modules []ModuleStep
	Mode    string // "managed" or "data"
	Type    string
	Name    string
	Key     string
}

// ParseAddress tokenizes a resource address. It never fails: a malformed
// address yields a best-effort value.
func ParseAddress(s string) Address {
	toks := splitAddress(s)
	var addr Address
	i := 0
	for i+1 < len(toks) && toks[i].name == moduleKeyword {
		addr.Modules = append(addr.Modules, ModuleStep{Name: toks[i+1].name, Key: toks[i+1].key})
		i += 2
	}
	addr.Mode = "managed"
	if i < len(toks) && toks[i].name == modeData {
		addr.Mode = modeData
		i++
	}
	if i < len(toks) {
		addr.Type = toks[i].name
		i++
	}
	if i < len(toks) {
		addr.Name = toks[i].name
		addr.Key = toks[i].key
	}
	return addr
}

// ConfigAddress returns the full address with every `[...]` index removed.
func (a Address) ConfigAddress() string {
	var parts []string
	for _, m := range a.Modules {
		parts = append(parts, moduleKeyword, m.Name)
	}
	if a.Mode == modeData {
		parts = append(parts, modeData)
	}
	if a.Type != "" {
		parts = append(parts, a.Type)
	}
	if a.Name != "" {
		parts = append(parts, a.Name)
	}
	return strings.Join(parts, ".")
}

type addrToken struct {
	name string
	key  string // raw key text including brackets, e.g. `["a.b"]` or `[0]`
}

// splitAddress splits on dots outside brackets and quotes and attaches each
// bracketed key to the preceding segment.
func splitAddress(s string) []addrToken {
	var toks []addrToken
	var cur strings.Builder
	var key strings.Builder
	depth := 0
	inQuote := false
	escaped := false
	flush := func() {
		if cur.Len() > 0 || key.Len() > 0 {
			toks = append(toks, addrToken{name: cur.String(), key: key.String()})
		}
		cur.Reset()
		key.Reset()
	}
	for _, r := range s {
		switch {
		case depth > 0:
			key.WriteRune(r)
			switch {
			case escaped:
				escaped = false
			case inQuote && r == '\\':
				escaped = true
			case r == '"':
				inQuote = !inQuote
			case !inQuote && r == '[':
				depth++
			case !inQuote && r == ']':
				depth--
			}
		case r == '[':
			depth = 1
			key.WriteRune(r)
		case r == '.':
			flush()
		default:
			cur.WriteRune(r)
		}
	}
	flush()
	return toks
}
