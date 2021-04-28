package deku_test

import (
	dql "github.com/fenos/deku"
	"github.com/stretchr/testify/require"
	"testing"
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
		Path:        "C:\\Users\\fabri\\go\\src\\deku\\t.go",
		PackageName: "deku",
	})

	require.NoError(t, err)
}
