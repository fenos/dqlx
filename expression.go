package deku

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

type QueryFn string

var (
	eqOperator         QueryFn = "eq"  // Done
	geOperator         QueryFn = "ge"  // Done
	gtOperator         QueryFn = "gt"  // Done
	leOperator         QueryFn = "le"  // Done
	ltOperator         QueryFn = "lt"  // Done
	hasOperator        QueryFn = "has" // Done
	alloftermsOperator QueryFn = "allofterms"
	anyoftermsOperator QueryFn = "anyofterms"
	regexpOperator     QueryFn = "regexp"
	matchOperator      QueryFn = "match"
	alloftextOperator  QueryFn = "alloftext"
	anyoftextOperator  QueryFn = "anyoftext"
	countOperator      QueryFn = "count"
	exactOperator      QueryFn = "exact"
	termOperator       QueryFn = "term"
	fulltextOperator   QueryFn = "fulltext"
	valOperator        QueryFn = "val"
	sumOperator        QueryFn = "sum"
	betweenOperator    QueryFn = "between"
	uidOperator        QueryFn = "uid"
	uidInOperator      QueryFn = "uid_in"
)

type FilterFn = func() DQLizer

type QueryFunction struct {
	operator QueryFn
	field    string
	value    interface{}
}

var placeholderSymbol = "??"

func (queryFunction QueryFunction) ToDQL() (query string, args []interface{}, err error) {

	placeholder := placeholderSymbol

	if isListType(queryFunction.value) {
		listValue, err := toInterfaceSlice(queryFunction.value)

		if err != nil {
			return "", nil, err
		}

		placeholders := make([]string, len(listValue))
		for index, value := range listValue {
			placeholders[index] = placeholderSymbol
			args = append(args, value)
		}

		placeholder = fmt.Sprintf("[%s]", strings.Join(placeholders, ","))
	} else {
		args = append(args, queryFunction.value)
	}

	query = fmt.Sprintf("%s(%s,%s)", queryFunction.operator, queryFunction.field, placeholder)

	return query, args, nil
}

type connector []DQLizer

func (connector connector) join(separator string) (query string, args []interface{}, err error) {
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

type Or connector
type And connector

func (or Or) ToDQL() (query string, args []interface{}, err error) {
	return connector(or).join(" OR ")
}

func (and And) ToDQL() (query string, args []interface{}, err error) {
	return connector(and).join(" AND ")
}

type mapExpression map[string]interface{}

func (expression mapExpression) toDQL(operator QueryFn) (query string, args []interface{}, err error) {
	var expressions []string
	sortedKeys := getSortedKeys(expression)

	for _, key := range sortedKeys {
		value := expression[key]

		queryFn := QueryFunction{
			operator: operator,
			field:    key,
			value:    value,
		}

		fnStatement, fnArgs, err := queryFn.ToDQL()

		if err != nil {
			return "", nil, err
		}

		expressions = append(expressions, fnStatement)

		args = append(args, fnArgs...)
	}

	return strings.Join(expressions, ", "), args, nil
}

func Fields(fields string) []string {
	return strings.Fields(fields)
}

// Eq eq expression eq(field, value)
type Eq mapExpression

func (eq Eq) ToDQL() (query string, args []interface{}, err error) {
	return mapExpression(eq).toDQL(eqOperator)
}

func EqFn(field string, value interface{}) FilterFn {
	return func() DQLizer {
		expression := Eq{}
		expression[field] = value
		return expression
	}
}

// Le le expression le(field, value)
type Le mapExpression

func (le Le) ToDQL() (query string, args []interface{}, err error) {
	return mapExpression(le).toDQL(leOperator)
}

func LeFn(field string, value interface{}) FilterFn {
	return func() DQLizer {
		expression := Le{}
		expression[field] = value
		return expression
	}
}

// Lt lt expression lt(field, value)
type Lt mapExpression

func (lt Lt) ToDQL() (query string, args []interface{}, err error) {
	return mapExpression(lt).toDQL(ltOperator)
}

func LtFn(field string, value interface{}) FilterFn {
	return func() DQLizer {
		expression := Lt{}
		expression[field] = value
		return expression
	}
}

// Ge ge expression ge(field, value)
type Ge mapExpression

func (ge Ge) ToDQL() (query string, args []interface{}, err error) {
	return mapExpression(ge).toDQL(geOperator)
}

func GeFn(field string, value interface{}) FilterFn {
	return func() DQLizer {
		expression := Ge{}
		expression[field] = value
		return expression
	}
}

// Gt gt expression gt(field, value)
type Gt mapExpression

func (gt Gt) ToDQL() (query string, args []interface{}, err error) {
	return mapExpression(gt).toDQL(gtOperator)
}

func GtFn(field string, value interface{}) FilterFn {
	return func() DQLizer {
		expression := Gt{}
		expression[field] = value
		return expression
	}
}

// Has has expression has(field, value)
type Has mapExpression

func (has Has) ToDQL() (query string, args []interface{}, err error) {
	return mapExpression(has).toDQL(gtOperator)
}

func HasFn(field string, value interface{}) FilterFn {
	return func() DQLizer {
		expression := Has{}
		expression[field] = value
		return expression
	}
}

// AllOfTerms allofterms expression allofterms(field, value)
type AllOfTerms mapExpression

func (allOfTerms AllOfTerms) ToDQL() (query string, args []interface{}, err error) {
	return mapExpression(allOfTerms).toDQL(alloftermsOperator)
}

func AllOfTermsFn(field string, value interface{}) FilterFn {
	return func() DQLizer {
		expression := AllOfTerms{}
		expression[field] = value
		return expression
	}
}

// AnyOfTerms anyofterms expression anyofterms(field, value)
type AnyOfTerms mapExpression

func (anyOfTerms AnyOfTerms) ToDQL() (query string, args []interface{}, err error) {
	return mapExpression(anyOfTerms).toDQL(anyoftermsOperator)
}

func AnyOfTermsFn(field string, value interface{}) FilterFn {
	return func() DQLizer {
		expression := AnyOfTerms{}
		expression[field] = value
		return expression
	}
}

// Regexp regexp expression regexp(field, /pattern/)
type Regexp mapExpression

func (regexp Regexp) ToDQL() (query string, args []interface{}, err error) {
	return mapExpression(regexp).toDQL(regexpOperator)
}

func RegexpFn(field string, pattern string) FilterFn {
	return func() DQLizer {
		expression := Regexp{}
		expression[field] = pattern
		return expression
	}
}

// Match match expression match(field, /pattern/)
type Match mapExpression

func (match Match) ToDQL() (query string, args []interface{}, err error) {
	return mapExpression(match).toDQL(matchOperator)
}

func MatchFn(field string, pattern string) FilterFn {
	return func() DQLizer {
		expression := Match{}
		expression[field] = pattern
		return expression
	}
}

type Pagination struct {
	First  int
	Offset int
	After  string
}

func (p Pagination) WantsPagination() bool {
	return p.Offset != 0 || p.First != 0 || p.After != ""
}

func (p Pagination) ToDQL() (query string, args []interface{}, err error) {

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

//alloftermsOperator QueryFn = "allofterms" // DONE
//anyoftermsOperator QueryFn = "anyofterms" // DONE
//regexpOperator     QueryFn = "regexp" // DONE
//matchOperator      QueryFn = "match" // DONE
//alloftextOperator  QueryFn = "alloftext"
//anyoftextOperator  QueryFn = "anyoftext"
//countOperator      QueryFn = "count"
//exactOperator      QueryFn = "exact"
//termOperator       QueryFn = "term"
//fulltextOperator   QueryFn = "fulltext"
//valOperator        QueryFn = "val"
//sumOperator        QueryFn = "sum"
//betweenOperator    QueryFn = "between"
//uidOperator        QueryFn = "uid"
//uidInOperator      QueryFn = "uid_in"

func getSortedVariables(exp map[string]interface{}) []string {
	sortedKeys := make([]string, 0, len(exp))
	for k := range exp {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Slice(sortedKeys, func(i, j int) bool {
		sNumA := strings.Replace(sortedKeys[i], "$", "", 1)
		sNumB := strings.Replace(sortedKeys[j], "$", "", 1)

		numA, _ := strconv.Atoi(sNumA)
		numB, _ := strconv.Atoi(sNumB)
		return numA < numB
	})
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
