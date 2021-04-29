package dqlx

import "strings"

func Minify(query string) string {
	parts := strings.Fields(query)
	return strings.Join(parts, " ")
}
