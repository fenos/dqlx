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
	alloftextFunc  FuncType = "alloftext"  // Done
	anyoftextFunc  FuncType = "anyoftext"  // Done
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
	funcType FuncType
	value    interface{}
}

func (filter Filter) ToDQL() (query string, args []interface{}, err error) {
	var placeholder string

	switch castValue := filter.value.(type) {
	case filterKV:
		return castValue.toDQL(filter.funcType)
	case Filter:
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
		args = append(args, valueArgs)
	}

	return fmt.Sprintf("%s(%s)", filter.funcType, placeholder), args, nil
}

type filterKV map[string]interface{}

func (filter filterKV) toDQL(funcType FuncType) (query string, args []interface{}, err error) {
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
	return Filter{
		funcType: eqFunc,
		value:    filterKV(eq),
	}.ToDQL()
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
	return Filter{
		funcType: leFunc,
		value:    filterKV(le),
	}.ToDQL()
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
	return Filter{
		funcType: ltFunc,
		value:    filterKV(lt),
	}.ToDQL()
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
	return Filter{
		funcType: geFunc,
		value:    filterKV(ge),
	}.ToDQL()
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
	return Filter{
		funcType: gtFunc,
		value:    filterKV(gt),
	}.ToDQL()
}

func GtFn(field string, value interface{}) FilterFn {
	return func() DQLizer {
		expression := Gt{}
		expression[field] = value
		return expression
	}
}

// Has has expression has(field, value)
type Has Filter

func (has Has) ToDQL() (query string, args []interface{}, err error) {
	return Filter{
		funcType: hasFunc,
		value:    has,
	}.ToDQL()
}

func HasFn(field string) FilterFn {
	return func() DQLizer {
		return Has{
			funcType: hasFunc,
			value:    RawVal(field),
		}
	}
}

// AllOfTerms allofterms expression allofterms(field, value)
type AllOfTerms filterKV

func (allOfTerms AllOfTerms) ToDQL() (query string, args []interface{}, err error) {
	return Filter{
		funcType: alloftermsFunc,
		value:    filterKV(allOfTerms),
	}.ToDQL()
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
	return Filter{
		funcType: anyoftermsFunc,
		value:    filterKV(anyOfTerms),
	}.ToDQL()
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
	return Filter{
		funcType: regexpFunc,
		value:    filterKV(regexp),
	}.ToDQL()
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
	return Filter{
		funcType: matchFunc,
		value:    filterKV(match),
	}.ToDQL()
}

func MatchFn(field string, pattern string) FilterFn {
	return func() DQLizer {
		expression := Match{}
		expression[field] = pattern
		return expression
	}
}

//countFunc      FuncType = "count"
//exactFunc      FuncType = "exact"
//termFunc       FuncType = "term"
//fulltextFunc   FuncType = "fulltext"
//sumFunc        FuncType = "sum"
//betweenFunc    FuncType = "between"
//uidFunc        FuncType = "uid"
//uid      FuncType = "uid_in"

// AllOfText alloftext expression alloftext(field, value)
type AllOfText filterKV

func (alloftext AllOfText) ToDQL() (query string, args []interface{}, err error) {
	return Filter{
		funcType: alloftextFunc,
		value:    filterKV(alloftext),
	}.ToDQL()
}

func AllOfTextFn(field string, pattern string) FilterFn {
	return func() DQLizer {
		expression := AllOfText{}
		expression[field] = pattern
		return expression
	}
}

// AnyOfText anyoftext expression anyoftext(field, value)
type AnyOfText filterKV

func (anyoftext AnyOfText) ToDQL() (query string, args []interface{}, err error) {
	return Filter{
		funcType: anyoftextFunc,
		value:    filterKV(anyoftext),
	}.ToDQL()
}

func AnyOfTextFn(field string, pattern string) FilterFn {
	return func() DQLizer {
		expression := AnyOfText{}
		expression[field] = pattern
		return expression
	}
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

func Val(ref string) Filter {
	return Filter{
		funcType: valFunc,
		value:    RawVal(ref),
	}
}

func UID(value interface{}) Filter {
	return Filter{
		funcType: uidFunc,
		value:    value,
	}
}

func UIDFn(value interface{}) FilterFn {
	return func() DQLizer {
		return UID(value)
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

type OrderDirection string

var (
	OrderAsc  OrderDirection = "orderasc"
	OrderDesc OrderDirection = "orderdesc"
)

type OrderBy struct {
	Direction OrderDirection
	Predicate interface{}
}

func (orderBy OrderBy) ToDQL() (query string, args []interface{}, err error) {
	placeholder, valueArgs, err := parseValue(orderBy.Predicate)

	if err != nil {
		return "", nil, err
	}

	query = fmt.Sprintf("%s:%s", orderBy.Direction, placeholder)

	args = append(args, valueArgs...)

	return
}

func parseValue(value interface{}) (valuePlaceholder string, args []interface{}, err error) {
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
	case RawValue:
		valuePlaceholder = fmt.Sprintf("%s", castType.val)
	default:
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
