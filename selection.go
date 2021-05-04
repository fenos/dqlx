package dqlx

import (
	"bytes"
	"fmt"
	"strings"
)

type selectionSet struct {
	Fields          DQLizer
	ParentName      string
	Edges           map[string][]QueryBuilder
	HasParentFields bool
}

// ToDQL returns the DQL statements for representing a selection set
func (selection selectionSet) ToDQL() (query string, args []interface{}, err error) {
	writer := bytes.Buffer{}

	if selection.Fields != nil {
		if err := addPart(selection.Fields, &writer, &args); err != nil {
			return "", nil, err
		}
	}

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
	predicates []interface{}
}

func Fields(fieldNames ...interface{}) DQLizer {
	return fields{fieldNames}
}

func (fields fields) ToDQL() (query string, args []interface{}, err error) {
	var selectedFields []string

	for _, field := range fields.predicates {
		switch requestField := field.(type) {
		case computedField:
			computedDql, computedArgs, err := requestField.ToDQL()

			if err != nil {
				return "", nil, err
			}

			args = append(args, computedArgs)
			selectedFields = append(selectedFields, computedDql)
		case string:
			fieldString := parseFields(requestField)
			selectedFields = append(selectedFields, fieldString...)
		default:
			return "", nil, fmt.Errorf("fields can only accept strings or computed values, givem %v", requestField)
		}
	}

	return strings.Join(selectedFields, " "), args, nil
}

type computedField struct {
	alias string
	value DQLizer
}

func Computed(alias string, predicate DQLizer) DQLizer {
	return computedField{
		alias: alias,
		value: predicate,
	}
}

func (computedField computedField) ToDQL() (query string, args []interface{}, err error) {
	computedValue, args, err := computedField.value.ToDQL()

	if err != nil {
		return "", nil, err
	}

	predicate := EscapePredicate(Minify(computedField.alias))

	return fmt.Sprintf("%s:%s", predicate, computedValue), args, nil
}

func parseFields(fields string) []string {
	var parsedFields []string
	fieldsParts := strings.Split(fields, "\n")

	for _, fieldPart := range fieldsParts {
		if strings.TrimSpace(fieldPart) == "" {
			continue
		}
		escapedField := EscapePredicate(fieldPart)
		parsedFields = append(parsedFields, escapedField)
	}
	return parsedFields
}

func splitDirective(predicate string) (string, string) {
	predicateParts := strings.Split(predicate, "@")
	directive := ""

	if len(predicateParts) > 1 {
		predicate = predicateParts[0]
		directive = "@" + strings.Join(predicateParts[1:], "")
	}

	return predicate, directive
}
