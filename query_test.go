package deku_test

import (
	dql "github.com/fenos/deku"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_Simple_Query(t *testing.T) {
	query, variables, err := dql.
		Query("bladerunner", dql.EqFn("item", "value")).
		Fields(`
			uid
			name
			initial_release_date
			netflix_id
		`).
		Filter(dql.Eq{"field1": "value1"}).
		ToDQL()

	require.NoError(t, err)
	require.Equal(t, map[string]interface{}{
		"$0": "value",
		"$1": "value1",
	}, variables)

	expected := dql.Minify(`
		query Bladerunner($0:string, $1:string) {
			bladerunner(func: eq(item,$0)) @filter(eq(field1,$1)) {
				uid
				name
				initial_release_date
				netflix_id
			}
		}
	`)

	require.Equal(t, expected, query)
}

func Test_Query_Nested(t *testing.T) {
	query, variables, err := dql.
		Query("bladerunner", dql.EqFn("name@en", "Blade Runner")).
		Fields(`
			uid
			name
			initial_release_date
			netflix_id
		`).
		Edge("authors", dql.Fields(`
			uid
			name
			surname
			age
		`)).
		Edge("actors", dql.Fields(`
			uid
			surname
		`)).
		Edge("actors->rewards", dql.Fields(`
			uid
			points
		`)).
		Edge("actors->rewards->venues", dql.Fields(`
			street
		`)).
		ToDQL()

	require.NoError(t, err)
	require.Equal(t, map[string]interface{}{
		"$0": "Blade Runner",
	}, variables)

	expected := dql.Minify(`
		query Bladerunner($0:string) {
			bladerunner(func: eq(name@en,$0)) {
				uid
				name
				initial_release_date
				netflix_id
				authors {
					uid
					name
					surname
					age
				}
				actors {
					uid
					surname
					rewards {
						uid
						points
						venues {
							street
						}
					}
				}
			}
		}
	`)

	require.Equal(t, expected, query)
}

func Test_Query_Filter_Nested(t *testing.T) {
	query, variables, err := dql.
		Query("bladerunner", dql.EqFn("name@en", "Blade Runner")).
		Fields(`
			uid
			name
			initial_release_date
			netflix_id
		`).
		Edge("authors", dql.Fields(`
			uid
			name
			surname
			age
		`), dql.Eq{"age": 20}).
		Edge("actors", dql.Fields(`
			uid
			surname
			age
		`), dql.Gt{"age": []int{18, 20, 30}}).
		Edge(dql.EdgePath("actors", "rewards"), dql.Fields(`
			uid
			points
		`), dql.Gt{"points": 3}).
		ToDQL()

	require.NoError(t, err)
	require.Equal(t, map[string]interface{}{
		"$0": "Blade Runner",
		"$1": 20,
		"$2": 18,
		"$3": 20,
		"$4": 30,
		"$5": 3,
	}, variables)

	expected := dql.Minify(`
		query Bladerunner($0:string, $1:int, $2:int, $3:int, $4:int, $5:int) {
			bladerunner(func: eq(name@en,$0)) {
				uid
				name
				initial_release_date
				netflix_id
				authors @filter(eq(age,$1)) {
					uid
					name
					surname
					age
				}
				actors @filter(gt(age,[$2,$3,$4])) {
					uid
					surname
					age
					rewards @filter(gt(points,$5)) {
						uid
						points
					}
				}
			}
		}
	`)

	require.Equal(t, expected, query)
}

func Test_Query_Connecting_Filter(t *testing.T) {
	query, variables, err := dql.
		Query("bladerunner", dql.EqFn("name@en", "Blade Runner")).
		Fields(`
			uid
			name
			initial_release_date
			netflix_id
		`).
		Filter(dql.Or{
			dql.Eq{"name": "actor1"},
			dql.Eq{"name": "actor2"},
		}).
		Edge("authors", dql.Fields(`
			uid
			name
			surname
			age
		`), dql.Or{
			dql.And{
				dql.Eq{"name": "author3"},
				dql.Gt{"age": 20},
			},
			dql.And{
				dql.Eq{"name": "author4"},
				dql.Lt{"age": 50},
			},
		}).
		ToDQL()

	require.NoError(t, err)
	require.Equal(t, map[string]interface{}{
		"$0": "Blade Runner",
		"$1": "actor1",
		"$2": "actor2",
		"$3": "author3",
		"$4": 20,
		"$5": "author4",
		"$6": 50,
	}, variables)

	expected := dql.Minify(`
		query Bladerunner($0:string, $1:string, $2:string, $3:string, $4:int, $5:string, $6:int) {
			bladerunner(func: eq(name@en,$0)) @filter((eq(name,$1) OR eq(name,$2))) {
				uid
				name
				initial_release_date
				netflix_id
				authors @filter(((eq(name,$3) AND gt(age,$4)) OR (eq(name,$5) AND lt(age,$6)))) {
					uid
					name
					surname
					age
				}
			}
		}
	`)

	require.Equal(t, expected, query)
}
