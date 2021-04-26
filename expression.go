package deku

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

type FuncType string

var (
	eqFunc         FuncType = "eq"         // Done
	geFunc         FuncType = "ge"         // Done
	gtFunc         FuncType = "gt"         // Done
	leFunc         FuncType = "le"         // Done
	ltFunc         FuncType = "lt"         // Done
	hasFunc        FuncType = "has"        // Done
	alloftermsFunc FuncType = "allofterms" // Done
	anyoftermsFunc FuncType = "anyofterms" // Done
	regexpFunc     FuncType = "regexp"     // Done
	matchFunc      FuncType = "match"      // Done
	alloftextFunc  FuncType = "alloftext"
	anyoftextFunc  FuncType = "anyoftext"
	countFunc      FuncType = "count"
	exactFunc      FuncType = "exact"
	termFunc       FuncType = "term"
	fulltextFunc   FuncType = "fulltext"
	valFunc        FuncType = "val"
	sumFunc        FuncType = "sum"
	betweenFunc    FuncType = "between"
	uidFunc        FuncType = "uid" // Done
	uid            FuncType = "uid_in"
)

type FilterFn = func() DQLizer

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

type Filter struct {
	value DQLizer
}

type filterValue struct {
	funcType FuncType
	value    interface{}
}

func (filter filterValue) ToDQL() (query string, args []interface{}, err error) {
	var placeholder string

	if innerFilterValue, ok := filter.value.(filterValue); ok {
		innerFn, innerArgs, err := innerFilterValue.ToDQL()

		if err != nil {
			return "", nil, err
		}

		placeholder = innerFn
		args = append(args, innerArgs...)
	} else {
		valuePlaceholder, valueArgs, err := parseValue(filter.value)

		if err != nil {
			return "", nil, err
		}

		placeholder = valuePlaceholder
		args = append(args, valueArgs)
	}

	return fmt.Sprintf("%s(%s)", filter.funcType, placeholder), args, nil
}

type filterKV map[string]interface{}

func (filter filterKV) ToDQL(funcType FuncType) (query string, args []interface{}, err error) {
	var expressions []string
	sortedKeys := getSortedKeys(filter)

	for _, key := range sortedKeys {
		value := filter[key]

		placeholder, fnArgs, err := parseValue(value)

		if err != nil {
			return "", nil, err
		}

		fnStatement := fmt.Sprintf("%s(%s,%s)", funcType, key, placeholder)

		expressions = append(expressions, fnStatement)
		args = append(args, fnArgs...)
	}

	return strings.Join(expressions, ", "), args, nil
}

// Eq eq expression eq(field, value)
type Eq filterKV

func (eq Eq) ToDQL() (query string, args []interface{}, err error) {
	return filterKV(eq).ToDQL(eqFunc)
}

func EqFn(field string, value interface{}) FilterFn {
	return func() DQLizer {
		expression := Eq{}
		expression[field] = value
		return expression
	}
}

// Le le expression le(field, value)
type Le filterKV

func (le Le) ToDQL() (query string, args []interface{}, err error) {
	return filterKV(le).ToDQL(leFunc)
}

func LeFn(field string, value interface{}) FilterFn {
	return func() DQLizer {
		expression := Le{}
		expression[field] = value
		return expression
	}
}

// Lt lt expression lt(field, value)
type Lt filterKV

func (lt Lt) ToDQL() (query string, args []interface{}, err error) {
	return filterKV(lt).ToDQL(ltFunc)
}

func LtFn(field string, value interface{}) FilterFn {
	return func() DQLizer {
		expression := Lt{}
		expression[field] = value
		return expression
	}
}

// Ge ge expression ge(field, value)
type Ge filterKV

func (ge Ge) ToDQL() (query string, args []interface{}, err error) {
	return filterKV(ge).ToDQL(geFunc)
}

func GeFn(field string, value interface{}) FilterFn {
	return func() DQLizer {
		expression := Ge{}
		expression[field] = value
		return expression
	}
}

// Gt gt expression gt(field, value)
type Gt filterKV

func (gt Gt) ToDQL() (query string, args []interface{}, err error) {
	return filterKV(gt).ToDQL(gtFunc)
}

func GtFn(field string, value interface{}) FilterFn {
	return func() DQLizer {
		expression := Gt{}
		expression[field] = value
		return expression
	}
}

// Has has expression has(field, value)
type Has filterKV

func (has Has) ToDQL() (query string, args []interface{}, err error) {
	return filterKV(has).ToDQL(gtFunc)
}

func HasFn(field string, value interface{}) FilterFn {
	return func() DQLizer {
		expression := Has{}
		expression[field] = value
		return expression
	}
}

// AllOfTerms allofterms expression allofterms(field, value)
type AllOfTerms filterKV

func (allOfTerms AllOfTerms) ToDQL() (query string, args []interface{}, err error) {
	return filterKV(allOfTerms).ToDQL(alloftermsFunc)
}

func AllOfTermsFn(field string, value interface{}) FilterFn {
	return func() DQLizer {
		expression := AllOfTerms{}
		expression[field] = value
		return expression
	}
}

// AnyOfTerms anyofterms expression anyofterms(field, value)
type AnyOfTerms filterKV

func (anyOfTerms AnyOfTerms) ToDQL() (query string, args []interface{}, err error) {
	return filterKV(anyOfTerms).ToDQL(anyoftermsFunc)
}

func AnyOfTermsFn(field string, value interface{}) FilterFn {
	return func() DQLizer {
		expression := AnyOfTerms{}
		expression[field] = value
		return expression
	}
}

// Regexp regexp expression regexp(field, /pattern/)
type Regexp filterKV

func (regexp Regexp) ToDQL() (query string, args []interface{}, err error) {
	return filterKV(regexp).ToDQL(regexpFunc)
}

func RegexpFn(field string, pattern string) FilterFn {
	return func() DQLizer {
		expression := Regexp{}
		expression[field] = pattern
		return expression
	}
}

// Match match expression match(field, /pattern/)
type Match filterKV

func (match Match) ToDQL() (query string, args []interface{}, err error) {
	return filterKV(match).ToDQL(matchFunc)
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

type RawValue struct {
	val interface{}
}

func RawVal(value interface{}) RawValue {
	return RawValue{value}
}

func Var(ref string) RawValue {
	return RawVal(ref)
}

func Val(ref string) filterValue {
	return filterValue{
		funcType: valFunc,
		value:    RawVal(ref),
	}
}

//type Val string
//
//func (val Val) ToDQL() (query string, args []interface{}, err error) {
//	query = string(valFunc) + "(" + string(val) + ")"
//	return
//}

func UID(value interface{}) filterValue {
	return filterValue{
		funcType: uidFunc,
		value:    value,
	}
}

func UIDFn(value interface{}) FilterFn {
	return func() DQLizer {
		return UID(value)
	}
}

//alloftextFunc  FuncType = "alloftext"
//anyoftextFunc  FuncType = "anyoftext"
//countFunc      FuncType = "count"
//exactFunc      FuncType = "exact"
//termFunc       FuncType = "term"
//fulltextFunc   FuncType = "fulltext"
//sumFunc        FuncType = "sum"
//betweenFunc    FuncType = "between"
//uidFunc        FuncType = "uid"
//uid      FuncType = "uid_in"

func parseValue(value interface{}) (valuePlaceholder string, args []interface{}, err error) {
	if isListType(value) {
		listValue, err := toInterfaceSlice(value)

		if err != nil {
			return "", nil, err
		}

		placeholders := make([]string, len(listValue))
		for index, value := range listValue {
			placeholders[index] = symbolValuePlaceholder
			args = append(args, value)
		}

		valuePlaceholder = fmt.Sprintf("[%s]", strings.Join(placeholders, ","))
	} else if varRef, ok := value.(RawValue); ok {
		valuePlaceholder = fmt.Sprintf("%s", varRef.val)

	} else {
		args = append(args, value)
		valuePlaceholder = symbolValuePlaceholder
	}

	return
}

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
