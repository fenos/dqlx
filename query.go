package deku

import (
	"strings"
)

type QueryBuilder struct {
	rootEdge  edge
	selection selectionSet
	edges     map[string][]QueryBuilder
	filters   []DQLizer
	variables []QueryBuilder
}

func Query(name string, rootQueryFn FilterFn) QueryBuilder {
	builder := QueryBuilder{
		rootEdge: edge{
			Name:   name,
			IsRoot: true,
			Filters: []DQLizer{
				rootQueryFn(),
			},
		},
		edges: map[string][]QueryBuilder{},
	}

	return builder
}

func (builder QueryBuilder) ToDQL() (query string, args map[string]interface{}, err error) {
	grammar := builder.toGrammar()

	operation := rootOperation{
		operations: []queryGrammar{grammar},
	}
	return operation.ToQuery()
}

func (builder QueryBuilder) toDQL() (query string, args []interface{}, err error) {
	return builder.toGrammar().ToDQL()
}

func (builder QueryBuilder) toGrammar() queryGrammar {
	return queryGrammar{
		RootEdge:  builder.rootEdge,
		Filters:   builder.filters,
		Selection: builder.selection,
		Variables: builder.variables,
	}
}

func (builder QueryBuilder) Variable(queryBuilder QueryBuilder) QueryBuilder {
	builder.variables = append(builder.variables, queryBuilder)
	return builder
}

func (builder QueryBuilder) Fields(fields ...string) QueryBuilder {
	if len(fields) == 1 {
		// templating Fields
		fields = strings.Fields(fields[0])
	}

	selectionSet := selectionSet{
		Parent: builder.rootEdge,
		Edges:  builder.edges,
		Fields: fields,
	}

	builder.selection = selectionSet

	return builder
}

func (builder QueryBuilder) Filter(filters ...DQLizer) QueryBuilder {
	for _, filter := range filters {
		builder.filters = append(builder.filters, filter)
	}
	return builder
}

func (builder QueryBuilder) Paginate(pagination Pagination) QueryBuilder {
	builder.rootEdge.Pagination = pagination
	return builder
}

func (builder QueryBuilder) Edge(fullPath string, fields []string, filters ...DQLizer) QueryBuilder {
	return builder.EdgeFn(fullPath, func(builder QueryBuilder) QueryBuilder {
		return builder.Fields(fields...).Filter(filters...)
	})
}

func (builder QueryBuilder) EdgeFn(fullPath string, fn func(builder QueryBuilder) QueryBuilder) QueryBuilder {
	edgePathParts := strings.Split(fullPath, "->")

	if len(edgePathParts) == 0 {
		return builder
	}

	edgeBuilder := QueryBuilder{
		rootEdge: edge{
			Name:   fullPath,
			IsRoot: false,
		},
		edges: builder.edges,
	}

	parentPath := fullPath

	if len(edgePathParts) == 1 {
		parentPath = builder.rootEdge.Name
	} else {
		parents := edgePathParts[0 : len(edgePathParts)-1]
		parentPath = strings.Join(parents, "->")
	}

	edgeBuilder = fn(edgeBuilder)

	builder.edges[parentPath] = append(builder.edges[parentPath], edgeBuilder)

	return builder
}
