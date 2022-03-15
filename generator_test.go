package dqlx_test

import (
	"testing"

	dql "github.com/getplexy/dqlx"
	"github.com/stretchr/testify/require"
)

func TestGenerateTypes(t *testing.T) {
	t.Skipf("Working Progress")
	schema := dql.NewSchema()

	schema.Type("User", func(user *dql.TypeBuilder) {
		user.String("name")
		user.String("surname")
		user.DateTime("birthday")
		user.Password("password")
	})

	schema.Type("Tag", func(tag *dql.TypeBuilder) {
		tag.String("name")
		tag.Type("posts", "Post").Reverse()
	})

	schema.Type("Post", func(post *dql.TypeBuilder) {
		post.String("title").Lang()
		post.String("content")
		post.Bool("published")
		post.DateTime("created_at")
		post.Type("tags", "Tag").Reverse().List()
		post.Int("views")
	})

	err := dql.GenerateTypes(schema, dql.GeneratorOption{
		Path:        "C:\\Users\\fabri\\go\\src\\dqlx\\t.go",
		PackageName: "dqlx",
	})

	require.NoError(t, err)
}
