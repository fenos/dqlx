package deku

import (
	"bytes"
	"strings"
)

type selectionSet struct {
	Fields []string
	Parent edge
	Edges  map[string][]QueryBuilder
}

func (selection selectionSet) ToDQL() (query string, args []interface{}, err error) {
	writer := bytes.Buffer{}

	fieldNames := make([]string, len(selection.Fields))
	for index, fieldName := range selection.Fields {
		fieldNames[index] = fieldName
	}

	writer.WriteString(strings.Join(fieldNames, " "))

	// nested edges
	nestedEdges, ok := selection.Edges[selection.Parent.Name]

	if !ok {
		return writer.String(), args, nil
	}

	for _, edge := range nestedEdges {
		nestedEdge, edgesArgs, err := edge.toDql()
		args = append(args, edgesArgs...)

		if err != nil {
			return "", nil, err
		}
		writer.WriteString(" " + nestedEdge)
	}

	return writer.String(), args, nil
}
