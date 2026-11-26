package terraform

import (
	"encoding/json"
	"sort"
)

// Configuration is the optional `configuration` section of a plan. It holds the
// static relations between resources (references and depends_on) that the
// diagram reads. It never holds attribute values of the plan.
type Configuration struct {
	RootModule ConfigModule `json:"root_module"`
}

// ConfigModule is a module body of the configuration.
type ConfigModule struct {
	Resources   []ConfigResource      `json:"resources"`
	ModuleCalls map[string]ModuleCall `json:"module_calls"`
}

// ModuleCall is a call to a child module; only its nested body is read.
type ModuleCall struct {
	Module ConfigModule `json:"module"`
}

// ConfigResource is a resource block of the configuration. Address is index-free
// and relative to its module.
type ConfigResource struct {
	Address     string         `json:"address"`
	Mode        string         `json:"mode"`
	Type        string         `json:"type"`
	Name        string         `json:"name"`
	DependsOn   []string       `json:"depends_on"`
	Expressions map[string]any `json:"expressions"`
}

// UnmarshalJSON is lenient on purpose: the configuration is an optional extra, so
// a section that does not match the expected shape is dropped instead of making a
// previously valid plan fail to parse.
func (c *Configuration) UnmarshalJSON(data []byte) error {
	type plain Configuration
	var p plain
	if err := json.Unmarshal(data, &p); err != nil {
		*c = Configuration{}
		return nil
	}
	*c = Configuration(p)
	return nil
}

// References returns every reference string found in the resource expressions
// and in depends_on, without duplicates, in lexical order.
func (r ConfigResource) References() []string {
	seen := make(map[string]bool)
	for _, d := range r.DependsOn {
		seen[d] = true
	}
	for _, v := range r.Expressions {
		collectReferences(v, seen)
	}
	out := make([]string, 0, len(seen))
	for s := range seen {
		out = append(out, s)
	}
	sort.Strings(out)
	return out
}

// collectReferences walks an expression value and gathers the strings of any
// `references` array. Unexpected shapes are ignored.
func collectReferences(v any, seen map[string]bool) {
	switch t := v.(type) {
	case map[string]any:
		for k, child := range t {
			if k == "references" {
				if list, ok := child.([]any); ok {
					for _, item := range list {
						if s, ok := item.(string); ok {
							seen[s] = true
						}
					}
					continue
				}
			}
			collectReferences(child, seen)
		}
	case []any:
		for _, child := range t {
			collectReferences(child, seen)
		}
	}
}
