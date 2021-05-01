package dqlx

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"time"
)

func toVariables(rawVariables map[int]interface{}) (variables map[string]string, placeholders []string) {
	variables = map[string]string{}
	placeholders = make([]string, len(rawVariables))

	queryPlaceholderNames := getSortedVariables(rawVariables)

	// Format Variables
	for index, placeholderName := range queryPlaceholderNames {
		variableName := fmt.Sprintf("$%d", placeholderName)

		variables[variableName] = toVariableValue(rawVariables[index])
		placeholders[index] = fmt.Sprintf("$%d:%s", placeholderName, goTypeToDQLType(rawVariables[placeholderName]))
	}

	return variables, placeholders
}

func toVariableValue(value interface{}) string {
	switch val := value.(type) {
	case time.Time:
		return val.Format(time.RFC3339)
	case *time.Time:
		return val.Format(time.RFC3339)
	default:
		return fmt.Sprintf("%v", val)
	}
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

func dqlTypeToGoType(value DGraphScalar) string {
	switch value {
	case ScalarPassword, ScalarString, ScalarUID:
		return "string"
	case ScalarInt:
		return "int64"
	case ScalarDateTime:
		return "time.Time"
	case ScalarBool:
		return "bool"
	case ScalarFloat:
		return "float64"
	}

	return string(value)
}

func replacePlaceholders(query string, args []interface{}, transform func(index int, value interface{}) string) (string, map[int]interface{}) {
	variables := map[int]interface{}{}
	buf := &bytes.Buffer{}
	i := 0

	for {
		p := strings.Index(query, "??")
		if p == -1 {
			break
		}

		buf.WriteString(query[:p])
		key := transform(i, args[i])
		buf.WriteString(key)
		query = query[p+2:]

		// Assign the variables
		variables[i] = args[i]

		i++
	}

	buf.WriteString(query)
	return buf.String(), variables
}

func isListType(val interface{}) bool {
	valVal := reflect.ValueOf(val)
	return valVal.Kind() == reflect.Array || valVal.Kind() == reflect.Slice
}

func addStatement(parts []DQLizer, statements *[]string, args *[]interface{}) error {
	for _, block := range parts {
		statement, queryArg, err := block.ToDQL()

		if err != nil {
			return err
		}

		*statements = append(*statements, statement)
		*args = append(*args, queryArg...)
	}

	return nil
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
