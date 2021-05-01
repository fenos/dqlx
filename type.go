package dqlx

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

func NewTypeBuilder(predicate string, options ...TypeBuilderOptionModifier) *TypeBuilder {
	builder := &TypeBuilder{
		DGraphType: &DGraphType{
			name:       predicate,
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

func (builder *TypeBuilder) String(predicate string) *PredicateStringBuilder {
	predicate = builder.normalizeName(predicate)

	field := &PredicateStringBuilder{
		PredicateBuilder: &PredicateBuilder{
			predicate: &DGraphPredicate{
				Name:       predicate,
				Tokenizers: nil,
				ScalarType: ScalarString,
			},
		},
	}

	builder.registerPredicate(field.PredicateBuilder.predicate)

	return field
}

func (builder *TypeBuilder) DateTime(predicate string) *PredicateDateBuilder {
	predicate = builder.normalizeName(predicate)

	field := &PredicateDateBuilder{
		PredicateBuilder: &PredicateBuilder{
			predicate: &DGraphPredicate{
				Name:       predicate,
				Tokenizers: nil,
				ScalarType: ScalarDateTime,
			},
		},
	}

	builder.registerPredicate(field.PredicateBuilder.predicate)

	return field
}

func (builder *TypeBuilder) Type(predicate string, kind string) *PredicateBuilder {
	castKind := DGraphScalar(kind)
	return builder.field(predicate, castKind)
}

func (builder *TypeBuilder) Predicate(predicate *PredicateBuilder) {
	builder.registerPredicate(predicate.predicate)
}

func (builder *TypeBuilder) UID(predicate string) *PredicateBuilder {
	return builder.field(predicate, ScalarUID)
}

func (builder *TypeBuilder) Int(predicate string) *PredicateBuilder {
	return builder.field(predicate, ScalarInt)
}

func (builder *TypeBuilder) Float(predicate string) *PredicateBuilder {
	return builder.field(predicate, ScalarFloat)
}

func (builder *TypeBuilder) Bool(predicate string) *PredicateBuilder {
	return builder.field(predicate, ScalarBool)
}

func (builder *TypeBuilder) Geo(predicate string) *PredicateBuilder {
	return builder.field(predicate, ScalarGeo)
}

func (builder *TypeBuilder) Password(predicate string) *PredicateBuilder {
	return builder.field(predicate, ScalarPassword)
}

func (builder *TypeBuilder) field(predicate string, scalar DGraphScalar) *PredicateBuilder {
	predicate = builder.normalizeName(predicate)

	field := &PredicateBuilder{
		predicate: &DGraphPredicate{
			Name:       predicate,
			Tokenizers: nil,
			ScalarType: scalar,
		},
	}

	builder.registerPredicate(field.predicate)

	return field
}

func (builder *TypeBuilder) HasPredicate(predicateName string) bool {
	for _, predicate := range builder.predicates {
		if predicate.Name == predicateName {
			return true
		}
	}
	return false
}

func (builder *TypeBuilder) normalizeName(predicate string) string {
	if builder.prefixFields {
		return builder.name + "." + predicate
	}

	return predicate
}

func (builder *TypeBuilder) registerPredicate(predicate *DGraphPredicate) {
	if builder.schema != nil {
		builder.schema.Predicates = append(builder.schema.Predicates, predicate)
	}

	builder.predicates = append(builder.predicates, predicate)
}

type DGraphType struct {
	name       string
	predicates []*DGraphPredicate
}

func (builder *DGraphType) ToString() (string, error) {
	writer := bytes.Buffer{}
	writer.WriteString("type " + builder.name + " { ")

	// make sure duplicate QueryFields are not allowed
	registeredPredicates := map[string]bool{}
	fields := make([]string, 0, len(builder.predicates))

	for _, field := range builder.predicates {

		if _, ok := registeredPredicates[field.Name]; ok {
			return "", fmt.Errorf("predicate '%s' already registered on type '%s'", field.Name, builder.name)
		}

		predicate := field.Name
		if field.Reverse {
			if isKnownScalarType(field.ScalarType) {
				return "", fmt.Errorf("attempted to use a reverse field on a scalar value on field '%s'", field.Name)
			}

			predicate = fmt.Sprintf("<~%s>", field.Name)
		}

		if !isKnownScalarType(field.ScalarType) {
			if field.List {
				predicate += fmt.Sprintf(": [%s]", field.ScalarType)
			} else {
				predicate += fmt.Sprintf(": %s", field.ScalarType)
			}
		}

		fields = append(fields, predicate)

		registeredPredicates[field.Name] = true
	}

	writer.WriteString(strings.Join(fields, "\n"))
	writer.WriteString("\n}")

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
