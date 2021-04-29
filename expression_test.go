package dqlx_test

import (
	dql "github.com/fenos/dqlx"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEq(t *testing.T) {
	t.Run("eq", func(t *testing.T) {
		query, args, err := dql.Eq{
			"field1": "value1",
			"field2": "value2",
		}.ToDQL()

		require.NoError(t, err)
		require.Equal(t, []interface{}{"value1", "value2"}, args)
		require.Equal(t, "eq(field1,??) AND eq(field2,??)", query)
	})

	t.Run("eqFn", func(t *testing.T) {
		query, args, err := dql.EqFn("field1", "value1").ToDQL()
		require.NoError(t, err)
		require.Equal(t, args, []interface{}{"value1"})
		require.Equal(t, "eq(field1,??)", query)
	})
}

func TestLe(t *testing.T) {
	t.Run("le", func(t *testing.T) {
		query, args, err := dql.Le{
			"field1": "value1",
			"field2": "value2",
		}.ToDQL()

		require.NoError(t, err)
		require.Equal(t, []interface{}{"value1", "value2"}, args)
		require.Equal(t, "le(field1,??) AND le(field2,??)", query)
	})

	t.Run("leFn", func(t *testing.T) {
		query, args, err := dql.LeFn("field1", "value1").ToDQL()
		require.NoError(t, err)
		require.Equal(t, args, []interface{}{"value1"})
		require.Equal(t, "le(field1,??)", query)
	})
}

func TestLt(t *testing.T) {
	t.Run("lt", func(t *testing.T) {
		query, args, err := dql.Lt{
			"field1": "value1",
			"field2": "value2",
		}.ToDQL()

		require.NoError(t, err)
		require.Equal(t, []interface{}{"value1", "value2"}, args)
		require.Equal(t, "lt(field1,??) AND lt(field2,??)", query)
	})

	t.Run("ltFn", func(t *testing.T) {
		query, args, err := dql.LtFn("field1", "value1").ToDQL()
		require.NoError(t, err)
		require.Equal(t, args, []interface{}{"value1"})
		require.Equal(t, "lt(field1,??)", query)
	})
}

func TestGe(t *testing.T) {
	t.Run("ge", func(t *testing.T) {
		query, args, err := dql.Ge{
			"field1": "value1",
			"field2": "value2",
		}.ToDQL()

		require.NoError(t, err)
		require.Equal(t, []interface{}{"value1", "value2"}, args)
		require.Equal(t, "ge(field1,??) AND ge(field2,??)", query)
	})

	t.Run("geFn", func(t *testing.T) {
		query, args, err := dql.GeFn("field1", "value1").ToDQL()
		require.NoError(t, err)
		require.Equal(t, args, []interface{}{"value1"})
		require.Equal(t, "ge(field1,??)", query)
	})
}

func TestGt(t *testing.T) {
	t.Run("gt", func(t *testing.T) {
		query, args, err := dql.Gt{
			"field1": "value1",
			"field2": "value2",
		}.ToDQL()

		require.NoError(t, err)
		require.Equal(t, []interface{}{"value1", "value2"}, args)
		require.Equal(t, "gt(field1,??) AND gt(field2,??)", query)
	})

	t.Run("gtFn", func(t *testing.T) {
		query, args, err := dql.GtFn("field1", "value1").ToDQL()
		require.NoError(t, err)
		require.Equal(t, []interface{}{"value1"}, args)
		require.Equal(t, "gt(field1,??)", query)
	})
}

func TestHas(t *testing.T) {
	t.Run("hasFn", func(t *testing.T) {
		query, args, err := dql.HasFn("field1").ToDQL()
		require.NoError(t, err)
		require.Len(t, args, 0)
		require.Equal(t, "has(field1)", query)
	})
}

func TestType(t *testing.T) {
	t.Run("typeFn", func(t *testing.T) {
		query, args, err := dql.TypeFn("field1").ToDQL()
		require.NoError(t, err)
		require.Len(t, args, 0)
		require.Equal(t, "type(field1)", query)
	})
}

func TestAllOfTerms(t *testing.T) {
	t.Run("allOfTerms", func(t *testing.T) {
		query, args, err := dql.AllOfTerms{
			"field1": "value1",
			"field2": "value2",
		}.ToDQL()

		require.NoError(t, err)
		require.Equal(t, []interface{}{"value1", "value2"}, args)
		require.Equal(t, "allofterms(field1,??) AND allofterms(field2,??)", query)
	})

	t.Run("allOfTermsFn", func(t *testing.T) {
		query, args, err := dql.AllOfTermsFn("field1", "value1").ToDQL()
		require.NoError(t, err)
		require.Equal(t, []interface{}{"value1"}, args)
		require.Equal(t, "allofterms(field1,??)", query)
	})
}

func TestAnyOfTerms(t *testing.T) {
	t.Run("anyOfTerms", func(t *testing.T) {
		query, args, err := dql.AnyOfTerms{
			"field1": "value1",
			"field2": "value2",
		}.ToDQL()

		require.NoError(t, err)
		require.Equal(t, []interface{}{"value1", "value2"}, args)
		require.Equal(t, "anyofterms(field1,??) AND anyofterms(field2,??)", query)
	})

	t.Run("anyOfTermsFn", func(t *testing.T) {
		query, args, err := dql.AnyOfTermsFn("field1", "value1").ToDQL()
		require.NoError(t, err)
		require.Equal(t, []interface{}{"value1"}, args)
		require.Equal(t, "anyofterms(field1,??)", query)
	})
}

func TestRegexp(t *testing.T) {
	t.Run("regexp", func(t *testing.T) {
		query, args, err := dql.Regexp{
			"field1": "/^Steven Sp.*$/",
			"field2": "/^Steven Sp.*$/",
		}.ToDQL()

		require.NoError(t, err)
		require.Len(t, args, 0)
		require.Equal(t, "regexp(field1,/^Steven Sp.*$/) AND regexp(field2,/^Steven Sp.*$/)", query)
	})

	t.Run("regexpFn", func(t *testing.T) {
		query, args, err := dql.RegexpFn("field1", "/^Steven Sp.*$/").ToDQL()
		require.NoError(t, err)
		require.Len(t, args, 0)
		require.Equal(t, "regexp(field1,/^Steven Sp.*$/)", query)
	})
}

func TestMatch(t *testing.T) {
	t.Run("match", func(t *testing.T) {
		query, args, err := dql.Match{
			"field1": "value1",
			"field2": "value2",
		}.ToDQL()

		require.NoError(t, err)
		require.Equal(t, []interface{}{"value1", "value2"}, args)
		require.Equal(t, "match(field1,??) AND match(field2,??)", query)
	})

	t.Run("matchFn", func(t *testing.T) {
		query, args, err := dql.MatchFn("field1", "value1").ToDQL()
		require.NoError(t, err)
		require.Equal(t, []interface{}{"value1"}, args)
		require.Equal(t, "match(field1,??)", query)
	})
}

func TestAllOfText(t *testing.T) {
	t.Run("allOfText", func(t *testing.T) {
		query, args, err := dql.AllOfText{
			"field1": "value1",
			"field2": "value2",
		}.ToDQL()

		require.NoError(t, err)
		require.Equal(t, []interface{}{"value1", "value2"}, args)
		require.Equal(t, "alloftext(field1,??) AND alloftext(field2,??)", query)
	})

	t.Run("allOfTextFn", func(t *testing.T) {
		query, args, err := dql.AllOfTextFn("field1", "value1").ToDQL()
		require.NoError(t, err)
		require.Equal(t, []interface{}{"value1"}, args)
		require.Equal(t, "alloftext(field1,??)", query)
	})
}

func TestAnyOfText(t *testing.T) {
	t.Run("anyOfText", func(t *testing.T) {
		query, args, err := dql.AnyOfText{
			"field1": "value1",
			"field2": "value2",
		}.ToDQL()

		require.NoError(t, err)
		require.Equal(t, []interface{}{"value1", "value2"}, args)
		require.Equal(t, "anyoftext(field1,??) AND anyoftext(field2,??)", query)
	})

	t.Run("anyOfTextFn", func(t *testing.T) {
		query, args, err := dql.AnyOfTextFn("field1", "value1").ToDQL()
		require.NoError(t, err)
		require.Equal(t, []interface{}{"value1"}, args)
		require.Equal(t, "anyoftext(field1,??)", query)
	})
}

func TestExact(t *testing.T) {
	t.Run("exact", func(t *testing.T) {
		query, args, err := dql.Exact{
			"field1": "value1",
			"field2": "value2",
		}.ToDQL()

		require.NoError(t, err)
		require.Equal(t, []interface{}{"value1", "value2"}, args)
		require.Equal(t, "exact(field1,??) AND exact(field2,??)", query)
	})

	t.Run("exactFn", func(t *testing.T) {
		query, args, err := dql.ExactFn("field1", "value1").ToDQL()
		require.NoError(t, err)
		require.Equal(t, []interface{}{"value1"}, args)
		require.Equal(t, "exact(field1,??)", query)
	})
}

func TestTerm(t *testing.T) {
	t.Run("term", func(t *testing.T) {
		query, args, err := dql.Term{
			"field1": "value1",
			"field2": "value2",
		}.ToDQL()

		require.NoError(t, err)
		require.Equal(t, []interface{}{"value1", "value2"}, args)
		require.Equal(t, "term(field1,??) AND term(field2,??)", query)
	})

	t.Run("termFn", func(t *testing.T) {
		query, args, err := dql.TermFn("field1", "value1").ToDQL()
		require.NoError(t, err)
		require.Equal(t, []interface{}{"value1"}, args)
		require.Equal(t, "term(field1,??)", query)
	})
}

func TestFullText(t *testing.T) {
	t.Run("fulltext", func(t *testing.T) {
		query, args, err := dql.FullText{
			"field1": "value1",
			"field2": "value2",
		}.ToDQL()

		require.NoError(t, err)
		require.Equal(t, []interface{}{"value1", "value2"}, args)
		require.Equal(t, "fulltext(field1,??) AND fulltext(field2,??)", query)
	})

	t.Run("fulltextFn", func(t *testing.T) {
		query, args, err := dql.FullTextFn("field1", "value1").ToDQL()
		require.NoError(t, err)
		require.Equal(t, []interface{}{"value1"}, args)
		require.Equal(t, "fulltext(field1,??)", query)
	})
}

func TestSum(t *testing.T) {
	expression := dql.Sum("field")
	require.Equal(t, dql.Expr("sum(val(field))"), expression)
}

func TestAvg(t *testing.T) {
	expression := dql.Avg("field")
	require.Equal(t, dql.Expr("avg(val(field))"), expression)
}

func TestMin(t *testing.T) {
	expression := dql.Min("field")
	require.Equal(t, dql.Expr("min(val(field))"), expression)
}

func TestMax(t *testing.T) {
	expression := dql.Max("field")
	require.Equal(t, dql.Expr("max(val(field))"), expression)
}

func TestCount(t *testing.T) {
	expression := dql.Count("field")
	require.Equal(t, dql.Expr("count(field)"), expression)
}

func TestBetween(t *testing.T) {
	query, args, err := dql.Between("field1", 1, 10).ToDQL()
	require.NoError(t, err)
	require.Equal(t, []interface{}{1, 10}, args)
	require.Equal(t, "between(field1,??,??)", query)
}

func TestExpr(t *testing.T) {
	expression := dql.Expr("val(field1)")
	require.Equal(t, dql.RawExpression{Val: "val(field1)"}, expression)
}

func TestVal(t *testing.T) {
	query, args, err := dql.Val("field1").ToDQL()
	require.NoError(t, err)
	require.Len(t, args, 0)
	require.Equal(t, "val(field1)", query)
}

func TestUID(t *testing.T) {
	query, args, err := dql.UID("field1").ToDQL()
	require.NoError(t, err)
	require.Equal(t, []interface{}{"field1"}, args)
	require.Equal(t, "uid(??)", query)
}

func TestUIDIn(t *testing.T) {
	t.Run("uuid_in", func(t *testing.T) {
		query, args, err := dql.UIDIn{
			"field1": "value1",
			"field2": "value2",
		}.ToDQL()

		require.NoError(t, err)
		require.Equal(t, []interface{}{"value1", "value2"}, args)
		require.Equal(t, "uid_in(field1,??) AND uid_in(field2,??)", query)
	})

	t.Run("uuidFn", func(t *testing.T) {
		query, args, err := dql.UIDInFn("field1", "value").ToDQL()
		require.NoError(t, err)
		require.Equal(t, []interface{}{"value"}, args)
		require.Equal(t, "uid_in(field1,??)", query)
	})
}
