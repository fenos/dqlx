package deku

import (
	"fmt"
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

func (schema *SchemaBuilder) ToString() string {
	writer := NewWriter()

	predicates := schema.PredicatesToString()
	types := schema.TypesToString()

	writer.AddLine(predicates)
	writer.Append(types)

	return writer.ToString()
}

func (schema *SchemaBuilder) PredicatesToString() string {
	writer := NewWriter()
	for _, predicate := range schema.predicates {
		writer.AddLine(predicate.ToString())
	}

	return writer.ToString()
}

func (schema *SchemaBuilder) TypesToString() string {
	writer := NewWriter()
	for _, dType := range schema.types {
		writer.AddLine(dType.ToString())
	}

	return writer.ToString()
}

func (schema *SchemaBuilder) HasType(name string) bool {
	for _, schemaType := range schema.types {
		if schemaType.Name == name {
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

func (schema *SchemaBuilder) Type(name string, options ...TypeBuilderOptionModifier) *TypeBuilder {
	if schema.HasType(name) {
		panic(fmt.Errorf("type '%s' already registered", name))
	}

	builder := &TypeBuilder{
		prefixFields: true,
		schema: schema,
		DGraphType: &DGraphType{
			Name:       name,
			predicates: nil,
		},
	}

	for _, modifier := range options {
		modifier(builder)
	}

	schema.types = append(schema.types, builder.DGraphType)

	return builder
}