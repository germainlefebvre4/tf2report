package diagram

import (
	"sort"

	"github.com/germainlefebvre4/tf2report/pkg/terraform"
)

// relation is a dependency between configuration objects, index-free. Exactly one
// of toResource and toModule is set.
type relation struct {
	from       string // config address of the dependent resource
	toResource string // config address of a resource
	toModule   string // index-free path of a module call
}

// resolveRelations walks the plan configuration module by module. A relative
// reference resolves inside the module that holds the resource; `module.X`
// resolves to the child module call. Data sources and non-resource references are
// skipped. A plan without configuration yields no relation.
func resolveRelations(cfg *terraform.Configuration) []relation {
	if cfg == nil {
		return nil
	}
	var rels []relation
	walkModule(cfg.RootModule, "", &rels)
	sort.Slice(rels, func(i, j int) bool {
		a, b := rels[i], rels[j]
		if a.from != b.from {
			return a.from < b.from
		}
		if a.toResource != b.toResource {
			return a.toResource < b.toResource
		}
		return a.toModule < b.toModule
	})
	out := rels[:0]
	for i, r := range rels {
		if i == 0 || r != rels[i-1] {
			out = append(out, r)
		}
	}
	return out
}

func walkModule(mod terraform.ConfigModule, path string, rels *[]relation) {
	for _, res := range mod.Resources {
		if res.Mode == modeData {
			continue
		}
		from := joinPath(path, res.Type+"."+res.Name)
		for _, ref := range res.References() {
			r := terraform.ParseReference(ref)
			switch r.Kind {
			case terraform.RefResource:
				*rels = append(*rels, relation{from: from, toResource: joinPath(path, r.Type+"."+r.Name)})
			case terraform.RefModule:
				*rels = append(*rels, relation{from: from, toModule: joinPath(path, "module."+r.ModuleCall)})
			}
		}
	}
	names := make([]string, 0, len(mod.ModuleCalls))
	for name := range mod.ModuleCalls {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		walkModule(mod.ModuleCalls[name].Module, joinPath(path, "module."+name), rels)
	}
}

// resourceEdges expands resource relations to every instance on both sides. When
// context is allowed, an unchanged target that is not drawn becomes a context
// node. Edges from or to a delete node, from a node not drawn, and self edges are
// dropped.
func resourceEdges(
	rels []relation, byConfig map[string][]*Node, drawn map[string]bool, allowContext bool,
) (edges []Edge, context map[string]bool) {
	context = make(map[string]bool)
	for _, r := range rels {
		if r.toResource == "" {
			continue
		}
		for _, f := range byConfig[r.from] {
			if !drawn[f.Address] || f.Status == StatusDelete {
				continue
			}
			for _, t := range byConfig[r.toResource] {
				if t == f || t.Status == StatusDelete {
					continue
				}
				if !drawn[t.Address] {
					if !allowContext || t.Status != StatusNoop {
						continue
					}
					context[t.Address] = true
				}
				edges = append(edges, Edge{From: f, To: t})
			}
		}
	}
	return edges, context
}

// moduleEdges links dependents to every container whose index-free path matches a
// referenced module call. A container that holds the dependent is skipped.
func moduleEdges(rels []relation, byConfig map[string][]*Node, drawn map[string]bool, roots []*Container) []Edge {
	byPath := make(map[string][]*Container)
	var collect func(c *Container)
	collect = func(c *Container) {
		if c.Kind == KindModule {
			byPath[c.configPath] = append(byPath[c.configPath], c)
		}
		for _, ch := range c.Children {
			collect(ch)
		}
	}
	for _, c := range roots {
		collect(c)
	}
	var edges []Edge
	for _, r := range rels {
		if r.toModule == "" {
			continue
		}
		for _, f := range byConfig[r.from] {
			if !drawn[f.Address] || f.Status == StatusDelete {
				continue
			}
			for _, c := range byPath[r.toModule] {
				if !holds(c, f) {
					edges = append(edges, Edge{From: f, ToContainer: c})
				}
			}
		}
	}
	return edges
}

// holds reports whether the node sits in the container or one of its descendants.
func holds(c *Container, n *Node) bool {
	for _, x := range c.Nodes {
		if x == n {
			return true
		}
	}
	for _, ch := range c.Children {
		if holds(ch, n) {
			return true
		}
	}
	return false
}

func edgeKey(e Edge) (from, to string) {
	if e.To != nil {
		return e.From.Address, "0" + e.To.Address
	}
	return e.From.Address, "1" + e.ToContainer.ID
}

// finalizeEdges sorts edges by (from, to) and removes duplicates.
func finalizeEdges(edges []Edge) []Edge {
	sort.Slice(edges, func(i, j int) bool {
		af, at := edgeKey(edges[i])
		bf, bt := edgeKey(edges[j])
		if af != bf {
			return af < bf
		}
		return at < bt
	})
	out := edges[:0]
	for i, e := range edges {
		if i > 0 {
			pf, pt := edgeKey(edges[i-1])
			cf, ct := edgeKey(e)
			if pf == cf && pt == ct {
				continue
			}
		}
		out = append(out, e)
	}
	return out
}
