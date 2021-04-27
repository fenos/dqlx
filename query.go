package deku

import (
	"strings"
)

type QueryBuilder struct {
	rootEdge      edge
	variables     []QueryBuilder
	childrenEdges map[string][]QueryBuilder
}

func Query(name string, rootQueryFn *FilterFn) QueryBuilder {
	var rootFilter DQLizer

	if rootQueryFn != nil {
		rootFilter = *rootQueryFn
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

func Variable(rootQueryFn *FilterFn) QueryBuilder {
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

func (builder QueryBuilder) Facets(predicates ...interface{}) QueryBuilder {
	builder.rootEdge.Facets = append(builder.rootEdge.Facets, facetExpr{
		Predicates: predicates,
	})

	return builder
}

func (builder QueryBuilder) Order(order DQLizer) QueryBuilder {
	builder.rootEdge.Order = append(builder.rootEdge.Order, order)
	return builder
}

func (builder QueryBuilder) OrderAsc(predicate interface{}) QueryBuilder {
	builder.rootEdge.Order = append(builder.rootEdge.Order, orderBy{
		Direction: OrderDirectionAsc,
		Predicate: predicate,
	})
	return builder
}

func (builder QueryBuilder) OrderDesc(predicate interface{}) QueryBuilder {
	builder.rootEdge.Order = append(builder.rootEdge.Order, orderBy{
		Direction: OrderDirectionDesc,
		Predicate: predicate,
	})
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

func (builder QueryBuilder) GroupBy(fields ...string) QueryBuilder {
	for _, field := range fields {
		builder.rootEdge.Group = append(builder.rootEdge.Group, GroupBy(field))
	}
	return builder
}

func (builder QueryBuilder) Edge(fullPath string, queryParts ...DQLizer) QueryBuilder {
	return builder.EdgeFn(fullPath, func(builder QueryBuilder) QueryBuilder {
		for _, part := range queryParts {
			switch cast := part.(type) {
			case filterExpr:
				builder = builder.Filter(part)
			case Fields:
				builder = builder.Fields(string(cast))
			case Pagination:
				builder = builder.Paginate(cast)
			case orderBy:
				builder = builder.Order(cast)
			case group:
				builder = builder.GroupBy(cast.Predicate)
			case facetExpr:
				builder = builder.Facets(cast.Predicates...)
			case DQLizer:
				builder = builder.Filter(cast)
			}
		}
		return builder
	})
}

func (builder QueryBuilder) EdgeFn(fullPath string, fn func(builder QueryBuilder) QueryBuilder) QueryBuilder {
	builder.addEdgeFn(Query(fullPath, nil), fn)
	return builder
}

func (builder QueryBuilder) EdgeFromQuery(edge QueryBuilder) QueryBuilder {
	builder.addEdgeFn(edge, nil)
	return builder
}

func (builder QueryBuilder) addEdgeFn(query QueryBuilder, fn func(builder QueryBuilder) QueryBuilder) QueryBuilder {
	edgePathParts := ParseEdge(query.rootEdge.Name)

	if len(edgePathParts) == 0 {
		return builder
	}

	edgeBuilder := query
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

	if fn != nil {
		edgeBuilder = fn(edgeBuilder)
	}

	builder.childrenEdges[parentPath] = append(builder.childrenEdges[parentPath], edgeBuilder)

	return builder
}
