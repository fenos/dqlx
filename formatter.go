package dqlx

import "strings"

// Minify minifies a dql query
func Minify(query string) string {
	parts := strings.Fields(query)
	return strings.Join(parts, " ")
}
