package terraform

// ReferenceKind classifies what a reference string points to.
type ReferenceKind string

// Reference kinds.
const (
	RefResource ReferenceKind = "resource"
	RefModule   ReferenceKind = "module"
	RefData     ReferenceKind = "data"
	RefOther    ReferenceKind = "other"
)

// Reference is a parsed reference string of the plan configuration.
type Reference struct {
	Kind       ReferenceKind
	ModuleCall string // set for RefModule
	Type       string // set for RefResource
	Name       string // set for RefResource
}

var nonResourcePrefixes = map[string]bool{
	"var": true, "local": true, "count": true, "each": true,
	"path": true, "terraform": true, "self": true,
}

// ParseReference classifies a reference such as `aws_subnet.a.id`,
// `module.net.vpc_id` or `var.x`. It never fails: unknown forms are RefOther.
func ParseReference(s string) Reference {
	toks := splitAddress(s)
	if len(toks) == 0 || toks[0].name == "" {
		return Reference{Kind: RefOther}
	}
	first := toks[0].name
	switch {
	case first == moduleKeyword:
		if len(toks) >= 2 && toks[1].name != "" {
			return Reference{Kind: RefModule, ModuleCall: toks[1].name}
		}
	case first == "data":
		return Reference{Kind: RefData}
	case nonResourcePrefixes[first]:
	case len(toks) >= 2 && toks[1].name != "":
		return Reference{Kind: RefResource, Type: first, Name: toks[1].name}
	}
	return Reference{Kind: RefOther}
}
