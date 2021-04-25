package deku

import (
	"bytes"
	"fmt"
	"strings"
)

type TypeBuilder struct {
	*DGraphType
	prefixFields bool
	schema       *SchemaBuilder
}

func NewTypeBuilder(name string, options ...TypeBuilderOptionModifier) *TypeBuilder {
	builder := &TypeBuilder{
		DGraphType: &DGraphType{
			Name:       name,
			predicates: nil,
		},
		prefixFields: true,
	}

	for _, option := range options {
		option(builder)
	}

	return builder
}

type TypeBuilderOptionModifier = func(*TypeBuilder)

func WithTypePrefix(usePrefix bool) TypeBuilderOptionModifier {
	return func(builder *TypeBuilder) {
		builder.prefixFields = usePrefix
	}
}

func (builder *TypeBuilder) String(name string) PredicateStringBuilder {
	name = builder.normalizeName(name)

	field := PredicateStringBuilder{
		PredicateBuilder: PredicateBuilder{
			predicate: &DGraphPredicate{
				name:       name,
				tokenizers: nil,
				scalarType: ScalarString,
			},
		},
	}

	builder.registerPredicate(field.PredicateBuilder.predicate)

	return field
}

func (builder *TypeBuilder) DateTime(name string) *PredicateDateBuilder {
	name = builder.normalizeName(name)

	field := &PredicateDateBuilder{
		PredicateBuilder: &PredicateBuilder{
			predicate: &DGraphPredicate{
				name:       name,
				tokenizers: nil,
				scalarType: ScalarDateTime,
			},
		},
	}

	builder.registerPredicate(field.PredicateBuilder.predicate)

	return field
}

func (builder *TypeBuilder) Type(kind string, name string) *PredicateBuilder {
	castKind := DGraphScalar(kind)
	return builder.field(name, castKind)
}

func (builder *TypeBuilder) UID(name string) *PredicateBuilder {
	return builder.field(name, ScalarUID)
}

func (builder *TypeBuilder) Int(name string) *PredicateBuilder {
	return builder.field(name, ScalarInt)
}

func (builder *TypeBuilder) Float(name string) *PredicateBuilder {
	return builder.field(name, ScalarFloat)
}

func (builder *TypeBuilder) Bool(name string) *PredicateBuilder {
	return builder.field(name, ScalarBool)
}

func (builder *TypeBuilder) Geo(name string) *PredicateBuilder {
	return builder.field(name, ScalarGeo)
}

func (builder *TypeBuilder) Password(name string) *PredicateBuilder {
	return builder.field(name, ScalarPassword)
}

func (builder *TypeBuilder) field(name string, scalar DGraphScalar) *PredicateBuilder {
	name = builder.normalizeName(name)

	field := &PredicateBuilder{
		predicate: &DGraphPredicate{
			name:       name,
			tokenizers: nil,
			scalarType: scalar,
		},
	}

	builder.registerPredicate(field.predicate)

	return field
}

func (builder *TypeBuilder) HasPredicate(name string) bool {
	for _, predicate := range builder.predicates {
		if predicate.name == name {
			return true
		}
	}
	return false
}

func (builder *TypeBuilder) normalizeName(name string) string {
	if builder.prefixFields {
		return builder.Name + "." + name
	}

	return name
}

func (builder *TypeBuilder) registerPredicate(predicate *DGraphPredicate) {
	if builder.schema != nil {
		builder.schema.predicates = append(builder.schema.predicates, predicate)
	}

	builder.predicates = append(builder.predicates, predicate)
}

type DGraphType struct {
	Name       string
	predicates []*DGraphPredicate
}

func (builder *DGraphType) ToString() (string, error) {
	writer := bytes.Buffer{}
	writer.WriteString("type " + builder.Name + " {")

	// make sure duplicate QueryFields are not allowed
	registeredPredicates := map[string]bool{}

	for _, field := range builder.predicates {

		if _, ok := registeredPredicates[field.name]; ok {
			return "", fmt.Errorf("field '%s' already registered on type '%s'", field.name, builder.Name)
		}

		writer.WriteString(" " + field.name + "")

		registeredPredicates[field.name] = true
	}

	writer.WriteString(" }")

	return writer.String(), nil
}

func (builder *DGraphType) PredicatesToString() string {
	writer := bytes.Buffer{}

	var parts []string
	for _, field := range builder.predicates {
		parts = append(parts, field.ToString())
	}

	writer.WriteString(strings.Join(parts, " "))

	return writer.String()
}
