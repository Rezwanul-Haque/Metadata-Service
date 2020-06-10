package converters

import (
	"fmt"
	"github.com/rezwanul-haque/Metadata-Service/src/utils/errors"
	"reflect"
	"strings"
)

type combinedSQL struct {
	query string
}

func IsInstanceOf(objectPtr, typePtr interface{}) bool {
	return reflect.TypeOf(objectPtr) == reflect.TypeOf(typePtr)
}

func GetSQLQueryClauseFromDSL(queryDSL interface{}, tableName string, jsonColumnName string) (*string, error) {
	var q = &combinedSQL{
		query: "",
	}
	for key, value := range queryDSL.(map[string]interface{}) {
		switch strings.ToLower(key) {
		case "match":
			q.buildMatchQueryClause(value, tableName, jsonColumnName)
		case "except":
			q.buildExceptQueryClause(value, tableName, jsonColumnName)
		case "range":
			q.buildRangeQueryClause(value, jsonColumnName)
		default:
			return nil, errors.NewError(fmt.Sprintf("'%v' search is not supported", key))
		}
	}
	return &q.query, nil
}

func translateOperator(operator string, match bool) (*string, *errors.RestErr) {
	var result string
	if strings.ToLower(operator) == "any" && !match {
		result = " AND "
	}
	if strings.ToLower(operator) == "all" && !match {
		result = " OR "
	}
	if strings.ToLower(operator) == "any" && match {
		result = " OR "
	}
	if strings.ToLower(operator) == "all" && match {
		result = " AND "
	}
	if strings.ToLower(operator) != "any" && strings.ToLower(operator) != "all" {
		return nil, errors.NewBadRequestError(fmt.Sprintf("Operator '%s' is unknown!", operator))
	}
	return &result, nil
}

func (q *combinedSQL) buildMatchQueryClause(matchObj interface{}, tableName string, jsonColumnName string) {
	for key, value := range matchObj.(map[string]interface{}) {
		q.buildJSONContainsClause(tableName+"."+jsonColumnName, value, "$."+key, true)
	}
}

func (q *combinedSQL) buildExceptQueryClause(exceptObj interface{}, tableName string, jsonColumnName string) {
	for key, value := range exceptObj.(map[string]interface{}) {
		q.buildJSONContainsClause(tableName+"."+jsonColumnName, value, "$."+key, false)
	}
}

func (q *combinedSQL) buildRangeQueryClause(exceptObj interface{}, jsonColumnName string) {
	for key, value := range exceptObj.(map[string]interface{}) {
		q.buildComparisonClause(value, jsonColumnName, "$."+key)
	}
}

func (q *combinedSQL) buildJSONSpecificContainsClause(target string, candidate interface{}, path string, match bool, appendOperator bool) {
	if appendOperator {
		q.query += " AND "
	}
	q.query += fmt.Sprintf("JSON_CONTAINS(%v, '%v', '%v')", target, candidate, path)
	if match {
		q.query += "=1"
	} else {
		q.query += "=0"
	}
}

func (q *combinedSQL) buildJSONContainsClause(target string, candidate interface{}, path string, match bool) *errors.RestErr {
	if IsInstanceOf(candidate, []interface{}{}) {
		slice := candidate.([]interface{})
		if len(slice) == 0 {
			return errors.NewBadRequestError(fmt.Sprintf("%v list contains no item.", path))
		} else {
			sqlList := sliceToMysqlList(slice)
			q.buildJSONSpecificContainsClause(target, sqlList, path, match, true)
		}
	} else if IsInstanceOf(candidate, map[string]interface{}{}) {
		for key, value := range candidate.(map[string]interface{}) {
			if !match {
				operator, _ := translateOperator(key, false)
				q.buildOperatorJSONContainsClause(target, value, path, false, *operator)
			} else {
				if strings.ToLower(key) == "any" {
					operator, _ := translateOperator(key, true)
					sqlList := sliceToMysqlList(value.([]interface{}))
					q.buildOperatorJSONContainsClause(target, sqlList, path, true, *operator)
				} else if strings.ToLower(key) == "all" {
					sqlList := sliceToMysqlList(value.([]interface{}))
					q.buildJSONSpecificContainsClause(target, sqlList, path, true, true)
				} else {
					return errors.NewBadRequestError(fmt.Sprintf("Operator %v is unknown", key))
				}
			}
		}
	} else {
		q.buildJSONSpecificContainsClause(target, candidate, path, match, true)
	}
	return nil
}

func (q *combinedSQL) buildComparisonClause(payload interface{}, jsonColumnName string, path string) *errors.RestErr {
	for key, value := range payload.(map[string]interface{}) {
		q.query += fmt.Sprintf(" AND %v -> '%v'", jsonColumnName, path)
		if strings.ToLower(key) == "lt" {
			q.query += fmt.Sprintf(" < %v", value)
		} else if strings.ToLower(key) == "lte" {
			q.query += fmt.Sprintf(" <= %v", value)
		} else if strings.ToLower(key) == "gt" {
			q.query += fmt.Sprintf(" > %v", value)
		} else if strings.ToLower(key) == "gte" {
			q.query += fmt.Sprintf(" >= %v", value)
		} else {
			return errors.NewBadRequestError(fmt.Sprintf("Comparison operator %v  is not supported.", key))
		}
	}
	return nil
}

func (q *combinedSQL) buildOperatorJSONContainsClause(target string, candidate interface{}, path string, match bool, operator string) *errors.RestErr {
	if IsInstanceOf(candidate, []interface{}{}) {
		slice := candidate.([]interface{})
		if len(slice) == 0 {
			return errors.NewBadRequestError(fmt.Sprintf("%v list contains no item.", path))
		} else {
			appendOperator := false
			for range slice {
				if !appendOperator {
					q.query += " AND ("
				} else {
					q.query += operator
				}
				sqlList := sliceToMysqlList(slice)
				q.buildJSONSpecificContainsClause(target, sqlList, path, match, false)
				appendOperator = true
			}
			q.query += ")"
		}
	} else {
		q.buildJSONSpecificContainsClause(target, candidate, path, match, true)
	}
	return nil
}

func sliceToMysqlList(slice []interface{}) string {
	var result string
	for i, v := range slice {
		if i == len(slice)-1 {
			result += fmt.Sprintf("\"%v\"", v)
		} else {
			result += fmt.Sprintf("\"%v\",", v)
		}

	}

	return fmt.Sprintf("[%v]", result)
}
