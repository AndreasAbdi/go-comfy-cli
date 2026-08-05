package workflows

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// Mapping is the versioned YAML contract used by --args-file.
type Mapping struct {
	Version int              `yaml:"version"`
	Aliases map[string]Alias `yaml:"aliases"`
}

// Alias maps a friendly CLI name to a jq selector.
type Alias struct {
	Selector        string `yaml:"selector"`
	Select          string `yaml:"select"`
	IndexedSelector string `yaml:"indexed_selector"`
	Type            string `yaml:"type"`
	ValueType       string `yaml:"value_type"`
	Cardinality     string `yaml:"cardinality"`
	Match           string `yaml:"match"`
	Operation       string `yaml:"operation,omitempty"`
}

// AliasReference identifies an alias, optionally selecting one indexed match.
type AliasReference struct {
	Name  string
	Index *int
	All   bool
}

// Assignment is one parsed alias=value command.
type Assignment struct {
	Reference AliasReference
	Value     any
}

// LoadMapping reads and validates a version 1 args mapping file.
func LoadMapping(path string) (Mapping, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Mapping{}, fmt.Errorf("read args file: %w", err)
	}

	var mapping Mapping
	if err := yaml.Unmarshal(data, &mapping); err != nil {
		return Mapping{}, fmt.Errorf("parse args file: %w", err)
	}
	if mapping.Version != 1 {
		return Mapping{}, fmt.Errorf("unsupported args mapping version %d", mapping.Version)
	}
	if len(mapping.Aliases) == 0 {
		return Mapping{}, errors.New("args mapping does not define any aliases")
	}

	for name, alias := range mapping.Aliases {
		if strings.TrimSpace(name) == "" {
			return Mapping{}, errors.New("args mapping contains an empty alias")
		}
		if strings.TrimSpace(alias.selector()) == "" {
			return Mapping{}, fmt.Errorf("alias %q has no selector", name)
		}
		if alias.cardinality() == "indexed" && strings.TrimSpace(alias.IndexedSelector) == "" {
			return Mapping{}, fmt.Errorf("alias %q uses indexed cardinality but has no indexed_selector", name)
		}
	}

	return mapping, nil
}

// ParseAliasAssignment parses ALIAS=VALUE or ALIAS::VALUE using mapping types.
func ParseAliasAssignment(raw string, mapping Mapping) (Assignment, error) {
	left, right, found := strings.Cut(raw, "=")
	if !found {
		left, right, found = strings.Cut(raw, "::")
	}
	if !found {
		return Assignment{}, fmt.Errorf("%q must use ALIAS=VALUE or ALIAS::VALUE syntax", raw)
	}

	reference, err := ParseAliasReference(strings.TrimSpace(left))
	if err != nil {
		return Assignment{}, err
	}

	alias, exists := mapping.Aliases[reference.Name]
	if !exists {
		return Assignment{}, fmt.Errorf("unknown alias %q", reference.Name)
	}

	value, err := parseTypedValue(alias.valueType(), right)
	if err != nil {
		return Assignment{}, fmt.Errorf("alias %q: %w", reference.Name, err)
	}

	return Assignment{Reference: reference, Value: value}, nil
}

// ParseAliasReference parses alias, alias[0], or alias[].
func ParseAliasReference(raw string) (AliasReference, error) {
	if strings.HasSuffix(raw, "[]") {
		name := strings.TrimSpace(strings.TrimSuffix(raw, "[]"))
		if name == "" {
			return AliasReference{}, errors.New("alias name is empty")
		}
		return AliasReference{Name: name, All: true}, nil
	}

	open := strings.LastIndex(raw, "[")
	if open == -1 || !strings.HasSuffix(raw, "]") {
		name := strings.TrimSpace(raw)
		if name == "" {
			return AliasReference{}, errors.New("alias name is empty")
		}
		return AliasReference{Name: name}, nil
	}

	name := strings.TrimSpace(raw[:open])
	indexText := raw[open+1 : len(raw)-1]
	if name == "" {
		return AliasReference{}, errors.New("alias name is empty")
	}
	index, err := strconv.Atoi(indexText)
	if err != nil || index < 0 {
		return AliasReference{}, fmt.Errorf("invalid alias index %q", indexText)
	}

	return AliasReference{Name: name, Index: &index}, nil
}

// AliasOperations converts assignments into transpiler operations.
func (mapping Mapping) AliasOperations(rawAssignments []string) ([]Operation, error) {
	operations := make([]Operation, 0, len(rawAssignments))
	for _, raw := range rawAssignments {
		assignment, err := ParseAliasAssignment(raw, mapping)
		if err != nil {
			return nil, err
		}
		alias := mapping.Aliases[assignment.Reference.Name]
		selector := alias.selector()
		variables := map[string]any{}
		cardinality := alias.cardinality()
		if assignment.Reference.Index != nil {
			if strings.TrimSpace(alias.IndexedSelector) == "" {
				return nil, fmt.Errorf("alias %q does not define indexed_selector", assignment.Reference.Name)
			}
			selector = alias.IndexedSelector
			variables["index"] = *assignment.Reference.Index
			cardinality = "indexed"
		}
		if assignment.Reference.All && cardinality == "one" {
			cardinality = "many"
		}

		operations = append(operations, Operation{
			Selector:    selector,
			Value:       assignment.Value,
			Cardinality: cardinality,
			Variables:   variables,
			Description: raw,
		})
	}
	return operations, nil
}

func (alias Alias) selector() string {
	if strings.TrimSpace(alias.Selector) != "" {
		return alias.Selector
	}
	return alias.Select
}

func (alias Alias) valueType() string {
	if strings.TrimSpace(alias.Type) != "" {
		return strings.ToLower(strings.TrimSpace(alias.Type))
	}
	return strings.ToLower(strings.TrimSpace(alias.ValueType))
}

func (alias Alias) cardinality() string {
	if strings.TrimSpace(alias.Cardinality) != "" {
		return strings.ToLower(strings.TrimSpace(alias.Cardinality))
	}
	return strings.ToLower(strings.TrimSpace(alias.Match))
}

func parseTypedValue(kind, raw string) (any, error) {
	switch kind {
	case "", "string":
		return raw, nil
	case "number":
		var value json.Number
		decoder := json.NewDecoder(strings.NewReader(strings.TrimSpace(raw)))
		decoder.UseNumber()
		if err := decoder.Decode(&value); err != nil {
			return nil, fmt.Errorf("expected a number: %w", err)
		}
		return value, nil
	case "integer":
		value, err := strconv.ParseInt(strings.TrimSpace(raw), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("expected an integer: %w", err)
		}
		return value, nil
	case "boolean":
		value, err := strconv.ParseBool(strings.TrimSpace(raw))
		if err != nil {
			return nil, fmt.Errorf("expected a boolean: %w", err)
		}
		return value, nil
	case "json":
		var value any
		if err := json.Unmarshal([]byte(strings.TrimSpace(raw)), &value); err != nil {
			return nil, fmt.Errorf("expected JSON: %w", err)
		}
		return value, nil
	default:
		return nil, fmt.Errorf("unsupported alias type %q", kind)
	}
}
