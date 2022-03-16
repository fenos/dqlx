package dqlx

import (
	"context"

	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"google.golang.org/grpc"
)

var symbolValuePlaceholder = "??"
var symbolEdgeTraversal = "->"

type Args []interface{}

// DQLizer implementors are able to define a custom dql statement
type DQLizer interface {
	ToDQL() (query string, args Args, err error)
}

// Executor implementors are able to define a custom way
// of executing queries and mutations
type Executor interface {
	ExecuteQueries(ctx context.Context, queries ...QueryBuilder) (*Response, error)
	ExecuteMutations(ctx context.Context, mutations ...MutationBuilder) (*Response, error)
}

// DB represents the public API for interacting with a DGraph
// Database
type DB interface {
	Query(rootFn *FilterFn) QueryBuilder
	QueryType(typeName string) QueryBuilder
	QueryEdge(edgeName string, rootQueryFn *FilterFn) QueryBuilder

	Mutation() MutationBuilder

	ExecuteQueries(ctx context.Context, queries []QueryBuilder, options ...OperationExecutorOptionFn) (*Response, error)
	ExecuteMutations(ctx context.Context, mutations []MutationBuilder, options ...OperationExecutorOptionFn) (*Response, error)

	Schema() *SchemaBuilder
	NewTxn() *dgo.Txn
	NewReadOnlyTxn() *dgo.Txn
	GetDgraph() *dgo.Dgraph
}

type dqlx struct {
	dgraph   *dgo.Dgraph
	executor Executor
}

// Connect connects to a DGraph Cluster.
// upon success a DB instance is returned ready to be used
// to interact with DGraph
func Connect(addresses ...string) (DB, error) {
	clients := make([]api.DgraphClient, len(addresses))

	for index, address := range addresses {
		dial, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
		clients[index] = api.NewDgraphClient(dial)
	}

	dgraph := dgo.NewDgraphClient(clients...)

	return FromClient(dgraph), nil
}

// FromClient creates a DB instance from a raw dgraph client
func FromClient(dgraph *dgo.Dgraph) DB {
	return &dqlx{
		dgraph: dgraph,
	}
}

// Query returns a QueryBuilder
func (dqlx *dqlx) Query(rootFn *FilterFn) QueryBuilder {
	return Query(rootFn).WithDClient(dqlx.dgraph)
}

// QueryType returns a QueryBuilder with a default type() filter
func (dqlx *dqlx) QueryType(typeName string) QueryBuilder {
	return QueryType(typeName).WithDClient(dqlx.dgraph)
}

// QueryEdge returns a QueryBuilder with the ability of providing the query name
func (dqlx *dqlx) QueryEdge(edgeName string, rootQueryFn *FilterFn) QueryBuilder {
	return QueryEdge(edgeName, rootQueryFn).WithDClient(dqlx.dgraph)
}

// Mutation returns a MutationBuilder
func (dqlx *dqlx) Mutation() MutationBuilder {
	return Mutation().WithDClient(dqlx.dgraph)
}

// ExecuteQueries executes multiple queries joining them into 1 request
func (dqlx *dqlx) ExecuteQueries(ctx context.Context, queries []QueryBuilder, options ...OperationExecutorOptionFn) (*Response, error) {
	executor := NewDGoExecutor(dqlx.dgraph)
	for _, option := range options {
		option(executor)
	}
	return executor.ExecuteQueries(ctx, queries...)
}

// ExecuteMutations executes multiple mutations
func (dqlx *dqlx) ExecuteMutations(ctx context.Context, mutations []MutationBuilder, options ...OperationExecutorOptionFn) (*Response, error) {
	executor := NewDGoExecutor(dqlx.dgraph)
	for _, option := range options {
		option(executor)
	}
	return executor.ExecuteMutations(ctx, mutations...)
}

// Schema returns a schema builder
func (dqlx *dqlx) Schema() *SchemaBuilder {
	return NewSchema().WithClient(dqlx.dgraph)
}

// NewTxn creates a new transaction
func (dqlx *dqlx) NewTxn() *dgo.Txn {
	return dqlx.dgraph.NewTxn()
}

// NewReadOnlyTxn creates a new read-only transaction
func (dqlx *dqlx) NewReadOnlyTxn() *dgo.Txn {
	return dqlx.NewReadOnlyTxn()
}

// GetDgraph returns the dgraph client
func (dqlx *dqlx) GetDgraph() *dgo.Dgraph {
	return dqlx.dgraph
}
