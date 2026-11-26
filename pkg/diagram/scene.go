package diagram

// Status is the change status of a node. The declaration order is the fixed
// order used for sorting and for the legend.
type Status int

// Statuses in their fixed order.
const (
	StatusCreate Status = iota
	StatusUpdate
	StatusDelete
	StatusReplace
	StatusNoop
)

// String returns the lower-case name of the status.
func (s Status) String() string {
	switch s {
	case StatusCreate:
		return "create"
	case StatusUpdate:
		return "update"
	case StatusDelete:
		return "delete"
	case StatusReplace:
		return "replace"
	default:
		return "noop"
	}
}

// ContainerKind tells a provider container from a module container.
type ContainerKind string

// Container kinds.
const (
	KindProvider ContainerKind = "provider"
	KindModule   ContainerKind = "module"
)

// Options tells Build which filters the user gave; it decides no-op visibility
// and dependency context.
type Options struct {
	TypeFiltered   bool
	ActionFiltered bool
}

// Point is an integer coordinate.
type Point struct{ X, Y int }

// Scene is the picture of one plan: everything a renderer needs, in integer
// coordinates. No field holds a resource attribute value.
type Scene struct {
	Containers []*Container // one per provider, sorted by provider name
	Nodes      []*Node      // sorted by address
	Edges      []Edge       // sorted by (from, to)
	Statuses   []Status     // statuses present, in fixed order
	Width      int          // content bounding box
	Height     int
	Empty      bool // nothing is drawn
}

// Container groups nodes: a provider, or a module instance inside a provider.
type Container struct {
	ID       string
	Kind     ContainerKind
	Label    string // provider name, or module step with key such as module.app[0]
	Path     string // full module instance path, used for edge targets
	Children []*Container
	Nodes    []*Node
	X, Y     int
	W, H     int

	configPath string // index-free module path, matches module references
	provider   string // provider the container sits under
}

// Node is one drawn resource instance.
//
// Address is unique across nodes; a node exists only for create, update, delete,
// replace or (requested or context) no-op; never for data sources.
type Node struct {
	ID            string
	Address       string // verbatim, with index
	ConfigAddress string // index-free, used to expand edges
	Type          string
	Name          string // resource name with its index
	Provider      string // short provider name
	Status        Status
	Context       bool     // unchanged, drawn only as dependency context
	Family        []string // ordered icon chain, never empty
	Layer         int
	X, Y          int
	W, H          int
}

// Edge links a dependent node to its dependency: a node, or a module container.
//
// No edge starts or ends at a delete node, no edge starts at a context node, edges
// towards nodes not drawn are dropped and duplicates are removed.
type Edge struct {
	From        *Node
	To          *Node      // nil when the target is a container
	ToContainer *Container // nil when the target is a node
	Points      []Point    // cubic Bezier: start, two controls, end; filled by layout
}
