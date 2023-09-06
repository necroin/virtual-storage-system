package db

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func (database *Database) DeleteRequest(request *Request) *Response {
	response := &Response{
		Records: nil,
		Success: true,
		Error:   "",
	}

	sqlCommand := fmt.Sprintf("DELETE FROM %s", request.Table)

	if len(request.Filters) != 0 {
		filters := []string{}
		for _, filter := range request.Filters {
			filters = append(filters, fmt.Sprintf("%s %s %s", filter.Name, filter.Operator, filter.Value))
		}
		sqlFilters := " WHERE " + strings.Join(filters, " AND ")
		sqlCommand = sqlCommand + sqlFilters
	}

	_, err := database.sql.Exec(sqlCommand)
	if err != nil {
		response.Success = false
		response.Error = fmt.Sprintf("[Database] [DeleteRequest] [Error] failed database request: %s", err)
	}

	return response
}

func (database *Database) DeleteHandler(data io.Reader, responseWriter io.Writer) error {
	request := &Request{}
	if err := json.NewDecoder(data).Decode(request); err != nil {
		return fmt.Errorf("[Database] [DeleteHandler] [Error] failed decode json request: %s", err)
	}

	response := database.DeleteRequest(request)
	if err := json.NewEncoder(responseWriter).Encode(response); err != nil {
		return fmt.Errorf("[Database] [DeleteHandler] [Error] failed encode json response: %s", err)
	}

	return nil
}
