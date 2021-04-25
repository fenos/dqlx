package deku

import (
	"strings"
)

type QueryBuilder struct {
	edge      edge
	selection selectionSet
	edges     map[string][]QueryBuilder
	filters   []DQLizer
}

func Query(name string, rootQueryFn FilterFn) QueryBuilder {
	builder := QueryBuilder{
		edge: edge{
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
	grammar := queryGrammar{
		RootEdge:  builder.edge,
		Filters:   builder.filters,
		Selection: builder.selection,
	}

	return grammar.ToQuery(builder.edge.Name)
}

func (builder QueryBuilder) toDql() (query string, args []interface{}, err error) {
	grammar := queryGrammar{
		RootEdge:  builder.edge,
		Filters:   builder.filters,
		Selection: builder.selection,
	}

	return grammar.ToDQL()
}

//func (builder *QueryBuilder) MergeEdge(edgeQuery QueryBuilder) *QueryBuilder {
//	edgeQuery.isRoot = false
//	edgeQuery.depth = builder.depth + 1
//
//	//builder.selection.Edges = append(builder.selection.Edges, edgeQuery)
//
//	return builder
//}

func (builder QueryBuilder) Fields(fields ...string) QueryBuilder {
	if len(fields) == 1 {
		// templating Fields
		fields = strings.Fields(fields[0])
	}

	selectionSet := selectionSet{
		Parent: builder.edge,
		Edges:  builder.edges,
	}

	for _, fieldName := range fields {
		selectionSet.Fields = append(selectionSet.Fields, fieldName)
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
	builder.edge.Pagination = pagination
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
		edge: edge{
			Name:   fullPath,
			IsRoot: false,
		},
		edges: builder.edges,
	}

	parentPath := fullPath

	if len(edgePathParts) == 1 {
		parentPath = builder.edge.Name
	} else {
		parents := edgePathParts[0 : len(edgePathParts)-1]
		parentPath = strings.Join(parents, "->")
	}

	edgeBuilder = fn(edgeBuilder)

	builder.edges[parentPath] = append(builder.edges[parentPath], edgeBuilder)

	return builder
}
