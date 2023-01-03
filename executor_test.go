package dqlx

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/dgraph-io/dgo/v210/protos/api"

	"github.com/stretchr/testify/require"
)

func UnmarshalTestHelper(name string, in interface{}, out interface{}) func(t *testing.T) {
	return func(t *testing.T) {
		var err error
		r := Response{
			Raw:         &api.Response{},
			dataKeyPath: name,
		}

		if r.dataKeyPath == "" {
			r.Raw.Json, err = json.Marshal(in)
		} else {
			r.Raw.Json, err = json.Marshal(map[string]interface{}{name: in})
		}
		require.NoError(t, err)

		err = r.Unmarshal(out)
		require.NoError(t, err)

		require.Equal(t, in, out)
	}
}

func TestUnmarshal(t *testing.T) {
	t.Helper()

	// Test cases to ensure the DQL schema types are properly unmarshalled
	// https://dgraph.io/docs/query-language/schema/#scalar-types
	t.Run("dateTime", func(t *testing.T) {
		t.Helper()

		var out map[string]time.Time
		in := map[string]time.Time{
			"now": time.Now().UTC(),
		}

		t.Run("unnamed", UnmarshalTestHelper("", &in, &out))
		t.Run("named", UnmarshalTestHelper("rootQuery", &in, &out))
	})
}
