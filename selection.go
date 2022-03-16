package dqlx

import (
	"bytes"
	"fmt"
	"strings"
)

type node struct {
	Attributes          DQLizer
	ParentName          string
	Edges               map[string][]QueryBuilder
	HasParentAttributes bool
}

// ToDQL returns the DQL statements for representing a selection set
func (node node) ToDQL() (query string, args Args, err error) {
	writer := bytes.Buffer{}

	if node.Attributes != nil {
		if err := addPart(node.Attributes, &writer, &args); err != nil {
			return "", nil, err
		}
	}

	// nested childrenEdges
	nestedEdges, ok := node.Edges[node.ParentName]

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

	// add a space if parent nodeAttributes are present
	if node.HasParentAttributes {
		writer.WriteString(" ")
	}

	writer.WriteString(strings.Join(statements, " "))

	return writer.String(), args, nil
}

type nodeAttributes struct {
	predicates []interface{}
}

// Select adds nodeAttributes to selection set
func Select(predicates ...interface{}) DQLizer {
	return nodeAttributes{predicates}
}

// Fields alias of Select
// @Deprecated use Select() instead
func Fields(predicates ...interface{}) DQLizer {
	return Select(predicates...)
}

// ToDQL returns the dql statement for selected nodeAttributes
func (fields nodeAttributes) ToDQL() (query string, args Args, err error) {
	var selectedFields []string

	for _, field := range fields.predicates {
		switch requestField := field.(type) {
		case DQLizer:
			computedDql, computedArgs, err := requestField.ToDQL()

			if err != nil {
				return "", nil, err
			}

			args = append(args, computedArgs...)
			selectedFields = append(selectedFields, computedDql)
		case string:
			fieldString := parsePredicates(requestField)
			selectedFields = append(selectedFields, fieldString...)
		default:
			return "", nil, fmt.Errorf("nodeAttributes can only accept strings or Dqlizer, given %v", requestField)
		}
	}

	return strings.Join(selectedFields, " "), args, nil
}

type aliasField struct {
	alias string
	value interface{}
}

// Alias allows to alias a field
func Alias(alias string, predicate interface{}) DQLizer {
	return aliasField{
		alias: alias,
		value: predicate,
	}
}

// ToDQL returns the alias dql statement of a field
func (aliasField aliasField) ToDQL() (query string, args Args, err error) {
	var value string

	switch cast := aliasField.value.(type) {
	case DQLizer:
		value, args, err = cast.ToDQL()

		if err != nil {
			return "", nil, err
		}
	case string:
		value = EscapePredicate(cast)
	default:
		return "", nil, fmt.Errorf("alias only accepts  string or DQlizers, given %v", value)
	}

	aliasName := EscapePredicate(aliasField.alias)

	return fmt.Sprintf("%s:%s", aliasName, value), args, nil
}

type as struct {
	variable  string
	predicate interface{}
}

// As makes a field a variable
// Example: dqlx.Query(...).Select(dqlx.As("C", "a"))
func As(varName string, predicate interface{}) DQLizer {
	return as{
		variable:  varName,
		predicate: predicate,
	}
}

// ToDQL returns the dql statement for a field variable
func (as as) ToDQL() (query string, args Args, err error) {
	var predicate string

	switch cast := as.predicate.(type) {
	case DQLizer:
		predicate, args, err = cast.ToDQL()

		if err != nil {
			return "", nil, err
		}
	case string:
		predicate = EscapePredicate(cast)
	default:
		return "", nil, fmt.Errorf("alias only accepts  string or DQlizers, given %v", predicate)
	}

	varName := escapeSpecialChars(as.variable)

	return fmt.Sprintf("%s AS %s", varName, predicate), args, nil
}

func parsePredicates(predicates string) []string {
	var parsedFields []string
	fieldsParts := strings.Split(predicates, "\n")

	for _, fieldPart := range fieldsParts {
		if strings.TrimSpace(fieldPart) == "" {
			continue
		}
		escapedField := EscapePredicate(fieldPart)
		parsedFields = append(parsedFields, escapedField)
	}
	return parsedFields
}
