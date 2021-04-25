package deku

import (
	"fmt"
	"strings"
)

type QueryVariable struct {
	QueryBuilder
}

func Variable(rootQueryFn FilterFn) QueryVariable {
	query := Query("", rootQueryFn)
	return QueryVariable{
		QueryBuilder: query,
	}
}

func (variable QueryVariable) As(name string) QueryVariable {
	variable.QueryBuilder.rootEdge.Name = name
	return variable
}

func (variable QueryVariable) ToDQL() (query string, args []interface{}, err error) {
	innerQuery, queryArgs, err := variable.QueryBuilder.toDQL()

	if err != nil {
		return "", nil, err
	}

	args = append(args, queryArgs...)

	if variable.QueryBuilder.rootEdge.Name != "" {
		query = strings.Replace(
			innerQuery,
			variable.QueryBuilder.rootEdge.Name,
			fmt.Sprintf("%s as var", variable.QueryBuilder.rootEdge.Name),
			1,
		)
	} else {
		query = "var" + innerQuery
	}

	return query, args, nil
}
