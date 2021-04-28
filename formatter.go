package deku

import "strings"

func Minify(query string) string {
	parts := strings.Fields(query)
	return strings.Join(parts, " ")
}

//func MinifySchema(schema string) string {
//
//}
