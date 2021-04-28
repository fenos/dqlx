package deku

import (
	"bytes"
	"fmt"
	"strings"
)

type SchemaBuilder struct {
	predicates []*DGraphPredicate
	types      []*DGraphType
}

func NewSchema() *SchemaBuilder {
	return &SchemaBuilder{
		predicates: nil,
		types:      nil,
	}
}

func (schema *SchemaBuilder) ToDQL() (string, error) {
	writer := bytes.Buffer{}

	predicates, err := schema.PredicatesToString()

	if err != nil {
		return "", err
	}

	types, err := schema.TypesToString()

	if err != nil {
		return "", err
	}

	writer.WriteString(predicates)

	if types != "" {
		writer.WriteString("\n\n")
		writer.WriteString(types)
	}

	return writer.String(), nil
}

func (schema *SchemaBuilder) PredicatesToString() (string, error) {
	writer := bytes.Buffer{}

	registeredPredicates := map[string]*DGraphPredicate{}

	var predicates []string
	for _, predicate := range schema.predicates {

		if registeredPredicate, ok := registeredPredicates[predicate.name]; ok {
			if registeredPredicate.scalarType != predicate.scalarType {
				return "", fmt.Errorf("Predicate '%s' already registered with a different scalar type '%s'", predicate.name, registeredPredicate.scalarType)
			}
			continue
		}

		predicates = append(predicates, predicate.ToString())

		registeredPredicates[predicate.name] = predicate
	}

	writer.WriteString(strings.Join(predicates, "\n"))
	return writer.String(), nil
}

func (schema *SchemaBuilder) TypesToString() (string, error) {
	writer := bytes.Buffer{}

	types := make([]string, len(schema.types))
	for index, dType := range schema.types {
		dqlExpression, err := dType.ToString()

		if err != nil {
			return "", err
		}

		types[index] = dqlExpression
	}

	writer.WriteString(strings.Join(types, "\n"))

	return writer.String(), nil
}

func (schema *SchemaBuilder) HasType(name string) bool {
	for _, schemaType := range schema.types {
		if schemaType.name == name {
			return true
		}
	}
	return false
}

func (schema *SchemaBuilder) HasPredicate(name string) bool {
	for _, predicate := range schema.predicates {
		if predicate.name == name {
			return true
		}
	}
	return false
}

type TypeBuilderFn func(builder *TypeBuilder)

func (schema *SchemaBuilder) Type(name string, builderFn TypeBuilderFn, options ...TypeBuilderOptionModifier) *TypeBuilder {
	if schema.HasType(name) {
		panic(fmt.Errorf("type '%s' already registered", name))
	}

	builder := &TypeBuilder{
		prefixFields: true,
		DGraphType: &DGraphType{
			name:       name,
			predicates: nil,
		},
		schema: schema,
	}

	for _, modifier := range options {
		modifier(builder)
	}

	builderFn(builder)

	schema.types = append(schema.types, builder.DGraphType)

	return builder
}

func (schema *SchemaBuilder) Predicate(name string, scalar DGraphScalar) *PredicateBuilder {
	builder := &PredicateBuilder{
		predicate: &DGraphPredicate{
			name:       name,
			scalarType: scalar,
		},
	}

	schema.predicates = append(schema.predicates, builder.predicate)

	return builder
}
