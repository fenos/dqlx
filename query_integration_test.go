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
				"uid":                  result[0]["uid"],
				"name@en":              "Blade Runner",
				"initial_release_date": "1982-06-25T00:00:00Z",
				"netflix_id":           "70083726",
			},
		}
		require.NoError(suite.T(), err)
		require.ElementsMatch(suite.T(), result, expected)
	})

	suite.Run("By Id", func() {
		var idRecords []map[string]interface{}
		_, err := suite.db.Query(dqlx.EqFn("name@en", "Blade Runner")).
			Select(`
				uid
				name@en
				initial_release_date
				netflix_id
			`).
			UnmarshalInto(&idRecords).
			Execute(ctx)

		require.NoError(suite.T(), err)

		var result []map[string]interface{}

		_, err = suite.db.
			Query(dqlx.UIDFn(idRecords[0]["uid"])).
			Fields(`
				uid
				name@en
				initial_release_date
				netflix_id
			`).
			UnmarshalInto(&result).
			Execute(ctx)

		require.NoError(suite.T(), err)
		require.ElementsMatch(suite.T(), result, idRecords)
	})

	suite.Run("Anyofterms", func() {
		var result []map[string]interface{}

		_, err := suite.db.
			Query(dqlx.AllOfTermsFn("name@en", "Blade Runner")).
			Fields(`
				name@en
				initial_release_date
				netflix_id
			`).
			UnmarshalInto(&result).
			Execute(ctx)

		expected := []map[string]interface{}{
			{
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

	suite.Run("Applying filters", func() {
		var result []map[string]interface{}

		_, err := suite.db.
			Query(dqlx.EqFn("name@en", "Ridley Scott")).
			Fields(`
				uid
				name@en
				initial_release_date
				netflix_id
			`).
			Edge("director.film", dqlx.Fields(`
				name@en
      			initial_release_date
			`), dqlx.Le{"initial_release_date": "2000"}).
			UnmarshalInto(&result).
			Execute(ctx)

		require.NoError(suite.T(), err)
		require.Len(suite.T(), result, 2)
		require.Len(suite.T(), result[0]["director.film"], 10)

		suite.recordsContainsProperties(result[0]["director.film"].([]interface{}), []string{
			"name@en",
			"initial_release_date",
		})

		require.NotContains(suite.T(), result[1], "director.film")
	})

	suite.Run("Language support", func() {
		var result []map[string]interface{}

		_, err := suite.db.
			Query(dqlx.AllOfTermsFn("name@en", "Farhan Akhtar")).
			Fields(`
				name@hi
    			name@en

			`).
			Edge("director.film", dqlx.Fields(`
				  name@ru:hi:en
				  name@en
				  name@hi
				  name@ru
			`)).
			UnmarshalInto(&result).
			Execute(ctx)

		require.NoError(suite.T(), err)
		require.Len(suite.T(), result, 1)

		suite.recordsContainsProperties(result, []string{
			"name@en",
			"name@hi",
		})

		suite.recordsContainsProperties(result[0]["director.film"], []string{
			"name@ru:hi:en",
			"name@en",
			"name@hi",
			"name@ru",
		})
	})

	suite.Run("Count", func() {
		var result []map[string]interface{}

		_, err := suite.db.
			Query(dqlx.EqFn("name@en", "Ridley Scott")).
			Fields(
				"uid",
				"name@en",
				"initial_release_date",
				"netflix_id",
				dqlx.Count("director.film"),
			).
			UnmarshalInto(&result).
			Execute(ctx)

		require.NoError(suite.T(), err)
		require.Len(suite.T(), result, 2)

		suite.containsProperties(result[0], []string{
			"name@en",
			"count(director.film)",
		})
	})

	suite.Run("Aliases", func() {
		var result []map[string]interface{}

		_, err := suite.db.
			Query(dqlx.EqFn("name@en", "Ridley Scott")).
			Select(
				"uid",
				"name@en",
				"alias_name:name@en",
				"alias_id:uid",
				":uid@en",
				dqlx.Alias("count", dqlx.Count("director.film")),
				dqlx.Alias("first_name_alias", "name@en"),
			).
			UnmarshalInto(&result).
			Execute(ctx)

		require.NoError(suite.T(), err)
		require.Len(suite.T(), result, 2)

		suite.containsProperties(result[0], []string{
			"alias_id",
			"name@en",
			"count",
			"alias_name",
			"first_name_alias",
		})
	})
}

func (suite *QueryIntegrationTest) recordsContainsProperties(slice interface{}, properties []string) {
	var records []interface{}

	switch castSlices := slice.(type) {
	case []interface{}:
		records = castSlices
	case []map[string]interface{}:
		for _, item := range castSlices {
			records = append(records, item)
		}
	default:
		require.Fail(suite.T(), "contains properties invalid value %v", slice)
	}

	for _, record := range records {
		for _, property := range properties {
			require.Contains(suite.T(), record, property)
			require.NotEmpty(suite.T(), record.(map[string]interface{})[property])
		}
	}
}

func (suite *QueryIntegrationTest) containsProperties(record interface{}, properties []string) {
	for _, property := range properties {
		require.Contains(suite.T(), record, property)
	}
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(QueryIntegrationTest))
}
