package dqlx_test

import (
	dql "github.com/fenos/dqlx"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_Simple_Query(t *testing.T) {
	query, variables, err := dql.
		QueryEdge("bladerunner", dql.EqFn("item", "value")).
		Fields(`
			uid
			name
			initial_release_date
			netflix_id
		`).
		Filter(dql.Eq{"field1": "value1"}).
		ToDQL()

	require.NoError(t, err)
	require.Equal(t, map[string]string{
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
		QueryEdge("bladerunner", dql.EqFn("name@en", "Blade Runner")).
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
	require.Equal(t, map[string]string{
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
		QueryEdge("bladerunner", dql.EqFn("name@en", "Blade Runner")).
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
	require.Equal(t, map[string]string{
		"$0": "Blade Runner",
		"$1": "20",
		"$2": "18",
		"$3": "20",
		"$4": "30",
		"$5": "3",
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
		QueryEdge("bladerunner", dql.EqFn("name@en", "Blade Runner")).
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
	require.Equal(t, map[string]string{
		"$0": "Blade Runner",
		"$1": "actor1",
		"$2": "actor2",
		"$3": "author3",
		"$4": "20",
		"$5": "author4",
		"$6": "50",
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

func Test_Query_Pagination(t *testing.T) {
	query, variables, err := dql.
		QueryEdge("bladerunner", dql.EqFn("name@en", "Blade Runner")).
		Fields(`
			uid
			name
			initial_release_date
			netflix_id
		`).
		Paginate(dql.Pagination{
			First:  20,
			Offset: 1,
			After:  "4567",
		}).
		EdgeFn("authors", func(builder dql.QueryBuilder) dql.QueryBuilder {
			return builder.
				Fields(`
					uid
					name
					surname
					age
				`).
				Paginate(dql.Pagination{
					First:  10,
					Offset: 2,
					After:  "1234",
				})
		}).
		EdgeFn("actors", func(builder dql.QueryBuilder) dql.QueryBuilder {
			return builder.
				Fields(`
					uid
					surname
					age
				`).
				Filter(dql.Gt{"age": 30}).
				Paginate(dql.Pagination{
					First:  2,
					Offset: 3,
					After:  "45",
				})
		}).
		ToDQL()

	require.NoError(t, err)
	require.Equal(t, map[string]string{
		"$0":  "Blade Runner",
		"$1":  "20",
		"$2":  "1",
		"$3":  "4567",
		"$4":  "10",
		"$5":  "2",
		"$6":  "1234",
		"$7":  "2",
		"$8":  "3",
		"$9":  "45",
		"$10": "30",
	}, variables)

	expected := dql.Minify(`
		query Bladerunner($0:string, $1:int, $2:int, $3:string, $4:int, $5:int, $6:string, $7:int, $8:int, $9:string, $10:int) {
			bladerunner(func: eq(name@en,$0),first:$1,offset:$2,after:$3) {
				uid
				name
				initial_release_date
				netflix_id
				authors(first:$4,offset:$5,after:$6) {
					uid
					name
					surname
					age
				}
				actors(first:$7,offset:$8,after:$9) @filter(gt(age,$10))  {
					uid
					surname
					age
				}
			}
		}
	`)

	require.Equal(t, expected, query)
}

func Test_Query_Variable(t *testing.T) {
	variable := dql.Variable(dql.EqFn("name", "test")).
		Edge("film").
		Edge("film->performance", dql.Fields(`
			 D AS genre
		`))

	query, variables, err := dql.
		QueryEdge("bladerunner", dql.EqFn("item", "value")).
		Fields(`
			uid
			name
			initial_release_date
			netflix_id
		`).
		Variable(variable).
		Filter(dql.Eq{"field1": dql.Expr("D")}).
		ToDQL()

	require.NoError(t, err)
	require.Equal(t, map[string]string{
		"$0": "test",
		"$1": "value",
	}, variables)

	expected := dql.Minify(`
		query Bladerunner($0:string, $1:string) {
			var(func: eq(name,$0)) {
				film {
					performance {
					 D AS genre
					}
				}
			}

			bladerunner(func: eq(item,$1)) @filter(eq(field1,D)) {
				uid
				name
				initial_release_date
				netflix_id
			}
		}
	`)

	require.Equal(t, expected, query)
}

func Test_Query_Value_Variable(t *testing.T) {
	variable := dql.Variable(dql.EqFn("name", "test")).
		Edge("film").
		Edge("film->performance", dql.Fields(`
			 D AS genre
		`))

	query, _, err := dql.
		QueryEdge("bladerunner", dql.EqFn("item", "value")).
		Fields(`
			uid
			name
			initial_release_date
			netflix_id
		`).
		Filter(dql.UID(dql.Val("D"))).
		Filter(dql.UID(dql.Expr("D"))).
		Variable(variable).
		ToDQL()

	expected := dql.Minify(`
		query Bladerunner($0:string, $1:string) {
			var(func: eq(name,$0)) {
				film {
					performance { 
						D AS genre
					}
				}
			}

			bladerunner(func: eq(item,$1)) @filter(uid(val(D)) AND uid(D)) {
				uid
				name
				initial_release_date
				netflix_id
			}
		}
	`)

	require.NoError(t, err)
	require.Equal(t, expected, query)
}

func Test_Query_OrderBy(t *testing.T) {
	query, variables, err := dql.
		QueryEdge("bladerunner", dql.EqFn("item", "value")).
		Fields(`
			uid
			name
			initial_release_date
			netflix_id
		`).
		OrderAsc("name").
		OrderDesc(dql.Val("initial_release_date")).
		Edge("films", dql.Fields(`
			id
			date
		`), dql.OrderDesc("date"), dql.Pagination{First: 10}).
		ToDQL()

	expected := dql.Minify(`
		query Bladerunner($0:string, $1:int) {
			bladerunner(func: eq(item,$0),orderasc:name,orderdesc:val(initial_release_date)) {
				uid
				name
				initial_release_date
				netflix_id
				films(first:$1)(orderdesc:date) {
					id
					date
				}
			}
		}
	`)

	require.Equal(t, map[string]string{
		"$0": "value",
		"$1": "10",
	}, variables)

	require.NoError(t, err)
	require.Equal(t, expected, query)
}

func Test_Query_GroupBy(t *testing.T) {
	variable := dql.Variable(dql.EqFn("name", "test")).
		Edge("film", dql.Fields(`
			 a AS genre
		`), dql.GroupBy("genre"))

	query, variables, err := dql.
		QueryEdge("bladerunner", dql.EqFn("item", "value")).
		Fields(`
			uid
			name
			initial_release_date
			netflix_id
		`).
		Edge("films", dql.Fields(`
			total_movies:val(a)
		`)).
		Variable(variable).
		ToDQL()

	expected := dql.Minify(`
		query Bladerunner($0:string, $1:string) {
			var(func: eq(name,$0)) {
				film @groupby(genre) {
					a AS genre
				}
			}

			bladerunner(func: eq(item,$1)) {
				uid
				name
				initial_release_date
				netflix_id
				films {
					total_movies:val(a)
				}
			}
		}
	`)

	require.Equal(t, map[string]string{
		"$0": "test",
		"$1": "value",
	}, variables)

	require.NoError(t, err)
	require.Equal(t, expected, query)
}

func Test_Query_Facets(t *testing.T) {
	query, variables, err := dql.
		QueryEdge("bladerunner", dql.EqFn("item", "value")).
		Facets().
		Fields(`
			uid
			name
			initial_release_date
			netflix_id
		`).
		Edge("films", dql.Fields(`
			name
		`), dql.Facets(dql.Eq{
			"close":    true,
			"relative": true,
		}), dql.Facets(dql.Expr("relative"))).
		ToDQL()

	expected := dql.Minify(`
		query Bladerunner($0:string, $1:bool, $2:bool) {
			bladerunner(func: eq(item,$0)) @facets {
				uid
				name
				initial_release_date
				netflix_id
				films @facets(eq(close,$1) AND eq(relative,$2)) @facets(relative) {
					name
				}
			}
		}
	`)

	require.Equal(t, map[string]string{
		"$0": "value",
		"$1": "true",
		"$2": "true",
	}, variables)

	require.NoError(t, err)
	require.Equal(t, expected, query)
}

func Test_Query_Edge_From_Query(t *testing.T) {
	edge := dql.QueryEdge("actors->rewards->venues", nil).
		Fields(`
			street
		`)

	query, variables, err := dql.
		QueryEdge("bladerunner", dql.EqFn("name@en", "Blade Runner")).
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
		EdgeFromQuery(edge).
		ToDQL()

	require.NoError(t, err)
	require.Equal(t, map[string]string{
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

func Test_List_function_Query(t *testing.T) {

	from := time.Date(2021, 04, 27, 0, 0, 0, 0, time.UTC)
	to := time.Date(2021, 04, 28, 0, 0, 0, 0, time.UTC)

	query, variables, err := dql.
		QueryEdge("bladerunner", dql.UIDFn("value")).
		Fields(`
			uid
			name
			initial_release_date
			netflix_id
		`).
		Filter(dql.UIDIn{"name": "value1"}).
		Filter(dql.Between("release_date", from, to)).
		ToDQL()

	require.NoError(t, err)
	require.Equal(t, map[string]string{
		"$0": "value",
		"$1": "value1",
		"$2": from.Format(time.RFC3339),
		"$3": to.Format(time.RFC3339),
	}, variables)

	expected := dql.Minify(`
		query Bladerunner($0:string, $1:string, $2:datetime, $3:datetime) {
			bladerunner(func: uid($0)) @filter(uid_in(name,$1) AND between(release_date,$2,$3)) {
				uid
				name
				initial_release_date
				netflix_id
			}
		}
	`)

	require.Equal(t, expected, query)
}
