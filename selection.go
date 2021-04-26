package deku

import (
	"bytes"
	"strings"
)

type selectionSet struct {
	Fields []string
	Parent *edge
	Edges  map[string][]QueryBuilder
}

func (selection selectionSet) ToDQL() (query string, args []interface{}, err error) {
	writer := bytes.Buffer{}

	var fieldNames []string
	for _, fieldName := range selection.Fields {
		if fieldName == "" {
			continue
		}
		fieldNames = append(fieldNames, fieldName)
	}

	writer.WriteString(strings.Join(fieldNames, " "))

	// nested childrenEdges
	nestedEdges, ok := selection.Edges[selection.Parent.Name]

	if !ok {
		return writer.String(), args, nil
	}

	statements := make([]string, 0, len(nestedEdges))
	for _, queryBuilder := range nestedEdges {
		nestedEdge, edgesArgs, err := queryBuilder.rootEdge.ToDQL()
		args = append(args, edgesArgs...)

		if err != nil {
			return "", nil, err
		}

		statements = append(statements, nestedEdge)
	}

	// Add a space if parent fields are there
	if selection.Parent != nil && len(selection.Parent.Selection.Fields) > 0 {
		writer.WriteString(" ")
	}

	writer.WriteString(strings.Join(statements, " "))

	return writer.String(), args, nil
}

func FieldList(fields []string) string {
	return strings.Join(fields, symbolEdgeTraversal)
}

func ParseFields(fields string) []string {
	return strings.Fields(fields)
}
