package deku

import (
	"strings"
)

type DGraphScalar string

var (
	ScalarDefault  DGraphScalar = "string"
	ScalarInt      DGraphScalar = "int"
	ScalarFloat    DGraphScalar = "float"
	ScalarString   DGraphScalar = "string"
	ScalarBool     DGraphScalar = "bool"
	ScalarDateTime DGraphScalar = "datetime"
	ScalarGeo      DGraphScalar = "geo"
	ScalarPassword DGraphScalar = "password"
	ScalarUID      DGraphScalar = "uid"
)

type StringIndexTokenizer string
type DateIndexTokenizer string

var (
	TokenixerExact    StringIndexTokenizer = "exact"
	TokenixerHash     StringIndexTokenizer = "hash"
	TokenixerTerm     StringIndexTokenizer = "term"
	TokenixerFulltext StringIndexTokenizer = "fulltext"
	TokenixerTrigram  StringIndexTokenizer = "trigram"

	TokenizerYear  DateIndexTokenizer = "year"
	TokenizerMonth DateIndexTokenizer = "month"
	TokenizerDay   DateIndexTokenizer = "day"
	TokenizerHour  DateIndexTokenizer = "hour"
)

type DGraphPredicate struct {
	Name       string
	Index      bool
	Upsert     bool
	List       bool
	Lang       bool
	Reverse    bool
	Tokenizers []string
	ScalarType DGraphScalar
}

func (builder *DGraphPredicate) ToString() string {
	predicateString := builder.Name + ": "

	scalarType := builder.ScalarType

	if !isKnownScalarType(scalarType) {
		scalarType = ScalarUID
	}

	if builder.List {
		predicateString += "[" + string(scalarType) + "] "
	} else {
		predicateString += string(scalarType) + " "
	}

	if builder.Index {
		if len(builder.Tokenizers) > 0 {
			predicateString += "@index(" + strings.Join(builder.Tokenizers, ",") + ") "
		} else {
			predicateString += "@index() "
		}
	}

	if builder.Upsert {
		predicateString += "@upsert "
	}

	if builder.Lang {
		predicateString += "@lang "
	}

	if builder.Reverse {
		predicateString += "@reverse "
	}

	predicateString += "."

	return predicateString
}

type PredicateBuilder struct {
	predicate *DGraphPredicate
}

func (builder *PredicateBuilder) Index() *PredicateBuilder {
	builder.predicate.Index = true
	return builder
}

func (builder *PredicateBuilder) Upsert() *PredicateBuilder {
	builder.predicate.Upsert = true
	return builder
}

func (builder *PredicateBuilder) List() *PredicateBuilder {
	builder.predicate.List = true
	return builder
}

func (builder *PredicateBuilder) Lang() *PredicateBuilder {
	builder.predicate.Lang = true
	return builder
}

func (builder *PredicateBuilder) Reverse() *PredicateBuilder {
	builder.predicate.Reverse = true
	return builder
}

type PredicateStringBuilder struct {
	*PredicateBuilder
}

func (builder *PredicateStringBuilder) HasIndex(tokenizer StringIndexTokenizer) bool {
	for _, token := range builder.predicate.Tokenizers {
		if string(tokenizer) == token {
			return true
		}
	}

	return false
}

func (builder *PredicateStringBuilder) Index(tokenizer StringIndexTokenizer) *PredicateStringBuilder {
	builder.predicate.Index = true

	if !builder.HasIndex(tokenizer) {
		builder.predicate.Tokenizers = append(builder.predicate.Tokenizers, string(tokenizer))
	}
	return builder
}

func (builder *PredicateStringBuilder) IndexExact() *PredicateStringBuilder {
	return builder.Index(TokenixerExact)
}

func (builder *PredicateStringBuilder) IndexHash() *PredicateStringBuilder {
	return builder.Index(TokenixerHash)
}

func (builder *PredicateStringBuilder) IndexTerm() *PredicateStringBuilder {
	return builder.Index(TokenixerTerm)
}

func (builder PredicateStringBuilder) IndexFulltext() *PredicateStringBuilder {
	return builder.Index(TokenixerFulltext)
}

func (builder *PredicateStringBuilder) IndexTrigram() *PredicateStringBuilder {
	return builder.Index(TokenixerTrigram)
}

type PredicateDateBuilder struct {
	*PredicateBuilder
}

func (builder *PredicateDateBuilder) HasIndex(tokenizer DateIndexTokenizer) bool {
	for _, token := range builder.predicate.Tokenizers {
		if string(tokenizer) == token {
			return true
		}
	}

	return false
}

func (builder *PredicateDateBuilder) Index(tokenizer DateIndexTokenizer) *PredicateDateBuilder {
	builder.predicate.Index = true

	if !builder.HasIndex(tokenizer) {
		builder.predicate.Tokenizers = append(builder.predicate.Tokenizers, string(tokenizer))
	}
	return builder
}

func (builder *PredicateDateBuilder) IndexYear() *PredicateDateBuilder {
	return builder.Index(TokenizerYear)
}

func (builder *PredicateDateBuilder) IndexMonth() *PredicateDateBuilder {
	return builder.Index(TokenizerMonth)
}

func (builder *PredicateDateBuilder) IndexDay() *PredicateDateBuilder {
	return builder.Index(TokenizerDay)
}

func (builder *PredicateDateBuilder) IndexHour() *PredicateDateBuilder {
	return builder.Index(TokenizerHour)
}

func isKnownScalarType(value DGraphScalar) bool {
	switch value {
	case
		ScalarPassword,
		ScalarUID,
		ScalarGeo,
		ScalarFloat,
		ScalarDateTime,
		ScalarBool,
		ScalarInt,
		ScalarDefault:
		return true
	}
	return false
}
