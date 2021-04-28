package deku_test

import (
	"github.com/fenos/deku"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGenerateTypes(t *testing.T) {
	schema := deku.NewSchema()

	schema.Type("User", func(user *deku.TypeBuilder) {
		user.String("name")
		user.String("surname")
		user.Bool("verified")
		user.DateTime("created_at")
		user.Int("age")
		user.Float("score")
		user.Password("password")
		user.Type("posts", "Post").List()
	})

	schema.Type("Post", func(post *deku.TypeBuilder) {
		post.String("title")
		post.String("description")
		post.Type("user", "User").Reverse()
	})

	err := deku.GenerateTypes(schema, deku.GeneratorOption{
		Path:        "C:\\Users\\fabri\\go\\src\\deku\\t.go",
		PackageName: "deku",
	})

	require.NoError(t, err)
}
