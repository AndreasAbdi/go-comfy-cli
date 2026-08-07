package workflows

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

// Prompt returns the API-format graph used by ComfyUI's /prompt endpoint.
// API-format files are passed through. Desktop-saved frontend workflows are
// flattened, including their embedded subgraph definitions.
func Prompt(workflow Workflow) (map[string]any, error) {
	var definition map[string]any
	if err := json.Unmarshal(workflow.Definition, &definition); err != nil {
		return nil, err
	}
	if isPrompt(definition) {
		return definition, nil
	}

	graph, err := decodeGraph(definition)
	if err != nil {
		return nil, err
	}
	result, _, err := flattenGraph(graph, nil)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, errors.New("workflow contains no executable nodes")
	}
	return result, nil
}

func isPrompt(definition map[string]any) bool {
	if len(definition) == 0 {
		return false
	}
	for _, value := range definition {
		node, ok := value.(map[string]any)
		if !ok {
			return false
		}
		if _, ok := node["class_type"].(string); !ok {
			return false
		}
		if _, ok := node["inputs"].(map[string]any); !ok {
			return false
		}
	}
	return true
}

type graph struct {
	Nodes     []node
	Links     map[int]link
	Subgraphs map[string]graph
	Inputs    []port
	Outputs   []port
}

type node struct {
	ID            int
	Type          string
	Mode          int
	Inputs        []input
	WidgetsValues []any
}

type input struct {
	Name   string
	Link   *int
	Widget bool
	Type   string
}

type link struct {
	OriginID   int
	OriginSlot int
}

type port struct {
	Name    string
	LinkIDs []int
}

type value struct {
	Reference []any
	Literal   any
	HasValue  bool
}

func decodeGraph(definition map[string]any) (graph, error) {
	return decodeGraphObject(definition)
}

func decodeGraphObject(object map[string]any) (graph, error) {
	result := graph{
		Links:     make(map[int]link),
		Subgraphs: make(map[string]graph),
	}

	for _, rawNode := range asSlice(object["nodes"]) {
		nodeObject, ok := rawNode.(map[string]any)
		if !ok {
			continue
		}
		id, ok := asInt(nodeObject["id"])
		if !ok {
			return graph{}, errors.New("workflow node has an invalid id")
		}
		current := node{ID: id, Type: asString(nodeObject["type"]), Mode: asIntDefault(nodeObject["mode"], 0)}
		for _, rawInput := range asSlice(nodeObject["inputs"]) {
			inputObject, ok := rawInput.(map[string]any)
			if !ok {
				continue
			}
			var linkID *int
			if candidate, ok := asInt(inputObject["link"]); ok {
				linkID = &candidate
			}
			widget, hasWidget := inputObject["widget"]
			widgetObject, _ := widget.(map[string]any)
			current.Inputs = append(current.Inputs, input{
				Name:   asString(inputObject["name"]),
				Link:   linkID,
				Widget: hasWidget && widgetObject != nil,
				Type:   asString(inputObject["type"]),
			})
		}
		current.WidgetsValues = asSlice(nodeObject["widgets_values"])
		result.Nodes = append(result.Nodes, current)
	}

	for _, rawLink := range asSlice(object["links"]) {
		var id, originID, originSlot int
		var idOK, originOK, slotOK bool
		if parts := asSlice(rawLink); len(parts) >= 4 {
			id, idOK = asInt(parts[0])
			originID, originOK = asInt(parts[1])
			originSlot, slotOK = asInt(parts[2])
		} else if linkObject, ok := rawLink.(map[string]any); ok {
			id, idOK = asInt(linkObject["id"])
			originID, originOK = asInt(linkObject["origin_id"])
			originSlot, slotOK = asInt(linkObject["origin_slot"])
		}
		if idOK && originOK && slotOK {
			result.Links[id] = link{OriginID: originID, OriginSlot: originSlot}
		}
	}

	definitions, _ := object["definitions"].(map[string]any)
	for _, rawSubgraph := range asSlice(definitions["subgraphs"]) {
		subgraphObject, ok := rawSubgraph.(map[string]any)
		if !ok {
			continue
		}
		id := asString(subgraphObject["id"])
		if id == "" {
			continue
		}
		subgraph, err := decodeGraphObject(subgraphObject)
		if err != nil {
			return graph{}, err
		}
		for _, rawPort := range asSlice(subgraphObject["inputs"]) {
			portObject, ok := rawPort.(map[string]any)
			if !ok {
				continue
			}
			subgraph.Inputs = append(subgraph.Inputs, port{
				Name:    asString(portObject["name"]),
				LinkIDs: asInts(portObject["linkIds"]),
			})
		}
		for _, rawPort := range asSlice(subgraphObject["outputs"]) {
			portObject, ok := rawPort.(map[string]any)
			if !ok {
				continue
			}
			subgraph.Outputs = append(subgraph.Outputs, port{
				Name:    asString(portObject["name"]),
				LinkIDs: asInts(portObject["linkIds"]),
			})
		}
		result.Subgraphs[id] = subgraph
	}

	return result, nil
}

func flattenGraph(current graph, external []value) (map[string]any, []value, error) {
	result := make(map[string]any)
	outputs := make([]value, 0, len(current.Outputs))
	nodeOutputs := make(map[int][]value)
	orderedNodes := make([]node, 0, len(current.Nodes))
	for _, currentNode := range current.Nodes {
		if _, ok := current.Subgraphs[currentNode.Type]; ok {
			orderedNodes = append(orderedNodes, currentNode)
		}
	}
	for _, currentNode := range current.Nodes {
		if isPrimitiveValueNode(currentNode.Type) {
			if resolved, ok := primitiveOutputValue(current, currentNode, nodeOutputs, external); ok {
				nodeOutputs[currentNode.ID] = []value{resolved}
			}
		}
		if currentNode.Mode == 4 {
			for _, currentInput := range currentNode.Inputs {
				if currentInput.Link == nil {
					continue
				}
				if resolved, ok := resolveLink(current, *currentInput.Link, nil, nodeOutputs, external); ok {
					nodeOutputs[currentNode.ID] = []value{resolved}
					break
				}
			}
		}
	}
	for _, currentNode := range current.Nodes {
		if _, ok := current.Subgraphs[currentNode.Type]; !ok {
			orderedNodes = append(orderedNodes, currentNode)
		}
	}

	for _, currentNode := range orderedNodes {
		if currentNode.Mode == 4 || currentNode.Mode == 2 {
			continue
		}
		if isVirtualNode(currentNode.Type) {
			if isPrimitiveValueNode(currentNode.Type) {
				if resolved, ok := primitiveOutputValue(current, currentNode, nodeOutputs, external); ok {
					nodeOutputs[currentNode.ID] = []value{resolved}
				}
			}
			continue
		}

		if subgraph, ok := current.Subgraphs[currentNode.Type]; ok {
			externalInputs := make([]value, len(subgraph.Inputs))
			for index := range subgraph.Inputs {
				if index >= len(currentNode.Inputs) {
					continue
				}
				externalInputs[index] = inputValue(current, currentNode, currentNode.Inputs[index], outputs, nodeOutputs, external)
			}

			inner, innerOutputs, err := flattenGraph(subgraph, externalInputs)
			if err != nil {
				return nil, nil, err
			}
			for id, nodeValue := range inner {
				result[id] = nodeValue
			}
			nodeOutputs[currentNode.ID] = innerOutputs
			continue
		}

		inputs := make(map[string]any)
		for _, currentInput := range currentNode.Inputs {
			if currentInput.Link != nil {
				resolved, ok := resolveLink(current, *currentInput.Link, outputs, nodeOutputs, external)
				if ok {
					if len(resolved.Reference) == 2 {
						inputs[currentInput.Name] = resolved.Reference
					} else if resolved.HasValue {
						inputs[currentInput.Name] = resolved.Literal
					}
				} else if resolved.HasValue {
					inputs[currentInput.Name] = resolved.Literal
				}
				continue
			}
			if currentInput.Widget {
				if widgetValue, ok := widgetValue(currentNode, currentInput); ok {
					inputs[currentInput.Name] = widgetValue
				}
			}
		}
		if currentNode.Type == "" {
			return nil, nil, fmt.Errorf("workflow node %d has no type", currentNode.ID)
		}
		result[strconv.Itoa(currentNode.ID)] = map[string]any{
			"class_type": currentNode.Type,
			"inputs":     inputs,
		}
	}
	for index, port := range current.Outputs {
		if len(port.LinkIDs) == 0 {
			continue
		}
		resolved, ok := resolveLink(current, port.LinkIDs[0], outputs, nodeOutputs, external)
		if !ok {
			continue
		}
		if index >= len(outputs) {
			outputs = append(outputs, make([]value, index-len(outputs)+1)...)
		}
		outputs[index] = resolved
	}

	return result, outputs, nil
}

// primitiveOutputValue resolves a Primitive*-node's output. A promoted
// subgraph widget wires the primitive's own "value" widget input to the
// subgraph's -10 boundary link, so an incoming link must take priority over
// the node's own (possibly stale) baked-in widgets_values.
func primitiveOutputValue(current graph, currentNode node, nodeOutputs map[int][]value, external []value) (value, bool) {
	for _, currentInput := range currentNode.Inputs {
		if currentInput.Link == nil {
			continue
		}
		if resolved, ok := resolveLink(current, *currentInput.Link, nil, nodeOutputs, external); ok {
			return resolved, true
		}
	}
	if len(currentNode.WidgetsValues) > 0 {
		return value{Literal: currentNode.WidgetsValues[0], HasValue: true}, true
	}
	return value{}, false
}

func isVirtualNode(nodeType string) bool {
	switch nodeType {
	case "MarkdownNote", "Note", "Reroute":
		return true
	default:
		return isPrimitiveValueNode(nodeType)
	}
}

func isPrimitiveValueNode(nodeType string) bool {
	switch nodeType {
	case "PrimitiveNode", "PrimitiveFloat", "PrimitiveStringMultiline":
		return true
	default:
		return false
	}
}

func inputValue(current graph, currentNode node, currentInput input, outputs []value, nodeOutputs map[int][]value, external []value) value {
	if currentInput.Link != nil {
		if resolved, ok := resolveLink(current, *currentInput.Link, outputs, nodeOutputs, external); ok {
			return resolved
		}
	}
	if currentInput.Widget {
		if widgetValue, ok := widgetValue(currentNode, currentInput); ok {
			return value{Literal: widgetValue, HasValue: true}
		}
	}
	return value{}
}

func resolveLink(current graph, linkID int, outputs []value, nodeOutputs map[int][]value, external []value) (value, bool) {
	connection, ok := current.Links[linkID]
	if !ok {
		return value{}, false
	}
	if connection.OriginID == -10 {
		if connection.OriginSlot >= 0 && connection.OriginSlot < len(external) {
			return external[connection.OriginSlot], external[connection.OriginSlot].HasValue || len(external[connection.OriginSlot].Reference) == 2
		}
		return value{}, false
	}
	if connection.OriginID == -20 {
		if connection.OriginSlot >= 0 && connection.OriginSlot < len(outputs) {
			return outputs[connection.OriginSlot], true
		}
		return value{}, false
	}
	if refs, ok := nodeOutputs[connection.OriginID]; ok {
		if connection.OriginSlot >= 0 && connection.OriginSlot < len(refs) {
			return refs[connection.OriginSlot], true
		}
		return value{}, false
	}
	return value{Reference: []any{strconv.Itoa(connection.OriginID), connection.OriginSlot}, HasValue: true}, true
}

func widgetValue(currentNode node, currentInput input) (any, bool) {
	index := 0
	widgetValues := currentNode.WidgetsValues
	if hasSeedControlWidget(currentNode) {
		widgetValues = withoutSeedControlWidget(currentNode, widgetValues)
	}
	for _, candidate := range currentNode.Inputs {
		if !candidate.Widget {
			continue
		}
		if index >= len(widgetValues) {
			return nil, false
		}
		value := widgetValues[index]
		index++
		if candidate.Name != currentInput.Name {
			continue
		}
		if candidate.Type == "IMAGEUPLOAD" {
			return nil, false
		}
		if _, isArray := value.([]any); isArray {
			return map[string]any{"__value__": value}, true
		}
		return value, true
	}
	return nil, false
}

func hasSeedControlWidget(currentNode node) bool {
	if currentNode.Type != "KSampler" && currentNode.Type != "KSamplerAdvanced" && currentNode.Type != "TextEncodeAceStepAudio1.5" {
		return false
	}
	seedIndex := seedWidgetIndex(currentNode)
	return seedIndex >= 0 && seedIndex+1 < len(currentNode.WidgetsValues)
}

func seedWidgetIndex(currentNode node) int {
	index := 0
	for _, currentInput := range currentNode.Inputs {
		if !currentInput.Widget {
			continue
		}
		if currentInput.Name == "seed" {
			return index
		}
		index++
	}
	return -1
}

func withoutSeedControlWidget(currentNode node, values []any) []any {
	seedIndex := seedWidgetIndex(currentNode)
	if seedIndex < 0 || seedIndex+1 >= len(values) {
		return values
	}
	result := make([]any, 0, len(values)-1)
	result = append(result, values[:seedIndex+1]...)
	result = append(result, values[seedIndex+2:]...)
	return result
}

func asSlice(value any) []any {
	result, _ := value.([]any)
	return result
}

func asInts(value any) []int {
	var result []int
	for _, item := range asSlice(value) {
		if number, ok := asInt(item); ok {
			result = append(result, number)
		}
	}
	return result
}

func asInt(value any) (int, bool) {
	switch number := value.(type) {
	case float64:
		return int(number), number == float64(int(number))
	case int:
		return number, true
	case json.Number:
		parsed, err := strconv.Atoi(number.String())
		return parsed, err == nil
	default:
		return 0, false
	}
}

func asIntDefault(value any, fallback int) int {
	if result, ok := asInt(value); ok {
		return result
	}
	return fallback
}

func asString(value any) string {
	result, _ := value.(string)
	return result
}
