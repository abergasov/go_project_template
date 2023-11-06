package utils

import (
	"fmt"
	"strings"
)

const (
	PQParamPlaceholder = "$"
	MSParamPlaceholder = "@p"
)

func GenerateInsertSQL(tableName string, fieldsValuesMapping map[string]any) (sqlI string, params []any) {
	fields := make([]string, 0, len(fieldsValuesMapping))
	placeholders := make([]string, 0, len(fieldsValuesMapping))
	params = make([]any, 0, len(fieldsValuesMapping))
	counter := 1
	for k, v := range fieldsValuesMapping {
		params = append(params, v)
		fields = append(fields, k)
		placeholders = append(placeholders, fmt.Sprintf("$%d", counter))
		counter++
	}
	sqlI = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, strings.Join(fields, ", "), strings.Join(placeholders, ", "))
	return sqlI, params
}

// GenerateBulkInsertSQL This method generates a bulk insert SQL statement based on entity mapping
// will panic if entityList is empty
func GenerateBulkInsertSQL[T any](
	tableName string,
	paramPlaceholder string,
	entityList []T,
	entityProcessor func(entity T) map[string]any,
) (sqlI string, params []any) {
	// processor is based on map, so it random iteration. generate columns first
	columns := make([]string, 0, 10)
	for k := range entityProcessor(entityList[0]) {
		columns = append(columns, k)
	}

	// generate values
	counter := 1
	placeholders := make([]string, 0, len(entityList)*len(columns))
	params = make([]any, 0, len(entityList)*len(columns))
	for i := range entityList {
		sqlMapping := entityProcessor(entityList[i])
		localPlaceholders := make([]string, 0, len(columns))
		for j := range columns {
			params = append(params, sqlMapping[columns[j]])
			localPlaceholders = append(localPlaceholders, fmt.Sprintf("%s%d", paramPlaceholder, counter))
			counter++
		}
		placeholders = append(placeholders, fmt.Sprintf("(%s)", strings.Join(localPlaceholders, ",")))
	}

	sqlI = fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", tableName, strings.Join(columns, ","), strings.Join(placeholders, ","))
	return sqlI, params
}
