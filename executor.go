package dqlx

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/dgraph-io/dgo/v200"
	"github.com/dgraph-io/dgo/v200/protos/api"
	"github.com/mitchellh/mapstructure"
	"reflect"
	"time"
)

type DGoExecutor struct {
	client *dgo.Dgraph
	tnx    *dgo.Txn

	readOnly bool
}

type ExecutorOptionFn func(executor *DGoExecutor)

func WithTnx(tnx *dgo.Txn) ExecutorOptionFn {
	return func(executor *DGoExecutor) {
		executor.tnx = tnx
	}
}

func WithClient(client *dgo.Dgraph) ExecutorOptionFn {
	return func(executor *DGoExecutor) {
		executor.client = client
	}
}

func WithReadOnly(readOnly bool) ExecutorOptionFn {
	return func(executor *DGoExecutor) {
		executor.readOnly = readOnly
	}
}

func NewDGoExecutor(client *dgo.Dgraph) *DGoExecutor {
	return &DGoExecutor{
		client: client,
	}
}

func (executor DGoExecutor) ExecuteQueries(ctx context.Context, queries ...QueryBuilder) (*QueryResponse, error) {
	if executor.client == nil {
		return nil, errors.New("cannot execute query without setting a dqlx. use DClient() to set one")
	}

	query, variables, err := QueriesToDQL(queries...)
	if err != nil {
		return nil, err
	}

	tx := executor.tnx

	if tx == nil {
		if executor.readOnly {
			tx = executor.client.NewReadOnlyTxn()
		} else {
			tx = executor.client.NewTxn()
		}
	}

	defer tx.Discard(ctx)

	resp, err := tx.QueryWithVars(ctx, query, variables)
	if err != nil {
		return nil, err
	}

	if !executor.readOnly {
		err := tx.Commit(ctx)
		if err != nil {
			return nil, err
		}
	}

	operationName := "data" // default operation name

	if len(queries) == 1 {
		operationName = queries[0].rootEdge.Name
	}

	queryResponse := &QueryResponse{
		operationName: operationName,
		Raw:           resp,
	}

	for _, query := range queries {
		if query.unmarshalInto != nil {
			err := queryResponse.Unmarshal(query.unmarshalInto)

			if err != nil {
				return nil, err
			}
		}
	}

	return queryResponse, nil
}

type QueryResponse struct {
	Raw           *api.Response
	operationName string
}

func (response QueryResponse) Unmarshal(value interface{}) error {
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

	return decoder.Decode(values[response.operationName])
}
