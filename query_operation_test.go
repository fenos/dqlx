package dqlx_test

import (
	"testing"

	dql "github.com/getplexy/dqlx"
	"github.com/stretchr/testify/require"
)

func Test_Multiple_Blocks(t *testing.T) {
	query1 := dql.
		QueryEdge("bladerunner", dql.EqFn("item", "value")).
		Fields(`
			uid
			super_alias:name
			initial_release_date
			d AS netflix_id
		`).
		Filter(dql.Eq{"field1": "value1"})

	query2 := dql.
		QueryEdge("bladerunner2", dql.EqFn("item", "value")).
		Fields(`
			uid
			name
		`).
		Filter(dql.Eq{"field1": "value1"})

	query, variables, err := dql.QueriesToDQL(query1, query2)

	require.NoError(t, err)
	require.Equal(t, map[string]string{
		"$0": "value",
		"$1": "value1",
		"$2": "value",
		"$3": "value1",
	}, variables.ToVariables())

	expected := dql.Minify(`
		query Bladerunner_Bladerunner2($0:string, $1:string, $2:string, $3:string) {
			<bladerunner>(func: eq(<item>,$0)) @filter(eq(<field1>,$1)) {
				<uid>
				<super_alias>:<name>
				<initial_release_date>
				d AS <netflix_id>
			}

			<bladerunner2>(func: eq(<item>,$2)) @filter(eq(<field1>,$3)) {
				<uid>
				<name>
			}
		}
	`)

	require.Equal(t, expected, query)
}

func Test_Multiple_Blocks_With_Select(t *testing.T) {
	q1 := dql.Query(dql.EqFn("id", "id_a")).
		Select(dql.As("id_a", "id"))

	q2 := dql.Query(dql.EqFn("id", "id_b")).
		Select(dql.As("id_b", "id"))

	query, variables, err := dql.QueriesToDQL(q1, q2)

	require.NoError(t, err)

	require.Equal(t, map[string]string{
		"$0": "id_a",
		"$1": "id_b",
	}, variables.ToVariables())

	expected := dql.Minify(`
		query Rootquery_Rootquery_1($0:string, $1:string) { 
          <rootQuery>(func: eq(<id>,$0)) { 
            id_a AS <id> 
          } 
          <rootQuery_1>(func: eq(<id>,$1)) {
            id_b AS <id>
          } 
        }
	`)

	require.Equal(t, expected, query)
}
