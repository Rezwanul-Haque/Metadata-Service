package converters

import (
	"fmt"
	"github.com/rezwanul-haque/Metadata-Service/src/utils/helpers"
	"math"
	"strings"
)

func RequestParamsToSqlQuery(userIds string, page float64, size float64, status string, tableName string) string {

	p := math.Max(0, page)
	s := math.Max(0, size)
	pageNum := int(p)
	pageSize := int(s)

	var sql string
	if !helpers.IsEmpty(userIds) {
		userIdsList := stringToMysqlList(userIds)
		sql += fmt.Sprintf(" AND %s.user_id IN (%s)", tableName, userIdsList)
	}

	if !helpers.IsEmpty(status) {
		sql += fmt.Sprintf(" AND %s.status IN ('%s')", tableName, status)
	}

	sql += fmt.Sprintf(" LIMIT %d, %d", (pageNum-1)*pageSize, pageSize)

	return sql
}

func stringToMysqlList(str string) string {
	var result string
	splitedString := strings.Split(str, ",")
	for i, v := range splitedString {
		if i == len(splitedString)-1 {
			result += fmt.Sprintf("'%v'", v)
		} else {
			result += fmt.Sprintf("'%v',", v)
		}

	}

	return fmt.Sprintf("%s", result)
}
