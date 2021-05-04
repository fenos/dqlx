package dqlx

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

// ToDQL returns the DQL statements for representing a selection set
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
		combinedFields := ParseFields(field)
		allFields = append(allFields, combinedFields...)
	}

	return fields{
		predicates: allFields,
	}
}

func (fields fields) ToDQL() (query string, args []interface{}, err error) {
	return strings.Join(fields.predicates, " "), nil, nil
}

func ParseFields(fields string) []string {
	var parsedFields []string
	fieldsParts := strings.Split(fields, "\n")

	for _, fieldPart := range fieldsParts {
		if fieldPart == "" {
			continue
		}
		escapedField := Minify(escapeField(fieldPart))
		parsedFields = append(parsedFields, escapedField)
	}
	return parsedFields
}

func escapeField(field string) string {
	removeCharacters := []string{"{", "}"}

	escapedField := field
	for _, char := range removeCharacters {
		escapedField = strings.ReplaceAll(field, char, "")
	}
	return escapedField
}
