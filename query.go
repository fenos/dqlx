package deku

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"time"
)

type QueryBuilder struct {
	name         string
	rootFn       *QueryFunction
	selectionSet *SelectionSet
	filters      []*Filter
	filterGroups []*FilterGroup
	isRootQuery  bool
	depth        int
}

type DQLizer interface {
	toDQL() (query string, args []interface{}, err error)
}

func (builder *QueryBuilder) Name(name string) *QueryBuilder {
	builder.name = name
	return builder
}

func Query(rootFn *QueryFunction) *QueryBuilder {
	builder := &QueryBuilder{
		name:        "query",
		rootFn:      rootFn,
		depth:       0,
		isRootQuery: true,
	}

	return builder
}

func (builder *QueryBuilder) ToDQL() (query string, args map[string]interface{}, err error) {
	variables := map[string]interface{}{}

	anonymousQuery, values, err := builder.toDQL()

	if err != nil {
		return "", nil, err
	}

	queryName := strings.Title(strings.ToLower(builder.name))

	query, err = replacePlaceholders(anonymousQuery, values, variables)

	if err != nil {
		return
	}

	queryPlaceholderNames := getSortedKeys(variables)
	placeholders := make([]string, len(queryPlaceholderNames))

	for index, placeholderName := range queryPlaceholderNames {
		placeholders[index] = fmt.Sprintf("%s: %s", placeholderName, goTypeToDQLType(variables[placeholderName]))
	}

	writer := NewWriter()
	writer.AddLine(fmt.Sprintf("query %s(%s) {", queryName, strings.Join(placeholders, ",")))
	writer.AddLine(query)
	writer.AddLine("}")


	return writer.ToString(), variables, nil
}

func goTypeToDQLType(value interface{}) string  {
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

func replacePlaceholders(query string, args []interface{}, variables map[string]interface{}) (string, error) {
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
		variables[key] = args[i]
		query = query[p+2:]

		i++
	}

	buf.WriteString(query)
	return buf.String(), nil
}

func getSortedKeys(exp map[string]interface{}) []string {
	sortedKeys := make([]string, 0, len(exp))
	for k := range exp {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	return sortedKeys
}

func (builder *QueryBuilder) toDQL() (query string, args []interface{}, err error) {
	writer := NewWriter()

	rootOperation := builder.name

	if builder.rootFn != nil {
		operation, rootFnArgs, err := builder.rootFn.toDQL()
		args = append(args, rootFnArgs...)

		if err != nil {
			return "", nil, err
		}
		rootOperation += fmt.Sprintf("(func: %s)", operation)
	}

	// compose filters
	filterOperations := ""
	if len(builder.filters) > 0 || len(builder.filterGroups) > 0 {
		filterOperations = " @filter("
		for index, filter := range builder.filters {
			if index > 0 {
				filterOperations += " " + filter.connection + " "
			}

			appliedFn, fnArgs, err := filter.function.toDQL()
			args = append(args, fnArgs...)

			if err != nil {
				return "", nil, err
			}

			filterOperations += appliedFn
		}

		for index, filterGroup := range builder.filterGroups {
			groupStatement, groupArgs, err := filterGroup.toDQL()
			args = append(args, groupArgs...)

			if index > 0 || len(builder.filters) > 0 {
				filterOperations += " " + filterGroup.connection + " "
			}

			if err != nil {
				return "", nil, err
			}

			filterOperations += groupStatement
		}
		filterOperations += ")"
	}

	// add filters
	rootOperation += filterOperations

	var mainOpWriter *Writer
	if builder.isRootQuery {
		// Add root operation
		mainOpWriter = writer.AddIndentedLine(rootOperation + " {")
		builder.selectionSet.SetDepth(mainOpWriter.indent)
	} else {
		// add nested operation
		mainOpWriter = writer.AddLine(rootOperation + " {")
		builder.selectionSet.SetDepth(builder.depth + mainOpWriter.indent + 1)
	}

	// Split out the fields
	selection, selectionArgs, err := builder.selectionSet.toDQL()
	args = append(args, selectionArgs...)

	if err != nil {
		return "", nil, err
	}

	mainOpWriter.AddLine(selection)

	if builder.isRootQuery {
		writer.AddIndentedLine("}")
	} else {
		mainOpWriter.indent = builder.depth + mainOpWriter.indent
		mainOpWriter.AddIndentedLine("}")
	}

	return writer.ToString(), args,nil
}

func (builder *QueryBuilder) Edge(name string, fn func(builder *QueryBuilder)) *QueryBuilder {
	nestedBuilder := Query(nil).Name(name)
	nestedBuilder.depth = builder.depth + 1
	nestedBuilder.isRootQuery = false
	fn(nestedBuilder)

	builder.selectionSet.edges = append(builder.selectionSet.edges, nestedBuilder)

	return builder
}

func (builder *QueryBuilder) MergeEdge(edgeQuery *QueryBuilder) *QueryBuilder {
	edgeQuery.isRootQuery = false
	edgeQuery.depth = builder.depth + 1

	builder.selectionSet.edges = append(builder.selectionSet.edges, edgeQuery)

	return builder
}

func (builder *QueryBuilder) Fields(fields ...string) *QueryBuilder {
	if len(fields) == 0 {
		panic("must provide at least 1 field")
	}

	if len(fields) == 1 {
		// templating fields
		fields = strings.Fields(fields[0])
	}

	selectionSet := &SelectionSet{
		fields: nil,
	}

	for _, fieldName := range fields {
		selectionSet.fields = append(selectionSet.fields, fieldName)
	}

	builder.selectionSet = selectionSet

	return builder
}

func (builder *QueryBuilder) Filter(function *QueryFunction) *QueryBuilder {
	filter := &Filter{
		connection: andConnector,
		function:   function,
	}

	builder.filters = append(builder.filters, filter)

	return builder
}

func (builder *QueryBuilder) OrFilter(function *QueryFunction) *QueryBuilder {
	filter := &Filter{
		connection: orConnector,
		function:   function,
	}

	builder.filters = append(builder.filters, filter)

	return builder
}

func (builder *QueryBuilder) FilterGroup(fn func(*FilterGroup)) *QueryBuilder {
	group := &FilterGroup{
		connection: andConnector,
	}

	fn(group)

	builder.filterGroups = append(builder.filterGroups, group)

	return builder
}

type SelectionSet struct {
	depth  int
	fields []string
	parent *SelectionSet
	edges  []*QueryBuilder
}

func (selection *SelectionSet) SetDepth(depth int) {
	selection.depth = depth
}

func (selection *SelectionSet) toDQL() (query string, args []interface{}, err error) {
	writer := NewWriter()

	if selection.depth > 0 {
		writer.indent = selection.depth
	}

	for _, fieldName := range selection.fields {
		writer.AddIndentedLine(fieldName)
	}

	for _, edge := range selection.edges {
		nestedEdge, edgesArgs, err := edge.toDQL()
		args = append(args, edgesArgs...)

		if err != nil {
			return "", nil, err
		}
		writer.AddIndentedLine(nestedEdge)
	}

	return writer.ToString(), args,nil
}
