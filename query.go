package dqlx

import (
	"context"
	"strings"

	"github.com/dgraph-io/dgo/v210"
)

// QueryBuilder represents the public API for building
// a secure dynamic Dgraph query
type QueryBuilder struct {
	rootEdge      edge
	variables     []QueryBuilder
	childrenEdges map[string][]QueryBuilder
	unmarshalInto interface{}

	client *dgo.Dgraph
}

// Query initialises the query builder
// with the provided root filter
// example: dqlx.Query(dqlx.EqFn(..,..))
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

	builder.rootEdge.Node = node{
		ParentName: builder.rootEdge.Name,
		Edges:      builder.childrenEdges,
	}
	return builder
}

// QueryType alias to initialise a query with
// the root function type()
// Example: dqlx.QueryType("User")
// Equivalent of: dqlx.Query(dqlx.TypeFn("User"))
func QueryType(typeName string) QueryBuilder {
	return Query(TypeFn(typeName))
}

// QueryEdge initialise a query builder with a specific name for
// the edge. Useful for reusable edges
// Example: dqlx.QueryEdge("name", dqlx.EqFn(..,..))
func QueryEdge(edgeName string, rootQueryFn *FilterFn) QueryBuilder {
	return Query(rootQueryFn).Name(edgeName)
}

// Variable initialise a variable query builder
// Example: dqlx.Variable(dqlx.EqFn(..,..))
func Variable(rootQueryFn *FilterFn) QueryBuilder {
	query := Query(rootQueryFn)
	query.rootEdge.IsVariable = true
	return query
}

// As sets an alias for the edge
// Example: dqlx.Query(...).As("C")
// DQL: { C AS rootQuery(func: ...) { ... } }
func (builder QueryBuilder) As(name string) QueryBuilder {
	builder.rootEdge.Alias = name
	return builder
}

// Name sets the name of the edge
// Exaple: dqlx.Query(...).Name("bladerunner")
// DQL: { bladerunner(func: ...) { ... }
func (builder QueryBuilder) Name(name string) QueryBuilder {
	builder.rootEdge.Name = name
	builder.rootEdge.Node.ParentName = name
	return builder
}

// ToDQL returns the current state of the query as DQL string
// Example: dqlx.Query(...).ToDQL()
func (builder QueryBuilder) ToDQL() (query string, args []interface{}, err error) {
	return QueriesToDQL(builder)
}

// Variable registers a variable within the query
// Example: dqlx.Query(...).Variable(variable)
func (builder QueryBuilder) Variable(queryBuilder QueryBuilder) QueryBuilder {
	builder.variables = append(builder.variables, queryBuilder)
	return builder
}

// Select assigns predicates to the selection set
// Example1: dqlx.Query(...).Select(`
// 	field1
//	field2
//	field3
//`)
//
// Example2: dqlx.Query(...).Select("field1", "field2", "field3")
func (builder QueryBuilder) Select(predicates ...interface{}) QueryBuilder {
	if len(predicates) == 0 {
		return builder
	}

	attributes := Select(predicates...).(nodeAttributes)

	selectedNode := node{
		ParentName:          builder.rootEdge.Name,
		HasParentAttributes: len(attributes.predicates) > 0,
		Edges:               builder.childrenEdges,
		Attributes:          attributes,
	}

	builder.rootEdge.Node = selectedNode

	return builder
}

// Fields alias of Select
// @Deprecated: use Select() instead
func (builder QueryBuilder) Fields(predicates ...interface{}) QueryBuilder {
	return builder.Select(predicates...)
}

// Facets requests facets for the current query
// Example1: dqlx.Query(...).Facets("field1")
// Example2: dqlx.Query(...).Facets(dqlx.Eq{"field1": "value"})
func (builder QueryBuilder) Facets(predicates ...interface{}) QueryBuilder {
	builder.rootEdge.Facets = append(builder.rootEdge.Facets, facetExpr{
		Predicates: predicates,
	})

	return builder
}

// Order requests an ordering for the result set
// Example1: dqlx.Query(...).Order(dqlx.OrderAsc("field1"))
// Example2: dqlx.Query(...).Order(dqlx.OrderDesc("field2"))
func (builder QueryBuilder) Order(order DQLizer) QueryBuilder {
	builder.rootEdge.Order = append(builder.rootEdge.Order, order)
	return builder
}

// OrderAsc alias for ordering in ascending order
// Example:    dqlx.Query(...).OrderAsc("field1")
// Equivalent: dqlx.Query(...).Order(dqlx.OrderAsc("field1"))
func (builder QueryBuilder) OrderAsc(predicate interface{}) QueryBuilder {
	builder.rootEdge.Order = append(builder.rootEdge.Order, orderBy{
		Direction: OrderDirectionAsc,
		Predicate: predicate,
	})
	return builder
}

// OrderDesc alias for ordering in descending order
// Example:    dqlx.Query(...).OrderDesc("field1")
// Equivalent: dqlx.Query(...).Order(dqlx.OrderDesc("field1"))
func (builder QueryBuilder) OrderDesc(predicate interface{}) QueryBuilder {
	builder.rootEdge.Order = append(builder.rootEdge.Order, orderBy{
		Direction: OrderDirectionDesc,
		Predicate: predicate,
	})
	return builder
}

// Filter requests filters for this query
// Example: dqlx.Query(...).Filter(dqlx.Eq{...}, dqlx.Gt{...})
func (builder QueryBuilder) Filter(filters ...DQLizer) QueryBuilder {
	for _, filter := range filters {
		builder.rootEdge.Filters = append(builder.rootEdge.Filters, filter)
	}
	return builder
}

// Paginate requests paginated results
// Example: dqlx.Query(...).Paginate(dqlx.Cursor{...})
func (builder QueryBuilder) Paginate(pagination Cursor) QueryBuilder {
	builder.rootEdge.Pagination = pagination
	return builder
}

// GroupBy requests group by
func (builder QueryBuilder) GroupBy(predicates ...string) QueryBuilder {
	for _, field := range predicates {
		builder.rootEdge.Group = append(builder.rootEdge.Group, GroupBy(field))
	}
	return builder
}

// Cascade adds cascade directive
func (builder QueryBuilder) Cascade(fields ...string) QueryBuilder {
	builder.rootEdge.Cascade = Cascade(fields...)

	return builder
}

// Edge adds an edge in the query selection
// Example1: dqlx.Query(...).Edge("path")
// Example2: dqlx.Query(...).Edge("parent->child->child")
// Example3: dqlx.Query(...).Edge("parent->child->child", dqlx.Select(""))
func (builder QueryBuilder) Edge(fullPath string, queryParts ...DQLizer) QueryBuilder {
	return builder.EdgeFn(fullPath, func(builder QueryBuilder) QueryBuilder {
		for _, part := range queryParts {
			switch cast := part.(type) {
			case filterExpr:
				builder = builder.Filter(part)
			case nodeAttributes:
				builder = builder.Select(cast.predicates...)
			case Cursor:
				builder = builder.Paginate(cast)
			case orderBy:
				builder = builder.Order(cast)
			case group:
				builder = builder.GroupBy(cast.Predicate)
			case facetExpr:
				builder = builder.Facets(cast.Predicates...)
			case cascadeExpr:
				builder = builder.Cascade(cast.fields...)
			case DQLizer:
				builder = builder.Filter(cast)
			}
		}
		return builder
	})
}

// EdgePath allows to defined an edge using the slice syntax for the path
// Example:    dqlx.Query(...).Edge([]string{"parent", "child", "child")
// Equivalent: dqlx.Query(...).Edge("parent->child->child")
func (builder QueryBuilder) EdgePath(fullPath []string, queryParts ...DQLizer) QueryBuilder {
	return builder.Edge(EdgePath(fullPath...), queryParts...)
}

// EdgeAs adds a new aliased edge
// Example: dqlx.Query(...).EdgeAs("C", "path", ...)
func (builder QueryBuilder) EdgeAs(as string, fullPath string, queryParts ...DQLizer) QueryBuilder {
	return builder.Edge(fullPath, queryParts...).As(as)
}

// EdgeFn allows to build an edge with a callback and query methods
// Example: dqlx.Query(...).EdgeFn("path", func(builder QueryBuilder) {
//  return builder.Select(...).Filter(...)
//})
func (builder QueryBuilder) EdgeFn(fullPath string, fn func(builder QueryBuilder) QueryBuilder) QueryBuilder {
	return builder.addEdgeFn(QueryEdge(fullPath, nil), fn)
}

// EdgeFromQuery allows to add an external constructed edge to the query
// Example: dqlx.Query(...).EdgeFromQuery(dqlx.Query(...))
func (builder QueryBuilder) EdgeFromQuery(edge QueryBuilder) QueryBuilder {
	return builder.addEdgeFn(edge, nil)
}

// UnmarshalInto requests to unmarshal the result set into this specific interface{}
// Example: dqlx.Query(...).UnmarshalInto(&value)+
func (builder QueryBuilder) UnmarshalInto(value interface{}) QueryBuilder {
	builder.unmarshalInto = value
	return builder
}

// WithDClient allows to swap the underline dgo.Dgraph client
// Example: dqlx.Query(...).WithDClient(dgoClient)
func (builder QueryBuilder) WithDClient(client *dgo.Dgraph) QueryBuilder {
	builder.client = client
	return builder
}

// Execute executes the current state of the query sending the operation
// to dgraph
// Example: dqlx.Query(...).Execute(ctx, ...)
func (builder QueryBuilder) Execute(ctx context.Context, options ...OperationExecutorOptionFn) (*Response, error) {
	executor := NewDGoExecutor(builder.client)

	for _, option := range options {
		option(executor)
	}
	return executor.ExecuteQueries(ctx, builder)
}

// GetName returns the name of the query edge
func (builder QueryBuilder) GetName() string {
	return builder.rootEdge.Name
}

func (builder QueryBuilder) addEdgeFn(query QueryBuilder, fn func(builder QueryBuilder) QueryBuilder) QueryBuilder {
	edgePathParts := ParseEdge(query.rootEdge.Name)

	if len(edgePathParts) == 0 {
		return builder
	}

	edgeBuilder := query
	edgeBuilder.rootEdge.IsRoot = false
	edgeBuilder.rootEdge.Node.Edges = builder.childrenEdges
	edgeBuilder.rootEdge.Node.HasParentAttributes = builder.rootEdge.Node.Attributes != nil
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

// IsEmptyQuery determine if a given query is an empty generated query
func IsEmptyQuery(query string) bool {
	return "query () {  {  } }" == query
}
