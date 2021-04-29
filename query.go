package deku

import (
	"strings"
)

type QueryBuilder struct {
	rootEdge      edge
	variables     []QueryBuilder
	childrenEdges map[string][]QueryBuilder
}

func Query(rootQueryFn *FilterFn) QueryBuilder {
	var rootFilter DQLizer

	if rootQueryFn != nil {
		rootFilter = *rootQueryFn
	}

	builder := QueryBuilder{
		rootEdge: edge{
			Name:       "rootQuery",
			RootFilter: rootFilter,
			Filters:    []DQLizer{},
			IsRoot:     true,
			IsVariable: false,
		},
		childrenEdges: map[string][]QueryBuilder{},
	}

	builder.rootEdge.Selection = selectionSet{
		ParentName: builder.rootEdge.Name,
		Edges:      builder.childrenEdges,
	}
	return builder
}

func QueryType(typeName string) QueryBuilder {
	return Query(TypeFn(typeName))
}

func QueryEdge(edgeName string, rootQueryFn *FilterFn) QueryBuilder {
	return Query(rootQueryFn).Name(edgeName)
}

func Variable(rootQueryFn *FilterFn) QueryBuilder {
	query := Query(rootQueryFn)
	query.rootEdge.IsVariable = true
	return query
}

func (builder QueryBuilder) As(name string) QueryBuilder {
	builder.rootEdge.Alias = name
	return builder
}

func (builder QueryBuilder) Name(name string) QueryBuilder {
	builder.rootEdge.Name = name
	builder.rootEdge.Selection.ParentName = name
	return builder
}

func (builder QueryBuilder) ToDQL() (query string, args map[string]string, err error) {
	return OperationQuery(builder)
}

func (builder QueryBuilder) Variable(queryBuilder QueryBuilder) QueryBuilder {
	builder.variables = append(builder.variables, queryBuilder)
	return builder
}

func (builder QueryBuilder) Fields(fieldNames ...string) QueryBuilder {
	if len(fieldNames) == 0 {
		return builder
	}

	allFields := Fields(fieldNames...).(fields)

	selection := selectionSet{
		ParentName:      builder.rootEdge.Name,
		HasParentFields: len(allFields.predicates) > 0,
		Edges:           builder.childrenEdges,
		Fields:          allFields.predicates,
	}

	builder.rootEdge.Selection = selection

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
			case fields:
				builder = builder.Fields(cast.predicates...)
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
	return builder.addEdgeFn(QueryEdge(fullPath, nil), fn)
}

func (builder QueryBuilder) EdgeFromQuery(edge QueryBuilder) QueryBuilder {
	return builder.addEdgeFn(edge, nil)
}

func (builder QueryBuilder) addEdgeFn(query QueryBuilder, fn func(builder QueryBuilder) QueryBuilder) QueryBuilder {
	edgePathParts := ParseEdge(query.rootEdge.Name)

	if len(edgePathParts) == 0 {
		return builder
	}

	edgeBuilder := query
	edgeBuilder.rootEdge.IsRoot = false
	edgeBuilder.rootEdge.Selection.Edges = builder.childrenEdges
	edgeBuilder.rootEdge.Selection.HasParentFields = len(builder.rootEdge.Selection.Fields) > 0
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
