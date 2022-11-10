package dqlx

import (
	"context"

	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
)

// SchemaExecutor executes schema operations
type SchemaExecutor struct {
	client *dgo.Dgraph

	dropAll         bool
	runInBackground bool
}

// SchemaExecutorOptionFn represents an function modifier
// for the SchemaExecutor options
type SchemaExecutorOptionFn func(*SchemaExecutor)

func WithDropAllSchema(dropAll bool) SchemaExecutorOptionFn {
	return func(schema *SchemaExecutor) {
		schema.dropAll = dropAll
	}
}

// WithRunInBackground instructs Dgraph to run indexes in the background
func WithRunInBackground(runInBackground bool) SchemaExecutorOptionFn {
	return func(schema *SchemaExecutor) {
		schema.runInBackground = runInBackground
	}
}

// NewSchemaExecutor creates a new schema executor
func NewSchemaExecutor(client *dgo.Dgraph) *SchemaExecutor {
	return &SchemaExecutor{
		client:          client,
		runInBackground: true,
	}
}

// AlterSchema alters the schema with new predicates or types
// No drop operation would occur if not specifying (DropAll)
func (executor SchemaExecutor) AlterSchema(ctx context.Context, schema *SchemaBuilder, options ...SchemaExecutorOptionFn) error {

	schemaDefinition, err := schema.ToDQL()

	if err != nil {
		return err
	}

	return executor.client.Alter(ctx, &api.Operation{
		Schema:          schemaDefinition,
		DropAll:         executor.dropAll,
		RunInBackground: executor.runInBackground,
	})
}

// DropType drops a type
func (executor SchemaExecutor) DropType(ctx context.Context, typeName string) error {
	return executor.client.Alter(ctx, &api.Operation{
		DropOp:          api.Operation_TYPE,
		DropValue:       typeName,
		RunInBackground: executor.runInBackground,
	})
}

// DropPredicate drops a predicate
func (executor SchemaExecutor) DropPredicate(ctx context.Context, predicateName string) error {
	return executor.client.Alter(ctx, &api.Operation{
		DropOp:          api.Operation_ATTR,
		DropValue:       predicateName,
		RunInBackground: executor.runInBackground,
	})
}
