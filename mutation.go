package dqlx

import (
	"bytes"
	"context"
	"github.com/dgraph-io/dgo/v200"
	"strings"
)

type MutationBuilder struct {
	condition DQLizer
	query     QueryBuilder

	setData       interface{}
	delData       interface{}
	unmarshalInto interface{}

	client *dgo.Dgraph
}

func Mutation() MutationBuilder {
	return MutationBuilder{}
}

func (builder MutationBuilder) Condition(condition DQLizer) MutationBuilder {
	if conditionType, ok := condition.(mutationCondition); ok {
		builder.condition = conditionType
	} else {
		builder.condition = Condition(condition)
	}

	return builder
}

func (builder MutationBuilder) Query(condition DQLizer) MutationBuilder {
	builder.condition = condition
	return builder
}

func (builder MutationBuilder) UnmarshalInto(value interface{}) MutationBuilder {
	builder.unmarshalInto = value
	return builder
}

func (builder MutationBuilder) Set(data interface{}) MutationBuilder {
	builder.setData = data
	return builder
}

func (builder MutationBuilder) Delete(data interface{}) MutationBuilder {
	builder.delData = data
	return builder
}

func (builder MutationBuilder) Execute(ctx context.Context, options ...ExecutorOptionFn) (*Response, error) {
	executor := NewDGoExecutor(builder.client)

	for _, option := range options {
		option(executor)
	}

	defer func() {
		builder.unmarshalInto = nil
	}()

	return executor.ExecuteMutations(ctx, builder)
}

func (builder MutationBuilder) WithDClient(client *dgo.Dgraph) MutationBuilder {
	builder.client = client
	return builder
}

type mutationCondition struct {
	Filters []DQLizer
}

func Condition(filters ...DQLizer) DQLizer {
	return mutationCondition{Filters: filters}
}

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
