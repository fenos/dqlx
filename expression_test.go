package deku_test

import (
	dql "github.com/fenos/deku"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEqFn(t *testing.T) {
	query, args, err := dql.Eq{
		"field1": "value1",
		"field2": "value2",
	}.ToDQL()

	require.NoError(t, err)
	require.Equal(t, args, []interface{}{"value1", "value2"})
	require.Equal(t, "eq(field1,??) AND eq(field2,??)", query)
}
