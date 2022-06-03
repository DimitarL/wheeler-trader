package storage

import (
	"fmt"
	"strings"
)

type databaseType int

const (
	character databaseType = iota
	numeric
)

type columnInfo struct {
	name   string
	dbType databaseType
}

type whereClauseBuilder struct {
	where     strings.Builder
	sqlParams []interface{}

	paramNameToColumnInfo func(string) columnInfo
}

func newWhereClauseBuilder(paramNameToColumnInfo func(string) columnInfo) *whereClauseBuilder {
	return &whereClauseBuilder{paramNameToColumnInfo: paramNameToColumnInfo}
}

func (b *whereClauseBuilder) Build(params map[string]interface{}) (string, []interface{}) {
	b.where.WriteString(" WHERE ")

	b.addParamsIfPresent(params)

	return b.where.String(), b.sqlParams
}

func (b *whereClauseBuilder) addParamsIfPresent(params map[string]interface{}) {
	for paramName, value := range params {
		b.addToWhereClause(paramName)
		b.sqlParams = append(b.sqlParams, value)
	}
}

func (b *whereClauseBuilder) addToWhereClause(paramName string) {
	if len(b.sqlParams) >= 1 {
		b.where.WriteString(" AND ")
	}

	var operator string

	columnInfo := b.paramNameToColumnInfo(paramName)
	if columnInfo.dbType == character {
		operator = "="
	} else if strings.HasPrefix(strings.ToLower(paramName), "min") {
		operator = ">="
	} else if strings.HasPrefix(strings.ToLower(paramName), "max") {
		operator = "<="
	} else {
		panic(fmt.Sprintf("cannot handle numeric parameter '%s'", paramName))
	}

	b.where.WriteString(fmt.Sprintf("%s %s $%d", columnInfo.name, operator, len(b.sqlParams)+1))
}
