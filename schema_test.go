package deku_test

import (
	dql "github.com/fenos/deku"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_Schema_ToDQL(t *testing.T) {
	t.Run("add a simple Predicate", func(t *testing.T) {
		schema := dql.NewSchema()
		schema.Predicate("name", dql.ScalarString)
		schema.Predicate("surname", dql.ScalarString)

		dqlSchema, err := schema.ToDQL()

		expected := dql.Minify(`
			name:string .
			surname:string .
		`)

		require.NoError(t, err)
		require.Equal(t, expected, dqlSchema)
	})

	t.Run("add types and its predicates to the schema", func(t *testing.T) {
		schema := dql.NewSchema()

		author := schema.Type("Author")
		author.String("name")
		author.Int("age")

		dqlSchema, err := schema.ToDQL()

		expected := dql.Minify(`
			Author.name:string .
			Author.age:int .
			
			type Author {
				Author.name
				Author.age
			}
		`)

		require.NoError(t, err)
		require.Equal(t, expected, dqlSchema)
	})

	t.Run("should not duplicate predicates", func(t *testing.T) {
		schema := dql.NewSchema()

		author := schema.Type("Author", dql.WithTypePrefix(false))
		author.String("name")
		author.Int("age")

		film := schema.Type("Film", dql.WithTypePrefix(false))
		film.String("name")

		dqlSchema, err := schema.ToDQL()

		expected := dql.Minify(`
			name:string .
			age:int .
			
			type Author {
				name
				age
			}

			type Film {
				name
			}
		`)

		require.NoError(t, err)
		require.Equal(t, expected, dqlSchema)
	})
}
