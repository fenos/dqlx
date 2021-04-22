package deku_test

import (
	"github.com/fenos/deku"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueryBuilder_Query(t *testing.T) {

	t.Run("simple query with root operation", func(t *testing.T) {
		query, variables, err := deku.
			Query(deku.EqFunc("name@en", "Blade Runner")).
			Name("bladerunner").
			Fields(`
				uid
				name
				initial_release_date
				netflix_id 
			`).
			ToDQL()

		require.NoError(t, err)
		require.Equal(t, variables, map[string]interface{}{
			"$0": "Blade Runner",
		})

		writer := deku.NewWriter()
		writer.AddLine("query Bladerunner($0: string) {")
		mainOpWriter := writer.AddIndentedLine("bladerunner(func: eq(name@en,$0)) {")
		mainOpWriter.AddIndentedLine("uid")
		mainOpWriter.AddIndentedLine("name")
		mainOpWriter.AddIndentedLine("initial_release_date")
		mainOpWriter.AddIndentedLine("netflix_id")
		writer.AddIndentedLine("}")
		writer.AddLine("}")

		require.Equal(t, writer.ToString(), query)
	})

	t.Run("deeply nested query", func(t *testing.T) {
		query, variables, err := deku.
			Query(deku.EqFunc("name@en", "Blade Runner")).
			Name("bladerunner").
			Fields(`
				uid
				name
				initial_release_date
				netflix_id
			`).
			Edge("deep1", func(builder *deku.QueryBuilder) {
				builder.
					Fields(`
						uid
						name
						second
						third
					`).
					Edge("deep2", func(builder *deku.QueryBuilder) {
						builder.Fields(`
							uid
							name
							forth
							fifth
						`)
					})
			}).
			ToDQL()

		require.NoError(t, err)
		require.Equal(t, variables, map[string]interface{}{
			"$0": "Blade Runner",
		})

		writer := deku.NewWriter()
		writer.AddLine("query Bladerunner($0: string) {")
		mainOpWriter := writer.AddIndentedLine("bladerunner(func: eq(name@en,$0)) {")
		mainOpWriter.AddIndentedLine("uid")
		mainOpWriter.AddIndentedLine("name")
		mainOpWriter.AddIndentedLine("initial_release_date")
		mainOpWriter.AddIndentedLine("netflix_id")
		deep1Writer := mainOpWriter.AddIndentedLine("deep1 {")
		deep1Writer.AddIndentedLine("uid")
		deep1Writer.AddIndentedLine("name")
		deep1Writer.AddIndentedLine("second")
		deep1Writer.AddIndentedLine("third")
		deep2Writer := deep1Writer.AddIndentedLine("deep2 {")
		deep2Writer.AddIndentedLine("uid")
		deep2Writer.AddIndentedLine("name")
		deep2Writer.AddIndentedLine("forth")
		deep2Writer.AddIndentedLine("fifth")
		deep1Writer.AddIndentedLine("}")
		mainOpWriter.AddIndentedLine("}")
		writer.AddIndentedLine("}")
		writer.AddLine("}")

		require.Equal(t, writer.ToString(), query)
	})

	t.Run("split and merge multiple edges", func(t *testing.T) {

		query := deku.
			Query(deku.EqFunc("name@en", "Blade Runner")).
			Name("bladerunner").
			Fields(`
				uid
				name
				initial_release_date
				netflix_id
			`)

		deep1Query := deku.
			Query(nil).
			Name("deep1").
			Fields(`
					uid
					name
					second
					third
			`)

		deep2Query := deku.
			Query(nil).
			Name("deep2").
			Fields(`
					uid
					name
					forth
					fifth
			`)

		query.MergeEdge(deep1Query)
		deep1Query.MergeEdge(deep2Query)

		fullQuery, variables, err := query.ToDQL()

		require.NoError(t, err)
		require.Equal(t, variables, map[string]interface{}{
			"$0": "Blade Runner",
		})

		writer := deku.NewWriter()
		writer.AddLine("query Bladerunner($0: string) {")
		mainOpWriter := writer.AddIndentedLine("bladerunner(func: eq(name@en,$0)) {")
		mainOpWriter.AddIndentedLine("uid")
		mainOpWriter.AddIndentedLine("name")
		mainOpWriter.AddIndentedLine("initial_release_date")
		mainOpWriter.AddIndentedLine("netflix_id")
		deep1Writer := mainOpWriter.AddIndentedLine("deep1 {")
		deep1Writer.AddIndentedLine("uid")
		deep1Writer.AddIndentedLine("name")
		deep1Writer.AddIndentedLine("second")
		deep1Writer.AddIndentedLine("third")
		deep2Writer := deep1Writer.AddIndentedLine("deep2 {")
		deep2Writer.AddIndentedLine("uid")
		deep2Writer.AddIndentedLine("name")
		deep2Writer.AddIndentedLine("forth")
		deep2Writer.AddIndentedLine("fifth")
		deep1Writer.AddIndentedLine("}")
		mainOpWriter.AddIndentedLine("}")
		writer.AddIndentedLine("}")
		writer.AddLine("}")

		require.Equal(t, writer.ToString(), fullQuery)
	})
}

func TestQuery_Filter(t *testing.T) {
	t.Run("add a simple filter", func(t *testing.T) {
		query, variables, err := deku.
			Query(deku.EqFunc("name@en", "Blade Runner")).
			Name("bladerunner").
			Filter(deku.EqFunc("initial_release_date", "2010-06-10")).
			Fields(`
				uid
				name
				initial_release_date
			`).
			ToDQL()

		require.NoError(t, err)
		require.Equal(t, variables, map[string]interface{}{
			"$0": "Blade Runner",
			"$1": "2010-06-10",
		})

		writer := deku.NewWriter()
		writer.AddLine("query Bladerunner($0: string,$1: string) {")
		mainOpWriter := writer.AddIndentedLine("bladerunner(func: eq(name@en,$0)) @filter(eq(initial_release_date,$1)) {")
		mainOpWriter.AddIndentedLine("uid")
		mainOpWriter.AddIndentedLine("name")
		mainOpWriter.AddIndentedLine("initial_release_date")
		writer.AddIndentedLine("}")
		writer.AddLine("}")

		require.Equal(t, writer.ToString(), query)
	})

	t.Run("multiples AND filters", func(t *testing.T) {
		query, variables, err := deku.
			Query(deku.EqFunc("name@en", "Blade Runner")).
			Name("bladerunner").
			Filter(deku.EqFunc("initial_release_date", "2010-06-10")).
			Filter(deku.EqFunc("initial_release_date", "2010-05-13")).
			Fields(`
				uid
				name
				initial_release_date
			`).
			ToDQL()

		require.NoError(t, err)
		require.Equal(t, variables, map[string]interface{}{
			"$0": "Blade Runner",
			"$1": "2010-06-10",
			"$2": "2010-05-13",
		})

		writer := deku.NewWriter()
		writer.AddLine("query Bladerunner($0: string,$1: string,$2: string) {")
		mainOpWriter := writer.AddIndentedLine("bladerunner(func: eq(name@en,$0)) @filter(eq(initial_release_date,$1) AND eq(initial_release_date,$2)) {")
		mainOpWriter.AddIndentedLine("uid")
		mainOpWriter.AddIndentedLine("name")
		mainOpWriter.AddIndentedLine("initial_release_date")
		writer.AddIndentedLine("}")
		writer.AddLine("}")

		require.Equal(t, writer.ToString(), query)
	})

	t.Run("multiples AND / OR filters", func(t *testing.T) {
		query, variables, err := deku.
			Query(deku.EqFunc("name@en", "Blade Runner")).
			Name("bladerunner").
			Filter(deku.EqFunc("initial_release_date", "2010-06-10")).
			Filter(deku.EqFunc("initial_release_date", "2010-05-13")).
			OrFilter(deku.EqFunc("initial_release_date", "2020-05-13")).
			Fields(`
				uid
				name
				initial_release_date
			`).
			ToDQL()

		require.NoError(t, err)
		require.Equal(t, variables, map[string]interface{}{
			"$0": "Blade Runner",
			"$1": "2010-06-10",
			"$2": "2010-05-13",
			"$3": "2020-05-13",
		})

		writer := deku.NewWriter()
		writer.AddLine("query Bladerunner($0: string,$1: string,$2: string,$3: string) {")
		mainOpWriter := writer.AddIndentedLine("bladerunner(func: eq(name@en,$0)) @filter(eq(initial_release_date,$1) AND eq(initial_release_date,$2) OR eq(initial_release_date,$3)) {")
		mainOpWriter.AddIndentedLine("uid")
		mainOpWriter.AddIndentedLine("name")
		mainOpWriter.AddIndentedLine("initial_release_date")
		writer.AddIndentedLine("}")
		writer.AddLine("}")

		require.Equal(t, writer.ToString(), query)
	})

	t.Run("single filter group", func(t *testing.T) {
		query, variables, err := deku.
			Query(deku.EqFunc("name@en", "Blade Runner")).
			Name("bladerunner").
			FilterGroup(func(group *deku.FilterGroup) {
				group.
					Filter(deku.EqFunc("initial_release_date", "2010-05-13")).
					OrFilter(deku.EqFunc("initial_release_date", "2020-05-13"))
			}).
			Fields(`
				uid
				name
				initial_release_date
			`).
			ToDQL()

		require.NoError(t, err)
		require.Equal(t, variables, map[string]interface{}{
			"$0": "Blade Runner",
			"$1": "2010-05-13",
			"$2": "2020-05-13",
		})

		writer := deku.NewWriter()
		writer.AddLine("query Bladerunner($0: string,$1: string,$2: string) {")
		mainOpWriter := writer.AddIndentedLine("bladerunner(func: eq(name@en,$0)) @filter((eq(initial_release_date,$1) OR eq(initial_release_date,$2))) {")
		mainOpWriter.AddIndentedLine("uid")
		mainOpWriter.AddIndentedLine("name")
		mainOpWriter.AddIndentedLine("initial_release_date")
		writer.AddIndentedLine("}")
		writer.AddLine("}")

		require.Equal(t, writer.ToString(), query)
	})

	t.Run("multiples AND / OR filters", func(t *testing.T) {
		query, variables, err := deku.
			Query(deku.EqFunc("name@en", "Blade Runner")).
			Name("bladerunner").
			Filter(deku.EqFunc("initial_release_date", "2010-06-10")).
			FilterGroup(func(group *deku.FilterGroup) {
				group.
					Filter(deku.EqFunc("initial_release_date", "2010-05-13")).
					OrFilter(deku.EqFunc("initial_release_date", "2020-05-13"))
			}).
			Fields(`
				uid
				name
				initial_release_date
			`).
			ToDQL()

		require.NoError(t, err)
		require.Equal(t, variables, map[string]interface{}{
			"$0": "Blade Runner",
			"$1": "2010-06-10",
			"$2": "2010-05-13",
			"$3": "2020-05-13",
		})

		writer := deku.NewWriter()
		writer.AddLine("query Bladerunner($0: string,$1: string,$2: string,$3: string) {")
		mainOpWriter := writer.AddIndentedLine("bladerunner(func: eq(name@en,$0)) @filter(eq(initial_release_date,$1) AND (eq(initial_release_date,$2) OR eq(initial_release_date,$3))) {")
		mainOpWriter.AddIndentedLine("uid")
		mainOpWriter.AddIndentedLine("name")
		mainOpWriter.AddIndentedLine("initial_release_date")
		writer.AddIndentedLine("}")
		writer.AddLine("}")

		require.Equal(t, writer.ToString(), query)
	})
}
