package dqlx

import (
	"bytes"
	"context"
	"github.com/dgraph-io/dgo/v200"
	"strings"
)

// MutationBuilder used to construct mutations
// in a fluent way.
//
// Example: Mutation().Set(data).Execute(ctx)
// Example: Mutation().Delete(data).Execute(ctx)
type MutationBuilder struct {
	condition DQLizer
	query     QueryBuilder

	setData       interface{}
	delData       interface{}
	unmarshalInto interface{}

	client *dgo.Dgraph
}

// Mutation creates a new MutationBuilder
func Mutation() MutationBuilder {
	return MutationBuilder{}
}

// Condition assigns a condition for this mutation
// Used for conditional Upserts
// DGraphDoc: https://dgraph.io/docs/mutations/conditional-upsert/
func (builder MutationBuilder) Condition(condition DQLizer) MutationBuilder {
	if conditionType, ok := condition.(mutationCondition); ok {
		builder.condition = conditionType
	} else {
		builder.condition = Condition(condition)
	}

	return builder
}

// Query assigns a query block for this mutation
// making it an Upsert
func (builder MutationBuilder) Query(query QueryBuilder) MutationBuilder {
	builder.query = query
	return builder
}

// UnmarshalInto defines where the data should be marshalled into,
// once the mutation gets executed
func (builder MutationBuilder) UnmarshalInto(value interface{}) MutationBuilder {
	builder.unmarshalInto = value
	return builder
}

// Set sets some data to be inserted or updated
func (builder MutationBuilder) Set(data interface{}) MutationBuilder {
	builder.setData = data
	return builder
}

// Delete sets some data to be deleted
func (builder MutationBuilder) Delete(data interface{}) MutationBuilder {
	builder.delData = data
	return builder
}

// Execute executes the mutation
func (builder MutationBuilder) Execute(ctx context.Context, options ...OperationExecutorOptionFn) (*Response, error) {
	executor := NewDGoExecutor(builder.client)

	for _, option := range options {
		option(executor)
	}
	return executor.ExecuteMutations(ctx, builder)
}

// WithDClient changes Dgraph client for this mutation
func (builder MutationBuilder) WithDClient(client *dgo.Dgraph) MutationBuilder {
	builder.client = client
	return builder
}

type mutationCondition struct {
	Filters []DQLizer
}

// Condition returns a condition statement
func Condition(filters ...DQLizer) DQLizer {
	return mutationCondition{Filters: filters}
}

// ToDQL returns a DQL statement for a mutation condition
func (condition mutationCondition) ToDQL() (query string, args []interface{}, err error) {
	writer := bytes.Buffer{}
	writer.WriteString(" @if(")

	var statements []string
	for _, filter := range condition.Filters {
		filterDql, filterArgs, err := filter.ToDQL()

		if err != nil {
			return "", nil, err
		}

		filterDql, _ = replacePlaceholders(filterDql, filterArgs, func(index int, value interface{}) string {
			return toVariableValue(value)
		})

		statements = append(statements, filterDql)
	}

	writer.WriteString(strings.Join(statements, " "))
	writer.WriteString(") ")

	return writer.String(), nil, err
}
