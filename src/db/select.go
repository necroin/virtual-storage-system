package db

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func (database *Database) Select(request *Request) *Response {
	response := &Response{
		Records: []Record{},
		Success: true,
		Error:   "",
	}

	sqlCommand := fmt.Sprintf("SELECT * FROM %s", request.Table)

	if len(request.Fields) != 0 {
		fields := []string{}
		for _, field := range request.Fields {
			fields = append(fields, field.Name)
		}
		sqlCommand = fmt.Sprintf("SELECT %s FROM %s", strings.Join(fields, ", "), request.Table)
	}

	if len(request.Filters) != 0 {
		filters := []string{}
		for _, filter := range request.Filters {
			filters = append(filters, fmt.Sprintf("%s %s %s", filter.Name, filter.Operator, filter.Value))
		}
		sqlFilters := " WHERE " + strings.Join(filters, " AND ")
		sqlCommand = sqlCommand + sqlFilters
	}

	result, err := database.sql.Query(sqlCommand)
	if err != nil {
		response.Success = false
		response.Error = fmt.Sprintf("[Database] [Select] [Error] failed database request: %s", err)
		return response
	}
	defer result.Close()

	columns, err := result.Columns()
	if err != nil {
		response.Success = false
		response.Error = fmt.Sprintf("[Database] [Select] [Error] failed get result columns: %s", err)
		return response
	}

	for result.Next() {
		values := make([]string, len(columns))
		valuesPointers := make([]interface{}, len(columns))
		for i := range values {
			valuesPointers[i] = &values[i]
		}

		if err := result.Scan(valuesPointers...); err != nil {
			response.Success = false
			response.Error = fmt.Sprintf("[Database] [Select] [Error] failed scan row values: %s", err)
			return response
		}

		record := Record{
			Fields: []Field{},
		}

		for i, column := range columns {
			record.Fields = append(record.Fields, Field{
				Name:  column,
				Value: values[i],
			})
		}

		response.Records = append(response.Records, record)
	}

	return response
}

func (database *Database) SelectHadler(data io.Reader, responseWriter io.Writer) error {
	request := &Request{}
	if err := json.NewDecoder(data).Decode(request); err != nil {
		return fmt.Errorf("[Database] [Select] [Error] failed decode json request: %s", err)
	}

	response := database.Select(request)

	if err := json.NewEncoder(responseWriter).Encode(response); err != nil {
		return fmt.Errorf("[Database] [Select] [Error] failed encode json response: %s", err)
	}

	return nil
}
