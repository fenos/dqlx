package dqlx

import (
	"bytes"
	"fmt"
	"strings"
)

type queryOperation struct {
	operations []edge
	variables  []edge
}

type KVPair struct {
	K string
	V string
}

// QueriesToDQL returns the DQL statement for 1 or more queries
// Example: dqlx.QueriesToDQL(query1,query2,query3)
func QueriesToDQL(queries ...QueryBuilder) (query string, args Args, err error) {
	mainOperation := queryOperation{}
	queries = ensureUniqueQueryNames(queries)

	for _, query := range queries {
		mainOperation.operations = append(mainOperation.operations, query.rootEdge)

		for _, variable := range query.variables {
			mainOperation.variables = append(mainOperation.variables, variable.rootEdge)
		}
	}

	return mainOperation.ToDQL()
}

// ToDQL returns the DQL statement for 1 or more queries
func (grammar queryOperation) ToDQL() (query string, variables Args, err error) {
	variables = []interface{}{}
	blocNames := make([]string, len(grammar.operations))

	for index, block := range grammar.operations {
		blocNames[index] = strings.Title(strings.ToLower(block.GetName()))
	}

	queryName := strings.Join(blocNames, "_")

	var args Args
	var statements []string

	if err := addOperation(grammar.variables, &statements, &args); err != nil {
		return "", nil, err
	}

	if err := addOperation(grammar.operations, &statements, &args); err != nil {
		return "", nil, err
	}

	innerQuery := strings.Join(statements, " ")

	query, rawVariables := replacePlaceholders(innerQuery, args, func(index int, value interface{}) string {
		return fmt.Sprintf("$%d", index)
	})
	variables, placeholders := toVariables(rawVariables)

	writer := bytes.Buffer{}
	writer.WriteString(fmt.Sprintf("query %s(%s) {", queryName, strings.Join(placeholders, ", ")))
	writer.WriteString(" " + query)
	writer.WriteString(" }")

	return writer.String(), variables, nil
}

func ensureUniqueQueryNames(queries []QueryBuilder) []QueryBuilder {
	queryNames := map[string]bool{}
	uniqueQueries := make([]QueryBuilder, len(queries))

	for index, query := range queries {
		if queryNames[query.rootEdge.Name] {
			query = query.Name(fmt.Sprintf("%s_%d", query.rootEdge.Name, index))
		}

		queryNames[query.rootEdge.Name] = true
		uniqueQueries[index] = query
	}

	return uniqueQueries
}

func addOperation(operations []edge, statements *[]string, args *Args) error {
	parts := make([]DQLizer, len(operations))

	for index, operation := range operations {
		parts[index] = operation
	}

	return addStatement(parts, statements, args)
}
