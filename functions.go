package deku

import "fmt"

type QueryOperator = string
type FilterConnector = string

var (
	orConnector  FilterConnector = "OR"
	andConnector FilterConnector = "AND"

	eqOperator         QueryOperator = "eq"
	geOperator         QueryOperator = "ge"
	gtOperator         QueryOperator = "gt"
	leOperator         QueryOperator = "le"
	ltOperator         QueryOperator = "lt"
	hasOperator        QueryOperator = "has"
	alloftermsOperator QueryOperator = "allofterms"
	anyoftermsOperator QueryOperator = "anyofterms"
	regexpOperator     QueryOperator = "regexp"
	matchOperator      QueryOperator = "match"
	alloftextOperator  QueryOperator = "alloftext"
	anyoftextOperator  QueryOperator = "anyoftext"
	countOperator      QueryOperator = "count"
	exactOperator      QueryOperator = "exact"
	termOperator       QueryOperator = "term"
	fulltextOperator   QueryOperator = "fulltext"
	IEOperator         QueryOperator = "IE"
	valOperator        QueryOperator = "val"
	sumOperator        QueryOperator = "sum"
	betweenOperator    QueryOperator = "between"
	uidOperator        QueryOperator = "uid"
	uidInOperator      QueryOperator = "uid_in"
)

type QueryFunction struct {
	operator QueryOperator
	field    string
	value    interface{}
	ranges   []string
}

type Filter struct {
	connection FilterConnector
	function   *QueryFunction
}

type FilterGroup struct {
	connection FilterConnector
	filters    []*Filter
}

func (group *FilterGroup) Filter(function *QueryFunction) *FilterGroup {
	filter := &Filter{
		connection: andConnector,
		function:   function,
	}

	group.filters = append(group.filters, filter)

	return group
}

func (group *FilterGroup) OrFilter(function *QueryFunction) *FilterGroup {
	filter := &Filter{
		connection: orConnector,
		function:   function,
	}

	group.filters = append(group.filters, filter)

	return group
}

func (group *FilterGroup) toDQL() (string, []interface{}, error) {
	allArgs := []interface{}{}

	allOperations := "("
	for index, filter := range group.filters {

		if index > 0 {
			allOperations += " " + filter.connection + " "
		}

		filterOperation, args, err := filter.function.toDQL()

		if err != nil {
			return "", nil, err
		}

		allArgs = append(allArgs, args...)

		allOperations += filterOperation
	}

	allOperations += ")"

	return allOperations, allArgs, nil
}

func (queryFunction *QueryFunction) toDQL() (string, []interface{}, error) {

	value := "??"

	// TODO: raw

	switch queryFunction.operator {
	case eqOperator:
		operation := fmt.Sprintf("eq(%s,%s)", queryFunction.field, value)

		return operation, []interface{}{queryFunction.value}, nil
	default:
		return "", nil, fmt.Errorf("operator %s not supported", queryFunction.operator)
	}
}

func EqFunc(field string, value interface{}) *QueryFunction {
	return &QueryFunction{
		operator: eqOperator,
		field:    field,
		value:    value,
		ranges:   nil,
	}
}

func GeFunc(field string, value interface{}) *QueryFunction {
	return &QueryFunction{
		operator: geOperator,
		field:    field,
		value:    value,
		ranges:   nil,
	}
}

func GtFunc(field string, value interface{}) *QueryFunction {
	return &QueryFunction{
		operator: gtOperator,
		field:    field,
		value:    value,
		ranges:   nil,
	}
}

func LeFunc(field string, value interface{}) *QueryFunction {
	return &QueryFunction{
		operator: leOperator,
		field:    field,
		value:    value,
		ranges:   nil,
	}
}

func LtFunc(field string, value interface{}) *QueryFunction {
	return &QueryFunction{
		operator: ltOperator,
		field:    field,
		value:    value,
		ranges:   nil,
	}
}
