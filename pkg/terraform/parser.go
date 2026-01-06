package terraform

import (
	"encoding/json"
	"fmt"
	"os"
)

// Parser handles parsing of Terraform plan files.
type Parser struct{}

// NewParser creates a new Parser instance.
func NewParser() *Parser {
	return &Parser{}
}

// ParseFile reads and parses a Terraform plan JSON file.
func (p *Parser) ParseFile(path string) (*Plan, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read plan file: %w", err)
	}

	return p.Parse(data)
}

// Parse parses Terraform plan JSON data.
func (p *Parser) Parse(data []byte) (*Plan, error) {
	var plan Plan
	if err := json.Unmarshal(data, &plan); err != nil {
		return nil, fmt.Errorf("failed to parse plan JSON: %w", err)
	}

	return &plan, nil
}

// FilterByType filters resource changes by resource type.
func FilterByType(changes []ResourceChange, types []string) []ResourceChange {
	if len(types) == 0 {
		return changes
	}

	typeMap := make(map[string]bool)
	for _, t := range types {
		typeMap[t] = true
	}

	filtered := make([]ResourceChange, 0)
	for _, rc := range changes {
		if typeMap[rc.Type] {
			filtered = append(filtered, rc)
		}
	}

	return filtered
}

// FilterByAction filters resource changes by action type.
func FilterByAction(changes []ResourceChange, actions []ActionType) []ResourceChange {
	if len(actions) == 0 {
		return changes
	}

	actionMap := make(map[ActionType]bool)
	for _, a := range actions {
		actionMap[a] = true
	}

	filtered := make([]ResourceChange, 0)
	for _, rc := range changes {
		action := rc.Change.GetPrimaryAction()
		if actionMap[action] {
			filtered = append(filtered, rc)
		}
	}

	return filtered
}
