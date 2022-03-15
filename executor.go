package dqlx

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"time"

	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"github.com/mitchellh/mapstructure"
)

// OperationExecutor represents a Dgraph executor for operations
// such as Queries and Mutations using the official dgo client
type OperationExecutor struct {
	client *dgo.Dgraph
	tnx    *dgo.Txn

	readOnly   bool
	bestEffort bool
}

// OperationExecutorOptionFn used to modify options of the executor
type OperationExecutorOptionFn func(executor *OperationExecutor)

// WithTnx configures a transaction to be used
// for the current execution
func WithTnx(tnx *dgo.Txn) OperationExecutorOptionFn {
	return func(executor *OperationExecutor) {
		executor.tnx = tnx
	}
}

// WithClient configures a client for the current execution
func WithClient(client *dgo.Dgraph) OperationExecutorOptionFn {
	return func(executor *OperationExecutor) {
		executor.client = client
	}
}

// WithReadOnly marks the execution as a read-only operation
// you can use this only on queries
func WithReadOnly(readOnly bool) OperationExecutorOptionFn {
	return func(executor *OperationExecutor) {
		executor.readOnly = readOnly
	}
}

// WithBestEffort sets the best effort flag for the current execution
func WithBestEffort(bestEffort bool) OperationExecutorOptionFn {
	return func(executor *OperationExecutor) {
		executor.bestEffort = bestEffort
	}
}

// NewDGoExecutor creates a new OperationExecutor
func NewDGoExecutor(client *dgo.Dgraph) *OperationExecutor {
	return &OperationExecutor{
		client: client,
	}
}

// ExecuteQueries executes a query operation. If multiple queries are provided they will
// get merged into a single one.
// the transaction will be automatically committed if a custom tnx is not provided.
// only non-readonly transactions will be committed.
func (executor OperationExecutor) ExecuteQueries(ctx context.Context, queries ...QueryBuilder) (*Response, error) {
	if err := executor.ensureClient(); err != nil {
		return nil, err
	}

	query, variables, err := QueriesToDQL(queries...)
	if err != nil {
		return nil, err
	}

	tx := executor.getTnx()

	defer tx.Discard(ctx)

	resp, err := tx.QueryWithVars(ctx, query, variables)
	if err != nil {
		return nil, err
	}

	if !executor.readOnly {
		if executor.tnx != nil {
			err := tx.Commit(ctx)
			if err != nil {
				return nil, err
			}
		}
	}

	return executor.toResponse(resp, queries...)
}

// ExecuteMutations executes one ore more mutations.
// the transaction will be automatically committed if a custom tnx is not provided.
func (executor OperationExecutor) ExecuteMutations(ctx context.Context, mutations ...MutationBuilder) (*Response, error) {
	if err := executor.ensureClient(); err != nil {
		return nil, err
	}

	var queries []QueryBuilder
	var mutationRequests []*api.Mutation

	for _, mutation := range mutations {
		var condition string

		if mutation.condition != nil {
			conditionDql, _, err := mutation.condition.ToDQL()
			if err != nil {
				return nil, err
			}
			condition = conditionDql
		}

		queries = append(queries, mutation.query)
		setData, deleteData, err := mutationData(mutation)

		if err != nil {
			return nil, err
		}

		mutationRequest := &api.Mutation{
			SetJson:    setData,
			DeleteJson: deleteData,
			Cond:       condition,
			CommitNow:  executor.tnx == nil,
		}

		mutationRequests = append(mutationRequests, mutationRequest)
	}

	query, variables, err := QueriesToDQL(queries...)

	if IsEmptyQuery(query) {
		query = ""
		variables = nil
	}

	request := &api.Request{
		Query:      query,
		Vars:       variables,
		ReadOnly:   executor.readOnly,
		BestEffort: executor.bestEffort,
		Mutations:  mutationRequests,
		CommitNow:  executor.tnx == nil,
		RespFormat: api.Request_JSON,
	}

	tx := executor.getTnx()
	defer tx.Discard(ctx)

	resp, err := tx.Do(ctx, request)

	if err != nil {
		return nil, err
	}

	return executor.toResponse(resp, queries...)
}

func (executor OperationExecutor) toResponse(resp *api.Response, queries ...QueryBuilder) (*Response, error) {
	var dataPathKey string

	if len(queries) == 1 {
		dataPathKey = queries[0].rootEdge.Name
	} else {
		dataPathKey = ""
	}

	queryResponse := &Response{
		dataKeyPath: dataPathKey,
		Raw:         resp,
	}

	queries = ensureUniqueQueryNames(queries)

	for _, queryBuilder := range queries {
		if queryBuilder.unmarshalInto == nil {
			continue
		}
		singleResponse := &Response{
			dataKeyPath: queryBuilder.rootEdge.Name,
			Raw:         resp,
		}

		err := singleResponse.Unmarshal(queryBuilder.unmarshalInto)

		if err != nil {
			return nil, err
		}
	}

	return queryResponse, nil
}

func mutationData(mutation MutationBuilder) (updateData []byte, deleteData []byte, err error) {
	var setDataBytes []byte
	var deleteDataBytes []byte

	if mutation.setData != nil {
		setBytes, err := json.Marshal(mutation.setData)
		if err != nil {
			return nil, nil, err
		}
		setDataBytes = setBytes
	}

	if mutation.delData != nil {
		deleteBytes, err := json.Marshal(mutation.delData)
		if err != nil {
			return nil, nil, err
		}
		deleteDataBytes = deleteBytes
	}

	return setDataBytes, deleteDataBytes, nil
}

func (executor OperationExecutor) ensureClient() error {
	if executor.client == nil {
		return errors.New("cannot execute query without setting a dqlx. use DClient() to set one")
	}
	return nil
}

func (executor OperationExecutor) getTnx() *dgo.Txn {
	tx := executor.tnx

	if tx == nil {
		if executor.readOnly {
			tx = executor.client.NewReadOnlyTxn()
		} else {
			tx = executor.client.NewTxn()
		}
	}
	return tx
}

// Response represents an operation response
type Response struct {
	Raw         *api.Response
	dataKeyPath string
}

// Unmarshal allows to dynamically marshal the result set of a query
// into an interface value
func (response Response) Unmarshal(value interface{}) error {
	values := map[string]interface{}{}
	err := json.Unmarshal(response.Raw.Json, &values)

	if err != nil {
		return err
	}

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: func(from reflect.Value, to reflect.Value) (interface{}, error) {
			if _, ok := to.Interface().(time.Time); ok {
				return time.Parse(time.RFC3339, from.String())
			}
			return from.Interface(), nil
		},
		Result:  value,
		TagName: "json",
	})

	if err != nil {
		return err
	}

	if response.dataKeyPath != "" {
		return decoder.Decode(values[response.dataKeyPath])
	}

	return decoder.Decode(values)
}
