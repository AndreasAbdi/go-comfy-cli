package workflows

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/itchyny/gojq"
)

// Operation describes one JSON replacement operation.
// Selector is a jq expression that identifies the value to replace.
type Operation struct {
	Selector    string
	Value       any
	Cardinality string
	Variables   map[string]any
	Description string
}

// ParseJSONReplacement parses SELECTOR::JSON_VALUE for --replace-json.
func ParseJSONReplacement(raw string) (Operation, error) {
	selector, encodedValue, found := strings.Cut(raw, "::")
	if !found {
		return Operation{}, fmt.Errorf("%q must use SELECTOR::JSON_VALUE syntax", raw)
	}
	selector = strings.TrimSpace(selector)
	if selector == "" {
		return Operation{}, errors.New("selector cannot be empty")
	}

	var value any
	if err := json.Unmarshal([]byte(strings.TrimSpace(encodedValue)), &value); err != nil {
		return Operation{}, fmt.Errorf("invalid JSON value %q: %w", encodedValue, err)
	}

	return Operation{
		Selector:    selector,
		Value:       value,
		Cardinality: "many",
		Description: raw,
	}, nil
}

// ApplyOperations applies jq assignment operations sequentially.
func ApplyOperations(document any, operations []Operation) (any, error) {
	result := document
	for index, operation := range operations {
		updated, err := applyOperation(result, operation)
		if err != nil {
			return nil, fmt.Errorf("operation %d: %w", index+1, err)
		}
		result = updated
	}
	return result, nil
}

// Transpile applies operations to a JSON document and returns the transformed JSON.
func Transpile(data json.RawMessage, operations []Operation) (json.RawMessage, error) {
	var document any
	if err := json.Unmarshal(data, &document); err != nil {
		return nil, fmt.Errorf("parse workflow JSON: %w", err)
	}

	result, err := ApplyOperations(document, operations)
	if err != nil {
		return nil, err
	}

	output, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("encode transformed workflow: %w", err)
	}
	return output, nil
}

func applyOperation(document any, operation Operation) (any, error) {
	if strings.TrimSpace(operation.Selector) == "" {
		return nil, errors.New("selector cannot be empty")
	}

	if err := validateCardinality(document, operation); err != nil {
		return nil, err
	}

	expression := fmt.Sprintf("(%s) = $value", operation.Selector)
	query, err := gojq.Parse(expression)
	if err != nil {
		return nil, fmt.Errorf("parse selector %q: %w", operation.Selector, err)
	}

	variableNames, variableValues := operationVariables(operation)
	code, err := gojq.Compile(query, gojq.WithVariables(variableNames))
	if err != nil {
		return nil, fmt.Errorf("compile selector %q: %w", operation.Selector, err)
	}

	iterator := code.Run(document, variableValues...)
	result, ok := iterator.Next()
	if !ok {
		return nil, errors.New("transformation produced no document")
	}
	if err, ok := result.(error); ok {
		return nil, fmt.Errorf("apply selector %q: %w", operation.Selector, err)
	}
	if _, ok := iterator.Next(); ok {
		return nil, fmt.Errorf("selector %q produced multiple documents", operation.Selector)
	}

	return result, nil
}

func validateCardinality(document any, operation Operation) error {
	cardinality := strings.ToLower(strings.TrimSpace(operation.Cardinality))
	if cardinality == "" {
		return nil
	}

	expression := fmt.Sprintf("[%s]", operation.Selector)
	query, err := gojq.Parse(expression)
	if err != nil {
		return fmt.Errorf("parse selector %q: %w", operation.Selector, err)
	}
	variableNames, variableValues := operationVariables(operation)
	code, err := gojq.Compile(query, gojq.WithVariables(variableNames))
	if err != nil {
		return fmt.Errorf("compile selector %q: %w", operation.Selector, err)
	}

	iterator := code.Run(document, variableValues...)
	result, ok := iterator.Next()
	if !ok {
		return fmt.Errorf("selector %q produced no result", operation.Selector)
	}
	if err, ok := result.(error); ok {
		return fmt.Errorf("evaluate selector %q: %w", operation.Selector, err)
	}

	matches, ok := result.([]any)
	if !ok {
		return fmt.Errorf("selector %q did not produce a match list", operation.Selector)
	}

	switch cardinality {
	case "one", "exactly-one", "exactly_one", "indexed":
		if len(matches) != 1 {
			return fmt.Errorf("selector %q matched %d values, expected exactly one", operation.Selector, len(matches))
		}
	case "many", "all":
		if len(matches) == 0 {
			return fmt.Errorf("selector %q matched no values", operation.Selector)
		}
	case "zero-or-more", "zero_or_more", "optional":
		return nil
	default:
		return fmt.Errorf("unsupported cardinality %q", operation.Cardinality)
	}

	return nil
}

func operationVariables(operation Operation) ([]string, []any) {
	names := []string{"$value"}
	values := []any{operation.Value}
	if index, ok := operation.Variables["index"]; ok {
		names = append(names, "$index")
		values = append(values, index)
	}
	return names, values
}
