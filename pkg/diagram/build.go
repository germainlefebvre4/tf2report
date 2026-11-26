package diagram

import (
	"sort"
	"strconv"
	"strings"

	"github.com/germainlefebvre4/tf2report/pkg/terraform"
)

// Build selects the resources to draw, groups them in containers, resolves the
// dependency edges and adds the dependency context. Coordinates are filled by
// Layout. The Scene holds no attribute value.
func Build(plan *terraform.Plan, opts Options) (*Scene, error) {
	drawn := Select(plan, opts)
	candidates := candidateNodes(plan)
	drawnSet := make(map[string]bool, len(drawn))
	for _, rc := range drawn {
		drawnSet[rc.Address] = true
	}

	relations := resolveRelations(plan.Configuration)
	byConfig := indexByConfig(candidates)
	edges, context := resourceEdges(relations, byConfig, drawnSet, !opts.TypeFiltered && !opts.ActionFiltered)

	nodes := make([]*Node, 0, len(drawn)+len(context))
	for _, n := range candidates {
		if drawnSet[n.Address] || context[n.Address] {
			n.Context = context[n.Address] && !drawnSet[n.Address]
			nodes = append(nodes, n)
		}
	}
	sort.Slice(nodes, func(i, j int) bool { return nodes[i].Address < nodes[j].Address })
	for i, n := range nodes {
		n.ID = "n" + strconv.Itoa(i+1)
	}

	scene := &Scene{Nodes: nodes, Empty: len(nodes) == 0}
	scene.Containers = buildContainers(nodes)
	scene.Edges = append(edges, moduleEdges(relations, byConfig, drawnSet, scene.Containers)...)
	scene.Edges = finalizeEdges(scene.Edges)
	scene.Statuses = presentStatuses(nodes)
	return scene, nil
}

// candidateNodes makes a node for every managed resource instance of the plan,
// whatever its status, so that edges can resolve their ends.
func candidateNodes(plan *terraform.Plan) []*Node {
	var nodes []*Node
	for _, rc := range plan.ResourceChanges {
		if !drawable(rc) {
			continue
		}
		addr := terraform.ParseAddress(rc.Address)
		provider := providerName(rc.ProviderName, rc.Type)
		nodes = append(nodes, &Node{
			Address:       rc.Address,
			ConfigAddress: addr.ConfigAddress(),
			Type:          rc.Type,
			Name:          addr.Name + addr.Key,
			Provider:      provider,
			Status:        statusOf(rc),
			Family:        Family(rc.Type, provider),
		})
	}
	return nodes
}

// providerName returns the short provider name: the last path segment of the
// registry address, or the type prefix when the address is empty.
func providerName(full, resourceType string) string {
	full = strings.Trim(full, `"]`)
	if i := strings.LastIndex(full, "/"); i >= 0 {
		full = full[i+1:]
	}
	if full == "" {
		full, _, _ = strings.Cut(resourceType, "_")
	}
	return full
}

func indexByConfig(nodes []*Node) map[string][]*Node {
	m := make(map[string][]*Node)
	for _, n := range nodes {
		m[n.ConfigAddress] = append(m[n.ConfigAddress], n)
	}
	for _, list := range m {
		sort.Slice(list, func(i, j int) bool { return list[i].Address < list[j].Address })
	}
	return m
}

// presentStatuses lists the statuses used by the nodes, in fixed order.
func presentStatuses(nodes []*Node) []Status {
	var seen [StatusNoop + 1]bool
	for _, n := range nodes {
		seen[n.Status] = true
	}
	var out []Status
	for s := StatusCreate; s <= StatusNoop; s++ {
		if seen[s] {
			out = append(out, s)
		}
	}
	return out
}

// buildContainers nests nodes in provider > module instance containers, sorted by
// label, and assigns container ids in depth-first order.
func buildContainers(nodes []*Node) []*Container {
	providers := make(map[string]*Container)
	for _, n := range nodes {
		pc := providers[n.Provider]
		if pc == nil {
			pc = &Container{Kind: KindProvider, Label: n.Provider}
			providers[n.Provider] = pc
		}
		cur := pc
		path := ""
		cfg := ""
		for _, step := range terraform.ParseAddress(n.Address).Modules {
			label := "module." + step.Name + step.Key
			path = joinPath(path, label)
			cfg = joinPath(cfg, "module."+step.Name)
			cur = childContainer(cur, label, path, cfg)
		}
		cur.Nodes = append(cur.Nodes, n)
	}
	out := make([]*Container, 0, len(providers))
	for _, c := range providers {
		out = append(out, c)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Label < out[j].Label })
	next := 0
	for _, c := range out {
		sortContainer(c, &next)
	}
	return out
}

func joinPath(base, step string) string {
	if base == "" {
		return step
	}
	return base + "." + step
}

// childContainer finds or creates the module container with the given label.
func childContainer(parent *Container, label, path, cfg string) *Container {
	for _, c := range parent.Children {
		if c.Label == label {
			return c
		}
	}
	c := &Container{Kind: KindModule, Label: label, Path: path, configPath: cfg, provider: parent.providerName()}
	parent.Children = append(parent.Children, c)
	return c
}

func (c *Container) providerName() string {
	if c.Kind == KindProvider {
		return c.Label
	}
	return c.provider
}

func sortContainer(c *Container, next *int) {
	*next++
	c.ID = "c" + strconv.Itoa(*next)
	sort.Slice(c.Nodes, func(i, j int) bool { return c.Nodes[i].Address < c.Nodes[j].Address })
	sort.Slice(c.Children, func(i, j int) bool { return c.Children[i].Label < c.Children[j].Label })
	for _, ch := range c.Children {
		sortContainer(ch, next)
	}
}
