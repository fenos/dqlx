package deku

import (
	"bytes"
)

type queryGrammar struct {
	RootEdge  edge
	Selection selectionSet
	Filters   []DQLizer
	Variables []QueryBuilder
}

func (grammar queryGrammar) Name() string {
	return grammar.RootEdge.Name
}

func (grammar queryGrammar) ToDQL() (query string, args []interface{}, err error) {
	writer := bytes.Buffer{}

	args = []interface{}{}

	if err := grammar.addEdge(&writer, &args); err != nil {
		return "", nil, err
	}

	if err := grammar.addFilters(&writer, &args); err != nil {
		return "", nil, err
	}

	writer.WriteString(" { ")

	if err := grammar.addSelection(&writer, &args); err != nil {
		return "", nil, err
	}

	writer.WriteString(" }")

	return writer.String(), args, nil
}

func (grammar queryGrammar) addEdge(writer *bytes.Buffer, args *[]interface{}) error {
	return addPart(grammar.RootEdge, writer, args)
}

func (grammar queryGrammar) addFilters(writer *bytes.Buffer, args *[]interface{}) error {
	if len(grammar.Filters) == 0 {
		return nil
	}

	writer.WriteString(" @filter(")
	for _, filter := range grammar.Filters {
		if err := addPart(filter, writer, args); err != nil {
			return err
		}
	}

	writer.WriteString(")")
	return nil
}

func (grammar queryGrammar) addSelection(writer *bytes.Buffer, args *[]interface{}) error {
	return addPart(grammar.Selection, writer, args)
}

func addPart(part DQLizer, writer *bytes.Buffer, args *[]interface{}) error {
	statement, statementArgs, err := part.ToDQL()
	*args = append(*args, statementArgs...)

	if err != nil {
		return err
	}

	writer.WriteString(statement)

	return nil
}
