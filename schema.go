package dqlx

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/dgraph-io/dgo/v210"
)

// SchemaBuilder used to compose a Dgraph schema in a fluent manner
type SchemaBuilder struct {
	Predicates []*DGraphPredicate
	Types      []*dGraphType

	client *dgo.Dgraph
}

// NewSchema creates a new SchemaBuilder
func NewSchema() *SchemaBuilder {
	return &SchemaBuilder{
		Predicates: nil,
		Types:      nil,
	}
}

// WithClient sets the dgraph client
func (schema *SchemaBuilder) WithClient(client *dgo.Dgraph) *SchemaBuilder {
	schema.client = client
	return schema
}

// ToDQL returns the schema in a DQL format
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

// PredicatesToString returns predicates to a DQL string
func (schema *SchemaBuilder) PredicatesToString() (string, error) {
	writer := bytes.Buffer{}

	registeredPredicates := map[string]*DGraphPredicate{}

	var predicates []string
	for _, predicate := range schema.Predicates {

		if registeredPredicate, ok := registeredPredicates[predicate.Name]; ok {
			if registeredPredicate.ScalarType != predicate.ScalarType {
				return "", fmt.Errorf("Predicate '%s' already registered with a different scalar type '%s'", predicate.Name, registeredPredicate.ScalarType)
			}
			continue
		}

		predicates = append(predicates, predicate.ToString())

		registeredPredicates[predicate.Name] = predicate
	}

	writer.WriteString(strings.Join(predicates, "\n"))
	return writer.String(), nil
}

// TypesToString returns type definitions to a DQL string
func (schema *SchemaBuilder) TypesToString() (string, error) {
	writer := bytes.Buffer{}

	types := make([]string, len(schema.Types))
	for index, dType := range schema.Types {
		dqlExpression, err := dType.ToString()

		if err != nil {
			return "", err
		}

		types[index] = dqlExpression
	}

	writer.WriteString(strings.Join(types, "\n"))

	return writer.String(), nil
}

// HasType determines if a type has been already registered
func (schema *SchemaBuilder) HasType(name string) bool {
	for _, schemaType := range schema.Types {
		if schemaType.name == name {
			return true
		}
	}
	return false
}

// HasPredicate determines if a predicate has been already registered
func (schema *SchemaBuilder) HasPredicate(name string) bool {
	for _, predicate := range schema.Predicates {
		if predicate.Name == name {
			return true
		}
	}
	return false
}

// TypeBuilderFn represents the closure function for the type builder
type TypeBuilderFn func(builder *TypeBuilder)

// Type registers a type definition to the schema
func (schema *SchemaBuilder) Type(name string, builderFn TypeBuilderFn, options ...TypeBuilderOptionModifier) *TypeBuilder {
	if schema.HasType(name) {
		panic(fmt.Errorf("type '%s' already registered", name))
	}

	builder := &TypeBuilder{
		prefixFields: true,
		dGraphType: &dGraphType{
			name:       name,
			predicates: nil,
		},
		schema: schema,
	}

	for _, modifier := range options {
		modifier(builder)
	}

	builderFn(builder)

	schema.Types = append(schema.Types, builder.dGraphType)

	return builder
}

// Predicate registers a predicate to the schema
func (schema *SchemaBuilder) Predicate(name string, scalar DGraphScalar) *PredicateBuilder {
	builder := &PredicateBuilder{
		predicate: &DGraphPredicate{
			Name:       name,
			ScalarType: scalar,
		},
	}

	schema.Predicates = append(schema.Predicates, builder.predicate)

	return builder
}

// PredicateString registers a predicate string type
func (schema *SchemaBuilder) PredicateString(name string) *PredicateStringBuilder {
	builder := &PredicateStringBuilder{
		PredicateBuilder: &PredicateBuilder{
			predicate: &DGraphPredicate{
				Name:       name,
				ScalarType: ScalarString,
			},
		},
	}

	schema.Predicates = append(schema.Predicates, builder.predicate)

	return builder
}

// PredicateDatetime registers a predicate datetime
func (schema *SchemaBuilder) PredicateDatetime(name string) *PredicateDateBuilder {
	builder := &PredicateDateBuilder{
		PredicateBuilder: &PredicateBuilder{
			predicate: &DGraphPredicate{
				Name:       name,
				ScalarType: ScalarDateTime,
			},
		},
	}

	schema.Predicates = append(schema.Predicates, builder.predicate)

	return builder
}

// Alter alters the schema with the current state. No drop operation will happen
// if DropAll is not explicitly set
func (schema *SchemaBuilder) Alter(ctx context.Context, options ...SchemaExecutorOptionFn) error {
	schemaExecutor := NewSchemaExecutor(schema.client)

	for _, option := range options {
		option(schemaExecutor)
	}

	return schemaExecutor.AlterSchema(ctx, schema)
}

// DropType drops a type
func (schema *SchemaBuilder) DropType(ctx context.Context, typeName string) error {
	schemaExecutor := NewSchemaExecutor(schema.client)

	return schemaExecutor.DropType(ctx, typeName)
}

// DropPredicate drops a predicate
func (schema *SchemaBuilder) DropPredicate(ctx context.Context, predicateName string) error {
	schemaExecutor := NewSchemaExecutor(schema.client)

	return schemaExecutor.DropPredicate(ctx, predicateName)
}
