package dqlx

import (
	"context"
	"github.com/dgraph-io/dgo/v200"
	"github.com/dgraph-io/dgo/v200/protos/api"
	"google.golang.org/grpc"
)

var symbolValuePlaceholder = "??"
var symbolEdgeTraversal = "->"

type DQLizer interface {
	ToDQL() (query string, args []interface{}, err error)
}

type Executor interface {
	ExecuteQueries(ctx context.Context, queries ...QueryBuilder) (*QueryResponse, error)
}

type Dqlx interface {
	Query(rootFn *FilterFn) QueryBuilder
	QueryType(typeName string) QueryBuilder
	QueryEdge(edgeName string, rootQueryFn *FilterFn) QueryBuilder
	Variable(rootQueryFn *FilterFn) QueryBuilder

	ExecuteQueries(ctx context.Context, queries []QueryBuilder, options ...ExecutorOptionFn) (*QueryResponse, error)
	NewTxn() *dgo.Txn
	NewReadOnlyTxn() *dgo.Txn
}

type dqlx struct {
	dgraph   *dgo.Dgraph
	executor Executor
}

func Connect(address string) (Dqlx, error) {
	dial, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	dgraph := dgo.NewDgraphClient(
		api.NewDgraphClient(dial),
	)

	return FromClient(dgraph), nil
}

func FromClient(dgraph *dgo.Dgraph) Dqlx {
	return &dqlx{
		dgraph: dgraph,
	}
}

func (dqlx *dqlx) Query(rootFn *FilterFn) QueryBuilder {
	return Query(rootFn).WithDClient(dqlx.dgraph)
}

func (dqlx *dqlx) QueryType(typeName string) QueryBuilder {
	return QueryType(typeName).WithDClient(dqlx.dgraph)
}

func (dqlx *dqlx) QueryEdge(edgeName string, rootQueryFn *FilterFn) QueryBuilder {
	return QueryEdge(edgeName, rootQueryFn).WithDClient(dqlx.dgraph)
}

func (dqlx *dqlx) Variable(rootQueryFn *FilterFn) QueryBuilder {
	return Variable(rootQueryFn).WithDClient(dqlx.dgraph)
}

func (dqlx *dqlx) ExecuteQueries(ctx context.Context, queries []QueryBuilder, options ...ExecutorOptionFn) (*QueryResponse, error) {
	executor := NewDGoExecutor(dqlx.dgraph)
	for _, option := range options {
		option(executor)
	}
	return executor.ExecuteQueries(ctx, queries...)
}

func (dqlx *dqlx) NewTxn() *dgo.Txn {
	return dqlx.dgraph.NewTxn()
}

func (dqlx *dqlx) NewReadOnlyTxn() *dgo.Txn {
	return dqlx.NewReadOnlyTxn()
}
