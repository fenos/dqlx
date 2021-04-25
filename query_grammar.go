package deku

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"time"
)

type queryGrammar struct {
	RootEdge  DQLizer
	Selection DQLizer
	Filters   []DQLizer
}

func (grammar queryGrammar) ToQuery(name string) (query string, variables map[string]interface{}, err error) {
	queryName := strings.Title(strings.ToLower(name))
	innerQuery, args, err := grammar.ToDQL()

	if err != nil {
		return "", nil, err
	}

	query, variables, err = replacePlaceholders(innerQuery, args)

	if err != nil {
		return
	}

	queryPlaceholderNames := getSortedVariables(variables)
	placeholders := make([]string, len(queryPlaceholderNames))

	for index, placeholderName := range queryPlaceholderNames {
		placeholders[index] = fmt.Sprintf("%s:%s", placeholderName, goTypeToDQLType(variables[placeholderName]))
	}

	writer := bytes.Buffer{}
	writer.WriteString(fmt.Sprintf("query %s(%s) {", queryName, strings.Join(placeholders, ", ")))
	writer.WriteString(" " + query)
	writer.WriteString(" }")

	return writer.String(), variables, nil
}

func (grammar queryGrammar) ToDQL() (query string, args []interface{}, err error) {
	writer := bytes.Buffer{}

	args = []interface{}{}

	if err := grammar.addEdge(&writer, &args); err != nil {
		return "", nil, err
	}

	if err := grammar.addFilters(&writer, &args); err != nil {
		return "", nil, err
	}

	writer.WriteString(" { ")

	if err := grammar.addSelection(&writer, &args); err != nil {
		return "", nil, err
	}

	writer.WriteString(" }")

	return writer.String(), args, nil
}

func (grammar queryGrammar) addEdge(writer *bytes.Buffer, args *[]interface{}) error {
	return addPart(grammar.RootEdge, writer, args)
}

func (grammar queryGrammar) addFilters(writer *bytes.Buffer, args *[]interface{}) error {
	if len(grammar.Filters) == 0 {
		return nil
	}

	writer.WriteString(" @filter(")
	for _, filter := range grammar.Filters {
		if err := addPart(filter, writer, args); err != nil {
			return err
		}
	}

	writer.WriteString(")")
	return nil
}

func (grammar queryGrammar) addSelection(writer *bytes.Buffer, args *[]interface{}) error {
	return addPart(grammar.Selection, writer, args)
}

func addPart(part DQLizer, writer *bytes.Buffer, args *[]interface{}) error {
	statement, statementArgs, err := part.ToDQL()
	*args = append(*args, statementArgs...)

	if err != nil {
		return err
	}

	writer.WriteString(statement)

	return nil
}

func goTypeToDQLType(value interface{}) string {
	switch value.(type) {
	case string:
		return "string"
	case int, int8, int32, int64:
		return "int"
	case float32, float64:
		return "float"
	case bool:
		return "bool"
	case time.Time, *time.Time:
		return "datetime"
	}

	return "string"
}

func replacePlaceholders(query string, args []interface{}) (string, map[string]interface{}, error) {
	variables := map[string]interface{}{}
	buf := &bytes.Buffer{}
	i := 0

	for {
		p := strings.Index(query, "??")
		if p == -1 {
			break
		}

		buf.WriteString(query[:p])
		key := fmt.Sprintf("$%d", i)
		buf.WriteString(key)
		query = query[p+2:]

		// Assign the variable
		variables[key] = args[i]

		i++
	}

	buf.WriteString(query)
	return buf.String(), variables, nil
}

func isListType(val interface{}) bool {
	valVal := reflect.ValueOf(val)
	return valVal.Kind() == reflect.Array || valVal.Kind() == reflect.Slice
}
