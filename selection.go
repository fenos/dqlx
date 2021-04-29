package deku

import (
	"bytes"
	"strings"
)

type selectionSet struct {
	Fields          []string
	ParentName      string
	Edges           map[string][]QueryBuilder
	HasParentFields bool
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
	nestedEdges, ok := selection.Edges[selection.ParentName]

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

	// add a space if parent fields are present
	if selection.HasParentFields {
		writer.WriteString(" ")
	}

	writer.WriteString(strings.Join(statements, " "))

	return writer.String(), args, nil
}

type fields struct {
	predicates []string
}

func Fields(fieldNames ...string) DQLizer {
	var allFields []string

	for _, field := range fieldNames {
		combinedFields := strings.Fields(field)
		allFields = append(allFields, combinedFields...)
	}

	return fields{
		predicates: allFields,
	}
}

func (fields fields) ToDQL() (query string, args []interface{}, err error) {
	return strings.Join(fields.predicates, " "), nil, nil
}

func FieldList(fields []string) string {
	return strings.Join(fields, symbolEdgeTraversal)
}

func ParseFields(fields string) []string {
	return strings.Fields(fields)
}
