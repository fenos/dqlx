package deku

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"time"
)

type rootOperation struct {
	operations []queryGrammar
}

func BatchQuery(queries ...QueryBuilder) (query string, args map[string]interface{}, err error) {
	mainOperation := rootOperation{}

	for _, query := range queries {
		mainOperation.operations = append(mainOperation.operations, query.toGrammar())
	}

	return mainOperation.ToQuery()
}

func (rootOperation rootOperation) ToQuery() (query string, variables map[string]interface{}, err error) {
	blocNames := make([]string, len(rootOperation.operations))

	for index, block := range rootOperation.operations {
		blocNames[index] = strings.Title(strings.ToLower(block.Name()))
	}

	queryName := strings.Join(blocNames, "_")

	var args []interface{}
	var statements []string

	for _, block := range rootOperation.operations {
		statement, queryArg, err := block.ToDQL()

		if err != nil {
			return "", nil, err
		}

		statements = append(statements, statement)
		args = append(args, queryArg...)
	}

	innerQuery := strings.Join(statements, " ")

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

		// Assign the variables
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
