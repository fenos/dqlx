package dqlx

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// FuncType represents a function type
type FuncType string

var (
	eqFunc         FuncType = "eq"         // Done
	geFunc         FuncType = "ge"         // Done
	gtFunc         FuncType = "gt"         // Done
	leFunc         FuncType = "le"         // Done
	ltFunc         FuncType = "lt"         // Done
	hasFunc        FuncType = "has"        // Done
	typeFunc       FuncType = "type"       // Done
	alloftermsFunc FuncType = "allofterms" // Done
	anyoftermsFunc FuncType = "anyofterms" // Done
	regexpFunc     FuncType = "regexp"     // Done
	matchFunc      FuncType = "match"      // Done
	alloftextFunc  FuncType = "alloftext"  // Done
	anyoftextFunc  FuncType = "anyoftext"  // Done
	countFunc      FuncType = "count"      // Done
	exactFunc      FuncType = "exact"      // Done
	termFunc       FuncType = "term"       // Done
	fulltextFunc   FuncType = "fulltext"   // Done
	valFunc        FuncType = "val"        // Done
	sumFunc        FuncType = "sum"        // Done
	betweenFunc    FuncType = "between"    // Done
	uidFunc        FuncType = "uid"        // Done
	uidInFunc      FuncType = "uid_in"     // Done
)

// FilterFn represents a filter function
type FilterFn struct {
	DQLizer
}

type filterExpr struct {
	funcType FuncType
	value    interface{}
}

func (filter filterExpr) ToDQL() (query string, args Args, err error) {
	var placeholder string

	switch castValue := filter.value.(type) {
	case filterKV:
		return castValue.toDQL(filter.funcType)
	case filterExpr:
		innerFn, innerArgs, err := castValue.ToDQL()

		if err != nil {
			return "", nil, err
		}

		placeholder = innerFn
		args = append(args, innerArgs...)
	default:
		valuePlaceholder, valueArgs, err := parseValue(castValue)

		if err != nil {
			return "", nil, err
		}

		placeholder = valuePlaceholder
		args = append(args, valueArgs...)
	}

	return fmt.Sprintf("%s(%s)", filter.funcType, placeholder), args, nil
}

type filterKV map[string]interface{}

func (filter filterKV) toDQL(funcType FuncType) (query string, args Args, err error) {
	var expressions []string
	sortedKeys := getSortedKeys(filter)

	for _, key := range sortedKeys {
		value := filter[key]

		placeholder, fnArgs, err := parseValue(value)

		if err != nil {
			return "", nil, err
		}

		fnStatement := fmt.Sprintf("%s(%s,%s)", funcType, EscapePredicate(key), placeholder)

		expressions = append(expressions, fnStatement)
		args = append(args, fnArgs...)
	}

	return strings.Join(expressions, " AND "), args, nil
}

// Or represents a OR conjunction statement
// Example: dqlx.Or{ dql.Eq{} }
type Or conjunction

// And represents a AND conjunction statement
// Example: dqlx.And{ dql.Eq{} }
type And conjunction

type conjunction []DQLizer

func (connector conjunction) join(separator string) (query string, args Args, err error) {
	if len(connector) == 0 {
		return "", []interface{}{}, nil
	}

	var parts []string
	for _, part := range connector {
		dqlPart, partArgs, err := part.ToDQL()
		if err != nil {
			return "", nil, err
		}
		if dqlPart != "" {
			parts = append(parts, dqlPart)
			args = append(args, partArgs...)
		}
	}
	if len(parts) > 0 {
		query = fmt.Sprintf("(%s)", strings.Join(parts, separator))
	}
	return
}

// ToDQL returns the DQL statement for the 'or' expression
func (or Or) ToDQL() (query string, args Args, err error) {
	return conjunction(or).join(" OR ")
}

// ToDQL returns the DQL statement for the 'and' expression
func (and And) ToDQL() (query string, args Args, err error) {
	return conjunction(and).join(" AND ")
}

// Eq syntactic sugar for the Eq expression,
// Expression: eq(predicate, value)
// Example: dql.Eq{"predicate": "value"}
type Eq filterKV

// ToDQL returns the DQL statement for the 'eq' expression
func (eq Eq) ToDQL() (query string, args Args, err error) {
	return filterExpr{
		funcType: eqFunc,
		value:    filterKV(eq),
	}.ToDQL()
}

// EqFn represents the eq expression,
// Expression: eq(predicate, value)
func EqFn(predicate string, value interface{}, values ...interface{}) *FilterFn {
	valueList := []interface{}{value}
	valueList = append(valueList, values...)
	expression := Eq{}
	expression[predicate] = valueList
	return &FilterFn{expression}
}

// Le syntactic sugar for the Le expression,
// Expression: le(predicate, value)
// Example: dql.Le{"predicate": "value"}
type Le filterKV

// ToDQL returns the DQL statement for the 'le' expression
func (le Le) ToDQL() (query string, args Args, err error) {
	return filterExpr{
		funcType: leFunc,
		value:    filterKV(le),
	}.ToDQL()
}

// LeFn represents the le expression,
// Expression: le(predicate, value)
func LeFn(predicate string, value interface{}) *FilterFn {
	expression := Le{}
	expression[predicate] = value
	return &FilterFn{expression}
}

// Lt syntactic sugar for the Lt expression,
// Expression: lt(predicate, value)
// Example: dql.Lt{"predicate": "value"}
type Lt filterKV

// ToDQL returns the DQL statement for the 'lt' expression
func (lt Lt) ToDQL() (query string, args Args, err error) {
	return filterExpr{
		funcType: ltFunc,
		value:    filterKV(lt),
	}.ToDQL()
}

// LtFn represents the lt expression,
// Expression: lt(predicate, value)
func LtFn(predicate string, value interface{}) *FilterFn {
	expression := Lt{}
	expression[predicate] = value
	return &FilterFn{expression}
}

// Ge syntactic sugar for the Ge expression,
// Expression: ge(predicate, value)
// Example: dql.Ge{"predicate": "value"}
type Ge filterKV

// ToDQL returns the DQL statement for the 'ge' expression
func (ge Ge) ToDQL() (query string, args Args, err error) {
	return filterExpr{
		funcType: geFunc,
		value:    filterKV(ge),
	}.ToDQL()
}

// GeFn represents the ge expression,
// Expression: ge(predicate, value)
func GeFn(predicate string, value interface{}) *FilterFn {
	expression := Ge{}
	expression[predicate] = value
	return &FilterFn{expression}
}

// Gt syntactic sugar for the Gt expression,
// Expression: gt(predicate, value)
// Example: dql.Gt{"predicate": "value"}
type Gt filterKV

// ToDQL returns the DQL statement for the 'gt' expression
func (gt Gt) ToDQL() (query string, args Args, err error) {
	return filterExpr{
		funcType: gtFunc,
		value:    filterKV(gt),
	}.ToDQL()
}

// GtFn represents the gt expression,
// Expression: gt(predicate, value)
func GtFn(predicate string, value interface{}) *FilterFn {
	expression := Gt{}
	expression[predicate] = value
	return &FilterFn{expression}
}

// HasFn represents the has expression,
// Expression: has(predicate)
func HasFn(predicate string) *FilterFn {
	expression := filterExpr{
		funcType: hasFunc,
		value:    Predicate(predicate),
	}
	return &FilterFn{expression}
}

// Has alias of HasFn,
// Expression: has(predicate)
func Has(predicate string) DQLizer {
	return HasFn(predicate)
}

// TypeFn represents the type expression,
// Expression: type(predicate)
func TypeFn(predicate string) *FilterFn {
	expression := filterExpr{
		funcType: typeFunc,
		value:    Predicate(predicate),
	}
	return &FilterFn{expression}
}

// Type alias of TypeFn,
// Expression: type(predicate)
func Type(predicate string) DQLizer {
	return TypeFn(predicate)
}

// AllOfTerms syntactic sugar for the AllOfTerms expression,
// Expression: allofterms(predicate, value)
// Example: dql.AllOfTerms{"predicate": "value"}
type AllOfTerms filterKV

// ToDQL returns the DQL statement for the 'allofterms' expression
func (allOfTerms AllOfTerms) ToDQL() (query string, args Args, err error) {
	return filterExpr{
		funcType: alloftermsFunc,
		value:    filterKV(allOfTerms),
	}.ToDQL()
}

// AllOfTermsFn represents the allofterms expression,
// Expression: allofterms(predicate, value)
func AllOfTermsFn(predicate string, value interface{}) *FilterFn {
	expression := AllOfTerms{}
	expression[predicate] = value
	return &FilterFn{expression}
}

// AnyOfTerms syntactic sugar for the AnyOfTerms expression,
// Expression: anyOfTerms(predicate, value)
// Example: dql.AnyOfTerms{"predicate": "value"}
type AnyOfTerms filterKV

// ToDQL returns the DQL statement for the 'anyofterms' expression
func (anyOfTerms AnyOfTerms) ToDQL() (query string, args Args, err error) {
	return filterExpr{
		funcType: anyoftermsFunc,
		value:    filterKV(anyOfTerms),
	}.ToDQL()
}

// AnyOfTermsFn represents the anyofterms expression,
// Expression: anyofterms(predicate, value)
func AnyOfTermsFn(predicate string, value interface{}) *FilterFn {
	expression := AnyOfTerms{}
	expression[predicate] = value
	return &FilterFn{expression}
}

// Regexp syntactic sugar for the Regexp expression,
// Expression: regexp(predicate, /pattern/)
// Example: dql.Regexp{"predicate": /pattern/}
type Regexp map[string]string

// ToDQL returns the DQL statement for the 'regexp' expression
func (regexp Regexp) ToDQL() (query string, args Args, err error) {
	// Value can't be escaped
	rawRegexp := filterKV{}
	for key, value := range regexp {
		rawRegexp[key] = Expr(value)
	}

	return filterExpr{
		funcType: regexpFunc,
		value:    rawRegexp,
	}.ToDQL()
}

// RegexpFn represents the regexp expression,
// Expression: regexp(predicate, /pattern/)
func RegexpFn(predicate string, pattern string) *FilterFn {
	expression := Regexp{}
	expression[predicate] = pattern
	return &FilterFn{expression}
}

// Match syntactic sugar for the Match expression,
// Expression: match(predicate, value)
// Example: dql.Match{"predicate": "value"}
type Match filterKV

// ToDQL returns the DQL statement for the 'match' expression
func (match Match) ToDQL() (query string, args Args, err error) {
	return filterExpr{
		funcType: matchFunc,
		value:    filterKV(match),
	}.ToDQL()
}

// MatchFn represents the match expression,
// Expression: match(predicate, value)
func MatchFn(predicate string, pattern string) *FilterFn {
	expression := Match{}
	expression[predicate] = pattern
	return &FilterFn{expression}
}

// AllOfText syntactic sugar for the AllOfText expression,
// Expression: alloftext(predicate, value)
// Example: dql.AllOfText{"predicate": "value"}
type AllOfText filterKV

// ToDQL returns the DQL statement for the 'alloftext' expression
func (alloftext AllOfText) ToDQL() (query string, args Args, err error) {
	return filterExpr{
		funcType: alloftextFunc,
		value:    filterKV(alloftext),
	}.ToDQL()
}

// AllOfTextFn represents the match expression,
// Expression: alloftext(predicate, value)
func AllOfTextFn(predicate string, pattern string) *FilterFn {
	expression := AllOfText{}
	expression[predicate] = pattern
	return &FilterFn{expression}
}

// AnyOfText syntactic sugar for the AnyOfText expression,
// Expression: anyoftext(predicate, value)
// Example: dql.AnyOfText{"predicate": "value"}
type AnyOfText filterKV

// ToDQL returns the DQL statement for the 'anyoftext' expression
func (anyoftext AnyOfText) ToDQL() (query string, args Args, err error) {
	return filterExpr{
		funcType: anyoftextFunc,
		value:    filterKV(anyoftext),
	}.ToDQL()
}

// AnyOfTextFn represents the anyoftext expression,
// Expression: anyoftext(predicate, value)
func AnyOfTextFn(predicate string, pattern string) *FilterFn {
	expression := AnyOfText{}
	expression[predicate] = pattern
	return &FilterFn{expression}
}

// Exact syntactic sugar for the Exact expression,
// Expression: exact(predicate, value)
// Example: dql.Exact{"predicate": "value"}
type Exact filterKV

// ToDQL returns the DQL statement for the 'exact' expression
func (exact Exact) ToDQL() (query string, args Args, err error) {
	return filterExpr{
		funcType: exactFunc,
		value:    filterKV(exact),
	}.ToDQL()
}

// ExactFn represents the exact expression,
// Expression: exact(predicate, value)
func ExactFn(predicate string, pattern string) *FilterFn {
	expression := Exact{}
	expression[predicate] = pattern
	return &FilterFn{expression}
}

// Term syntactic sugar for the Term expression,
// Expression: term(predicate, value)
// Example: dql.Term{"predicate": "value"}
type Term filterKV

// ToDQL returns the DQL statement for the 'term' expression
func (term Term) ToDQL() (query string, args Args, err error) {
	return filterExpr{
		funcType: termFunc,
		value:    filterKV(term),
	}.ToDQL()
}

// TermFn represents the term expression,
// Expression: term(predicate, value)
func TermFn(predicate string, pattern string) *FilterFn {
	expression := Term{}
	expression[predicate] = pattern
	return &FilterFn{expression}
}

// FullText syntactic sugar for the FullText expression,
// Expression: fulltext(predicate, value)
// Example: dql.FullText{"predicate": "value"}
type FullText filterKV

// ToDQL returns the DQL statement for the 'fulltext' expression
func (fulltext FullText) ToDQL() (query string, args Args, err error) {
	return filterExpr{
		funcType: fulltextFunc,
		value:    filterKV(fulltext),
	}.ToDQL()
}

// FullTextFn represents the term expression,
// Expression: fulltext(predicate, value)
func FullTextFn(predicate string, pattern string) *FilterFn {
	expression := FullText{}
	expression[predicate] = pattern
	return &FilterFn{expression}
}

// Sum represent the 'sum' expression
// Expression: sum(val(predicate))
func Sum(predicate string) RawExpression {
	return Expr(string(sumFunc) + "(val(" + EscapePredicate(predicate) + "))")
}

// Avg represent the 'avg' expression
// Expression: avg(val(predicate))
func Avg(predicate string) RawExpression {
	return Expr("avg(val(" + EscapePredicate(predicate) + "))")
}

// Min represent the 'min' expression
// Expression: min(val(predicate))
func Min(predicate string) RawExpression {
	return Expr("min(val(" + EscapePredicate(predicate) + "))")
}

// Max represent the 'max' expression
// Expression: max(val(predicate))
func Max(predicate string) RawExpression {
	return Expr("max(val(" + EscapePredicate(predicate) + "))")
}

// Count represent the 'count' expression
// Expression: count(predicate)
func Count(predicate string) RawExpression {
	return Expr(string(countFunc) + "(" + EscapePredicate(predicate) + ")")
}

// P represent a predicate expression
// Expression: <predicate>
func P(predicate string) RawExpression {
	return Expr(EscapePredicate(predicate))
}

type between struct {
	predicate string
	from      interface{}
	to        interface{}
}

// Between represents the between expression,
// Expression: between(predicate, from, to)
func Between(predicate string, from interface{}, to interface{}) DQLizer {
	return between{
		predicate: predicate,
		from:      from,
		to:        to,
	}
}

// ToDQL returns the DQL statement for the 'between' expression
func (between between) ToDQL() (query string, args Args, err error) {
	placeholderFrom, argsFrom, err := parseValue(between.from)

	if err != nil {
		return "", nil, err
	}

	args = append(args, argsFrom...)

	placeholderTo, argsTo, err := parseValue(between.to)

	if err != nil {
		return "", nil, err
	}

	args = append(args, argsTo...)

	return fmt.Sprintf("%s(%s,%s,%s)", betweenFunc, EscapePredicate(between.predicate), placeholderFrom, placeholderTo), args, nil
}

// RawExpression represents a raw expression
// what you write what you get
type RawExpression struct {
	Val string
}

// ToDQL returns the DQL statement for a RawExpression expression, use with care
func (rawExpression RawExpression) ToDQL() (query string, args Args, err error) {
	return rawExpression.Val, nil, nil
}

// Expr returns a RawExpression
func Expr(value string) RawExpression {
	return RawExpression{value}
}

// Predicate returns an escaped predicate
func Predicate(predicate string) RawExpression {
	return Expr(EscapePredicate(predicate))
}

// Val returns val expression
// Expression: val(predicate)
func Val(ref string) filterExpr {
	return filterExpr{
		funcType: valFunc,
		value:    Predicate(ref),
	}
}

// UID returns uid expression
// Expression: uid(predicate)
func UID(value interface{}) filterExpr {
	return filterExpr{
		funcType: uidFunc,
		value:    value,
	}
}

// UIDFn returns uid expression
// Expression: uid(predicate)
func UIDFn(value interface{}) *FilterFn {
	return &FilterFn{UID(value)}
}

// UIDIn syntactic sugar for the UIDIn expression,
// Expression: uid_in(predicate, value)
// Example: dql.UIDIn{"predicate": "value"}
type UIDIn filterKV

// ToDQL returns the DQL statement for the 'uid_in' expression
func (uidin UIDIn) ToDQL() (query string, args Args, err error) {
	return filterExpr{
		funcType: uidInFunc,
		value:    filterKV(uidin),
	}.ToDQL()
}

// UIDInFn represents the uid_in expression,
// Expression: uid_in(predicate, value)
func UIDInFn(predicate string, value interface{}) *FilterFn {
	expression := UIDIn{}
	expression[predicate] = value
	return &FilterFn{expression}
}

// Cursor represents pagination parameters
type Cursor struct {
	First  int
	Offset int
	After  string
}

// WantsPagination determines if pagination is requested
func (p Cursor) WantsPagination() bool {
	return p.Offset != 0 || p.First != 0 || p.After != ""
}

// ToDQL returns the DQL statement for the 'pagination' expression
func (p Cursor) ToDQL() (query string, args Args, err error) {
	var paginationExpressions []string
	if p.First != 0 {
		paginationExpressions = append(paginationExpressions, "first:??")
		args = append(args, p.First)
	}

	if p.Offset != 0 {
		paginationExpressions = append(paginationExpressions, "offset:??")
		args = append(args, p.Offset)
	}

	if p.After != "" {
		paginationExpressions = append(paginationExpressions, "after:??")
		args = append(args, p.After)
	}

	return strings.Join(paginationExpressions, ","), args, nil
}

// OrderDirection represent an order direction
type OrderDirection string

var (
	OrderDirectionAsc  OrderDirection = "orderasc"
	OrderDirectionDesc OrderDirection = "orderdesc"
)

type orderBy struct {
	Direction OrderDirection
	Predicate interface{}
}

// OrderAsc returns an orderasc expression
func OrderAsc(predicate interface{}) DQLizer {
	return orderBy{
		Direction: OrderDirectionAsc,
		Predicate: predicate,
	}
}

// OrderDesc returns an orderdesc expression
func OrderDesc(predicate interface{}) DQLizer {
	return orderBy{
		Direction: OrderDirectionDesc,
		Predicate: predicate,
	}
}

// ToDQL returns the DQL statement for the 'order' expression
func (orderBy orderBy) ToDQL() (query string, args Args, err error) {
	predicate := orderBy.Predicate

	switch val := orderBy.Predicate.(type) {
	case filterExpr:
		if val.funcType != valFunc {
			return "", nil, fmt.Errorf("invalid function %s on order expression", val.funcType)
		}

		valDql, _, err := val.ToDQL()

		if err != nil {
			return "", nil, err
		}

		predicate = valDql
	case string:
		predicate = EscapePredicate(val)
	default:
		return "", nil, fmt.Errorf("order clause only accept Val() or string predicates given %v", val)
	}

	query = fmt.Sprintf("%s:%s", orderBy.Direction, predicate)
	return
}

type group struct {
	Predicate string
}

// ToDQL returns the DQL statement for the 'group' expression
func (group group) ToDQL() (query string, args Args, err error) {
	query = EscapePredicate(group.Predicate)
	return
}

// GroupBy returns an expression for grouping
func GroupBy(name string) DQLizer {
	return group{name}
}

type cascadeExpr struct {
	fields []string
}

// Cascade returns an expression for cascade directive
func Cascade(fields ...string) DQLizer {
	return cascadeExpr{fields}
}

// ToDQL returns the DQL statement for the 'cascade' expression
func (cascade cascadeExpr) ToDQL() (query string, args Args, err error) {
	if len(cascade.fields) > 0 {
		return fmt.Sprintf("@cascade(%s)", strings.Join(cascade.fields, ",")), nil, nil
	}

	return "@cascade", nil, nil
}

type facetExpr struct {
	Predicates []interface{}
}

// ToDQL returns the DQL statement for the 'facets' expression
func (facet facetExpr) ToDQL() (query string, args Args, err error) {
	var predicates []string
	for _, predicate := range facet.Predicates {

		switch predicateCast := predicate.(type) {
		case DQLizer:
			if err := addStatement([]DQLizer{predicateCast}, &predicates, &args); err != nil {
				return "", nil, err
			}
		case string:
			predicates = append(predicates, EscapePredicate(predicateCast))
		default:
			return "", nil, fmt.Errorf("facets accepts only DQlizers or string as value, given %v", predicateCast)
		}
	}

	writer := bytes.Buffer{}
	writer.WriteString("@facets")

	if len(predicates) > 0 {
		writer.WriteString(fmt.Sprintf("(%s)", strings.Join(predicates, ",")))
	}

	return writer.String(), args, nil
}

// Facets returns the expression for representing facets
func Facets(predicates ...interface{}) DQLizer {
	return facetExpr{Predicates: predicates}
}

func parseValue(value interface{}) (valuePlaceholder string, args Args, err error) {
	if isListType(value) {
		var listValue []interface{}

		listValue, err = toInterfaceSlice(value)

		if err != nil {
			return "", nil, err
		}

		placeholders := make([]string, len(listValue))
		for index, value := range listValue {
			placeholders[index] = symbolValuePlaceholder
			args = append(args, value)
		}

		valuePlaceholder = fmt.Sprintf("[%s]", strings.Join(placeholders, ","))
		return
	}

	switch castType := value.(type) {
	case RawExpression:
		valuePlaceholder = fmt.Sprintf("%s", castType.Val)
	default:
		args = append(args, value)
		valuePlaceholder = symbolValuePlaceholder
	}

	return
}

func getSortedVariables(exp map[int]interface{}) []int {
	sortedKeys := make([]int, 0, len(exp))
	for k := range exp {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Ints(sortedKeys)
	return sortedKeys
}

func getSortedKeys(exp map[string]interface{}) []string {
	sortedKeys := make([]string, 0, len(exp))
	for k := range exp {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	return sortedKeys
}

func toInterfaceSlice(slice interface{}) ([]interface{}, error) {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		return nil, errors.New("toInterfaceSlice given a non-slice type")
	}

	// Keep the distinction between nil and empty slice input
	if s.IsNil() {
		return nil, nil
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret, nil
}
