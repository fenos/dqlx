package dqlx_test

import (
	dql "github.com/fenos/dqlx"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTypeBuilder_Fields(t *testing.T) {

	t.Run("boolean", func(t *testing.T) {
		actor := dql.NewTypeBuilder("Actor")
		actor.Bool("verified")
		actor.Bool("married").Index()

		expectedPredicates := dql.Minify(`
			Actor.verified: bool .
			Actor.married: bool @index() .
		`)

		expectedType := dql.Minify(`
			type Actor {
				Actor.verified
				Actor.married
			}
		`)

		dqlType, err := actor.ToString()
		require.NoError(t, err)
		require.Equal(t, expectedType, dql.Minify(dqlType))

		predicates := actor.PredicatesToString()
		require.NoError(t, err)
		require.Equal(t, expectedPredicates, dql.Minify(predicates))
	})

	t.Run("string", func(t *testing.T) {
		actor := dql.NewTypeBuilder("Actor")
		actor.String("name")
		actor.String("description").IndexFulltext().IndexTerm().IndexHash().IndexTrigram().IndexExact()

		expectedPredicates := dql.Minify(`
			Actor.name: string .
			Actor.description: string @index(fulltext,term,hash,trigram,exact) .
		`)

		expectedType := dql.Minify(`
			type Actor {
				Actor.name
				Actor.description
			}
		`)

		dqlType, err := actor.ToString()
		require.NoError(t, err)
		require.Equal(t, expectedType, dql.Minify(dqlType))

		predicates := actor.PredicatesToString()
		require.NoError(t, err)
		require.Equal(t, expectedPredicates, dql.Minify(predicates))
	})

	t.Run("int", func(t *testing.T) {
		actor := dql.NewTypeBuilder("Actor")
		actor.Int("age")
		actor.Int("films").Index()

		expectedPredicates := dql.Minify(`
			Actor.age: int .
			Actor.films: int @index() .
		`)

		expectedType := dql.Minify(`
			type Actor {
				Actor.age
				Actor.films
			}
		`)

		dqlType, err := actor.ToString()
		require.NoError(t, err)
		require.Equal(t, expectedType, dql.Minify(dqlType))

		predicates := actor.PredicatesToString()
		require.NoError(t, err)
		require.Equal(t, expectedPredicates, dql.Minify(predicates))
	})

	t.Run("float", func(t *testing.T) {
		actor := dql.NewTypeBuilder("Actor")
		actor.Float("float1")
		actor.Float("float2").Index()

		expectedPredicates := dql.Minify(`
			Actor.float1: float .
			Actor.float2: float @index() .
		`)

		expectedType := dql.Minify(`
			type Actor {
				Actor.float1
				Actor.float2
			}
		`)

		dqlType, err := actor.ToString()
		require.NoError(t, err)
		require.Equal(t, expectedType, dql.Minify(dqlType))

		predicates := actor.PredicatesToString()
		require.NoError(t, err)
		require.Equal(t, expectedPredicates, dql.Minify(predicates))
	})

	t.Run("datetime", func(t *testing.T) {
		actor := dql.NewTypeBuilder("Actor")
		actor.DateTime("birthday")
		actor.DateTime("married_at").IndexYear().IndexMonth().IndexDay().IndexHour().IndexDay()

		expectedPredicates := dql.Minify(`
			Actor.birthday: datetime .
			Actor.married_at: datetime @index(year,month,day,hour) .
		`)

		expectedType := dql.Minify(`
			type Actor {
				Actor.birthday
				Actor.married_at
			}
		`)

		dqlType, err := actor.ToString()
		require.NoError(t, err)
		require.Equal(t, expectedType, dql.Minify(dqlType))

		predicates := actor.PredicatesToString()
		require.NoError(t, err)
		require.Equal(t, expectedPredicates, dql.Minify(predicates))
	})

	t.Run("password", func(t *testing.T) {
		actor := dql.NewTypeBuilder("Actor")
		actor.Password("password")
		actor.Password("secret").Index()

		expectedPredicates := dql.Minify(`
			Actor.password: password .
			Actor.secret: password @index() .
		`)

		expectedType := dql.Minify(`
			type Actor {
				Actor.password
				Actor.secret
			}
		`)

		dqlType, err := actor.ToString()
		require.NoError(t, err)
		require.Equal(t, expectedType, dql.Minify(dqlType))

		predicates := actor.PredicatesToString()
		require.NoError(t, err)
		require.Equal(t, expectedPredicates, dql.Minify(predicates))
	})

	t.Run("uid", func(t *testing.T) {
		actor := dql.NewTypeBuilder("Actor")
		actor.UID("id")
		actor.UID("uid").Index()

		expectedPredicates := dql.Minify(`
			Actor.id: uid .
			Actor.uid: uid @index() .
		`)

		expectedType := dql.Minify(`
			type Actor {
				Actor.id
				Actor.uid
			}
		`)

		dqlType, err := actor.ToString()
		require.NoError(t, err)
		require.Equal(t, expectedType, dql.Minify(dqlType))

		predicates := actor.PredicatesToString()
		require.NoError(t, err)
		require.Equal(t, expectedPredicates, dql.Minify(predicates))
	})

	t.Run("type predicate", func(t *testing.T) {
		actor := dql.NewTypeBuilder("Actor")
		actor.UID("id")
		actor.String("name").Lang()
		actor.Type("film", "Film")
		actor.Type("rewards", "Rewards").Reverse().List()

		expectedPredicates := dql.Minify(`
			Actor.id: uid .
			Actor.name: string @lang .
			Actor.film: uid .
			Actor.rewards: [uid] @reverse .
		`)

		expectedType := dql.Minify(`
			type Actor {
				Actor.id
				Actor.name
				Actor.film: Film
				<~Actor.rewards>: [Rewards]
			}
		`)

		dqlType, err := actor.ToString()
		require.NoError(t, err)
		require.Equal(t, expectedType, dql.Minify(dqlType))

		predicates := actor.PredicatesToString()
		require.NoError(t, err)
		require.Equal(t, expectedPredicates, dql.Minify(predicates))
	})

	t.Run("global predicate on a type", func(t *testing.T) {
		schema := dql.NewSchema()

		name := schema.Predicate("name", dql.ScalarString)

		schema.Type("User", func(user *dql.TypeBuilder) {
			user.Predicate(name)
			user.String("surname")
		})

		types, err := schema.TypesToString()

		expected := dql.Minify(`
			type User {
				name
				User.surname
			}
		`)

		require.NoError(t, err)
		require.Equal(t, expected, dql.Minify(types))
	})

	// TODO:  geo indexes

	t.Run("can't register multiple QueryFields with the same name", func(t *testing.T) {
		actor := dql.NewTypeBuilder("Actor")
		actor.String("name")
		actor.String("name")

		_, err := actor.ToString()
		require.Error(t, err, "predicate 'Actor.name' already registered on type 'Actor'")
	})

	t.Run("can't register multiple QueryFields with the same name", func(t *testing.T) {
		actor := dql.NewTypeBuilder("Actor", dql.WithTypePrefix(false))
		actor.String("name")

		expectedPredicates := dql.Minify(`
			name: string .
		`)

		expectedType := dql.Minify(`
			type Actor {
				name
			}
		`)

		dqlType, err := actor.ToString()
		require.NoError(t, err)
		require.Equal(t, expectedType, dql.Minify(dqlType))

		predicates := actor.PredicatesToString()
		require.NoError(t, err)
		require.Equal(t, expectedPredicates, dql.Minify(predicates))
	})
}
