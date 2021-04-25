package deku

import (
	"bytes"
	"strings"
)

type edge struct {
	Name    string
	Filters []DQLizer
	IsRoot  bool
}

func (edge edge) RelativeName() string {
	path := strings.Split(edge.Name, "->")

	if len(path) == 0 {
		return ""
	}
	return path[len(path)-1]
}

func (edge edge) ToDQL() (query string, args []interface{}, err error) {
	writer := bytes.Buffer{}
	writer.WriteString(edge.RelativeName())

	if edge.IsRoot {
		writer.WriteString("(func: ")
	}

	if len(edge.Filters) > 0 {
		for _, filter := range edge.Filters {
			if err := addPart(filter, &writer, &args); err != nil {
				return "", nil, err
			}
		}
	}

	if edge.IsRoot {
		writer.WriteString(")")
	}

	return writer.String(), args, nil
}

func EdgePath(abstractPath ...string) string {
	return strings.Join(abstractPath, "->")
}
