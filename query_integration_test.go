package dqlx_test

import (
	"context"
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/fenos/dqlx/testdata"

	"github.com/dgraph-io/dgo/v200/protos/api"

	"github.com/fenos/dqlx"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

//go:embed testdata/dataset.json
var data string

type QueryIntegrationTest struct {
	suite.Suite
	db dqlx.DB
}

func (suite *QueryIntegrationTest) SetupTest() {
	ctx := context.TODO()
	db, err := dqlx.Connect("localhost:9080")
	require.NoError(suite.T(), err)

	suite.db = db

	schema := testdata.TestSchema()
	err = db.GetDgraph().Alter(ctx, &api.Operation{
		Schema: schema,
		//DropOp:          api.Operation_ALL,
		RunInBackground: false,
	})

	require.NoError(suite.T(), err)

	var structure interface{}
	err = json.Unmarshal([]byte(data), &structure)
	require.NoError(suite.T(), err)

	resp, err := suite.db.Mutation().Set(structure).Execute(ctx)
	require.NoError(suite.T(), err)

	println(resp.Raw.Uids)

}

func (suite *QueryIntegrationTest) TearDownTest() {
	err := suite.db.GetDgraph().Alter(context.Background(), &api.Operation{
		DropOp: api.Operation_ALL,
	})
	require.NoError(suite.T(), err)
}

func (suite *QueryIntegrationTest) TestFundamentals() {
	ctx := context.TODO()
	suite.Run("First Example", func() {
		var result []map[string]interface{}

		_, err := suite.db.
			Query(dqlx.EqFn("name@en", "Blade Runner")).
			Fields(`
				uid
				name@en
				initial_release_date
				netflix_id
			`).
			UnmarshalInto(&result).
			Execute(ctx)

		expected := []map[string]interface{}{
			{
				"uid":                  "0x1",
				"name@en":              "Blade Runner",
				"initial_release_date": "1982-06-25T00:00:00Z",
				"netflix_id":           "70083726",
			},
		}
		require.NoError(suite.T(), err)
		require.ElementsMatch(suite.T(), result, expected)
	})

	suite.Run("By Id", func() {
		var result []map[string]interface{}

		_, err := suite.db.
			Query(dqlx.UIDFn("0x1")).
			Fields(`
				uid
				name@en
				initial_release_date
				netflix_id
			`).
			UnmarshalInto(&result).
			Execute(ctx)

		expected := []map[string]interface{}{
			{
				"uid":                  "0x1",
				"name@en":              "Blade Runner",
				"initial_release_date": "1982-06-25T00:00:00Z",
				"netflix_id":           "70083726",
			},
		}
		require.NoError(suite.T(), err)
		require.ElementsMatch(suite.T(), result, expected)
	})

	suite.Run("Anyofterms", func() {
		var result []map[string]interface{}

		_, err := suite.db.
			Query(dqlx.AllOfTermsFn("name@en", "Blade Runner")).
			Fields(`
				uid
				name@en
				initial_release_date
				netflix_id
			`).
			UnmarshalInto(&result).
			Execute(ctx)

		expected := []map[string]interface{}{
			{
				"uid":                  "0x1",
				"name@en":              "Blade Runner",
				"initial_release_date": "1982-06-25T00:00:00Z",
				"netflix_id":           "70083726",
			},
		}
		require.NoError(suite.T(), err)
		require.ElementsMatch(suite.T(), result, expected)
	})

	suite.Run("Expanding edge", func() {
		var result []map[string]interface{}

		_, err := suite.db.
			Query(dqlx.EqFn("name@en", "Blade Runner")).
			Fields(`
				uid
				name@en
				initial_release_date
				netflix_id
			`).
			Edge("starring").
			Edge("starring->performance.actor", dqlx.Fields("name@en")).
			Edge("starring->performance.character", dqlx.Fields("name@en")).
			UnmarshalInto(&result).
			Execute(ctx)

		require.NoError(suite.T(), err)

		starring := result[0]["starring"].([]interface{})
		require.Len(suite.T(), starring, 15)
	})
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(QueryIntegrationTest))
}
