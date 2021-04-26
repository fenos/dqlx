package deku

import (
	"strings"
)

type QueryBuilder struct {
	rootEdge      edge
	childrenEdges map[string][]QueryBuilder
	variables     []QueryBuilder
}

func Query(name string, rootQueryFn FilterFn) QueryBuilder {
	var rootFilter DQLizer

	if rootQueryFn != nil {
		rootFilter = rootQueryFn()
	}

	builder := QueryBuilder{
		rootEdge: edge{
			Name:       name,
			RootFilter: rootFilter,
			Filters:    []DQLizer{},
			IsRoot:     true,
			IsVariable: false,
		},
		childrenEdges: map[string][]QueryBuilder{},
	}

	builder.rootEdge.Selection = selectionSet{
		Parent: &builder.rootEdge,
		Edges:  builder.childrenEdges,
	}
	return builder
}

func Variable(rootQueryFn FilterFn) QueryBuilder {
	query := Query("", rootQueryFn)
	query.rootEdge.IsVariable = true
	return query
}

func (builder QueryBuilder) As(name string) QueryBuilder {
	builder.rootEdge.Name = name
	return builder
}

func (builder QueryBuilder) ToDQL() (query string, args map[string]interface{}, err error) {
	return OperationQuery(builder)
}

func (builder QueryBuilder) Variable(queryBuilder QueryBuilder) QueryBuilder {
	builder.variables = append(builder.variables, queryBuilder)
	return builder
}

func (builder QueryBuilder) Fields(fields ...string) QueryBuilder {
	if len(fields) == 0 {
		return builder
	}

	if len(fields) == 1 {
		// templating ParseFields
		fields = strings.Fields(fields[0])
	}

	selectionSet := selectionSet{
		Parent: &builder.rootEdge,
		Edges:  builder.childrenEdges,
		Fields: fields,
	}

	builder.rootEdge.Selection = selectionSet

	return builder
}

func (builder QueryBuilder) Filter(filters ...DQLizer) QueryBuilder {
	for _, filter := range filters {
		builder.rootEdge.Filters = append(builder.rootEdge.Filters, filter)
	}
	return builder
}

func (builder QueryBuilder) Paginate(pagination Pagination) QueryBuilder {
	builder.rootEdge.Pagination = pagination
	return builder
}

func (builder QueryBuilder) Edge(fullPath string, fields string, filters ...DQLizer) QueryBuilder {
	return builder.EdgeFn(fullPath, func(builder QueryBuilder) QueryBuilder {
		return builder.Fields(ParseFields(fields)...).Filter(filters...)
	})
}

func (builder QueryBuilder) EdgeFn(fullPath string, fn func(builder QueryBuilder) QueryBuilder) QueryBuilder {
	edgePathParts := ParseEdge(fullPath)

	if len(edgePathParts) == 0 {
		return builder
	}

	edgeBuilder := Query(fullPath, nil)
	edgeBuilder.rootEdge.IsRoot = false

	edgeBuilder.rootEdge.Selection.Edges = builder.childrenEdges
	edgeBuilder.childrenEdges = builder.childrenEdges

	var parentPath string

	if len(edgePathParts) == 1 {
		parentPath = builder.rootEdge.Name
	} else {
		parents := edgePathParts[0 : len(edgePathParts)-1]
		parentPath = strings.Join(parents, symbolEdgeTraversal)
	}

	edgeBuilder = fn(edgeBuilder)

	builder.childrenEdges[parentPath] = append(builder.childrenEdges[parentPath], edgeBuilder)

	return builder
}
