package dqlx

import (
	"bytes"
	"fmt"
	"strings"
)

type edge struct {
	Name       string
	Alias      string
	Node       node
	RootFilter DQLizer
	Filters    []DQLizer
	Pagination Cursor
	Order      []DQLizer
	Group      []DQLizer
	Facets     []DQLizer
	Cascade    DQLizer
	IsRoot     bool
	IsVariable bool
}

func (edge edge) RelativeName() string {
	path := strings.Split(edge.Name, symbolEdgeTraversal)

	if len(path) == 0 {
		return ""
	}
	return path[len(path)-1]
}

func (edge edge) GetName() string {
	return edge.RelativeName()
}

func (edge edge) ToDQL() (query string, args []interface{}, err error) {
	writer := bytes.Buffer{}
	edgeName := edge.RelativeName()

	if edge.Alias != "" {
		writer.WriteString(fmt.Sprintf("%s as", edge.Alias))
	} else {
		if !(edge.IsRoot && edge.IsVariable) {
			if edgeName != "" {
				writer.WriteString(fmt.Sprintf("<%s>", edgeName))
			}
		}
	}

	if edge.IsVariable {
		writer.WriteString("var")
	}

	if edge.IsRoot {
		writer.WriteString("(")
		if edge.RootFilter != nil {
			writer.WriteString("func: ")
		}
	}

	if edge.RootFilter != nil {
		if err := addPart(edge.RootFilter, &writer, &args); err != nil {
			return "", nil, err
		}
	}

	if edge.IsRoot {
		// Pagination
		if edge.Pagination.WantsPagination() {
			writer.WriteString(",")

			if err := addPart(edge.Pagination, &writer, &args); err != nil {
				return "", nil, err
			}
		}

		// Order
		if len(edge.Order) > 0 {
			writer.WriteString(",")
			var statements []string
			if err := addStatement(edge.Order, &statements, &args); err != nil {
				return "", nil, err
			}

			writer.WriteString(strings.Join(statements, ","))
		}

		writer.WriteString(")")
	} else {

		// Pagination
		if edge.Pagination.WantsPagination() {
			writer.WriteString("(")
			if err := addPart(edge.Pagination, &writer, &args); err != nil {
				return "", nil, err
			}
			writer.WriteString(")")
		}

		// Order
		if len(edge.Order) > 0 {
			writer.WriteString("(")
			var statements []string
			if err := addStatement(edge.Order, &statements, &args); err != nil {
				return "", nil, err
			}

			writer.WriteString(strings.Join(statements, ","))
			writer.WriteString(")")
		}
	}

	if err := edge.addFacets(&writer, &args); err != nil {
		return "", nil, err
	}

	if err := edge.addFilters(&writer, &args); err != nil {
		return "", nil, err
	}

	if err := edge.addGroupBy(&writer, &args); err != nil {
		return "", nil, err
	}

	if err := edge.addCascade(&writer, &args); err != nil {
		return "", nil, err
	}

	writer.WriteString(" { ")

	if err := edge.addSelection(&writer, &args); err != nil {
		return "", nil, err
	}

	writer.WriteString(" }")

	return writer.String(), args, nil
}

func (edge edge) addFacets(writer *bytes.Buffer, args *[]interface{}) error {
	if len(edge.Facets) == 0 {
		return nil
	}

	writer.WriteString(" ")

	var statements []string
	if err := addStatement(edge.Facets, &statements, args); err != nil {
		return err
	}

	writer.WriteString(strings.Join(statements, " "))
	return nil
}

func (edge edge) addFilters(writer *bytes.Buffer, args *[]interface{}) error {
	if len(edge.Filters) == 0 {
		return nil
	}

	writer.WriteString(" @filter(")

	var statements []string
	if err := addStatement(edge.Filters, &statements, args); err != nil {
		return err
	}

	writer.WriteString(strings.Join(statements, " AND "))

	writer.WriteString(")")
	return nil
}

func (edge edge) addSelection(writer *bytes.Buffer, args *[]interface{}) error {
	return addPart(edge.Node, writer, args)
}

func (edge edge) addGroupBy(writer *bytes.Buffer, args *[]interface{}) error {
	if len(edge.Group) == 0 {
		return nil
	}

	writer.WriteString(" @groupby(")

	var statements []string
	if err := addStatement(edge.Group, &statements, args); err != nil {
		return err
	}

	writer.WriteString(strings.Join(statements, ","))

	writer.WriteString(")")
	return nil
}

func (edge edge) addCascade(writer *bytes.Buffer, args *[]interface{}) error {
	if edge.Cascade == nil {
		return nil
	}

	writer.WriteString(" ")

	var statements []string
	if err := addStatement([]DQLizer{edge.Cascade}, &statements, args); err != nil {
		return err
	}

	writer.WriteString(strings.Join(statements, " "))
	return nil
}

// EdgePath returns the abstract representation of an edge
// Example: dqlx.EdgePath("field1", "field2")
// Returns: "field1->field2"
func EdgePath(abstractPath ...string) string {
	return strings.Join(abstractPath, symbolEdgeTraversal)
}

// ParseEdge parses the abstract representation of an edge
// Example: dqlx.ParseEdge("field1->field2")
// Returns: []string{"field1", "field2"}
func ParseEdge(abstractPath string) []string {
	return strings.Split(abstractPath, symbolEdgeTraversal)
}
