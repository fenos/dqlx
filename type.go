package deku

import "fmt"

type DGraphType struct {
	Name       string
	predicates []*DGraphPredicate
}

func (builder *DGraphType) ToString() string {
	writer := NewWriter()
	writer.AddLine("type " + builder.Name + " {")

	for _, field := range builder.predicates {

		if field.reverse {
			writer.AddIndentedLine("<~" + field.name + ">")
		} else {
			writer.AddIndentedLine(field.name)
		}
	}

	writer.AddLine("}")

	return writer.ToString()
}

type TypeBuilder struct {
	*DGraphType
	schema *SchemaBuilder
	prefixFields bool
}

type TypeBuilderOptionModifier = func(*TypeBuilder)

func WithPrefix(usePrefix bool) TypeBuilderOptionModifier {
	return func(builder *TypeBuilder) {
		builder.prefixFields = usePrefix
	}
}

func (builder *TypeBuilder) String(name string) *PredicateStringBuilder {
	name = builder.normalizeName(name)

	field := &PredicateStringBuilder{
		PredicateBuilder: &PredicateBuilder{
			predicate: &DGraphPredicate{
				name:       name,
				tokenizers: nil,
				scalarType: ScalarString,
			},
			TypeBuilder: builder,
		},
	}

	builder.registerPredicate(name, field.PredicateBuilder.predicate)

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
			TypeBuilder: builder,
		},
	}

	builder.registerPredicate(name, field.PredicateBuilder.predicate)

	return field
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
		TypeBuilder: builder,
	}

	builder.registerPredicate(name, field.predicate)

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

func (builder *TypeBuilder) registerPredicate(name string, field *DGraphPredicate) {
	if builder.HasPredicate(name) {
		panic(fmt.Errorf("predicate '%s' already registered on type '%s'", name, builder.Name))
	}

	if !builder.schema.HasPredicate(name) {
		builder.schema.predicates = append(builder.schema.predicates, field)
	}

	builder.predicates = append(builder.predicates, field)
}
