package db

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func (database *Database) Insert(table string, columns []string, values []string) error {
	sqlCommand := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, strings.Join(columns, ", "), strings.Join(values, ", "))
	_, err := database.sql.Exec(sqlCommand)
	if err != nil {
		return fmt.Errorf("[Database] [Insert] [Error] failed database request: %s", err)
	}
	return nil
}

func (database *Database) InsertRequest(request *Request) *Response {
	response := &Response{
		Records: nil,
		Success: true,
		Error:   "",
	}

	names := []string{}
	for _, field := range request.Fields {
		names = append(names, field.Name)
	}

	values := []string{}
	for _, field := range request.Fields {
		values = append(values, field.Value)
	}

	if err := database.Insert(request.Table, names, values); err != nil {
		response.Success = false
		response.Error = err.Error()
		return response
	}

	return response
}

func (database *Database) InsertHandler(data io.Reader, responseWriter io.Writer) error {
	request := &Request{}
	if err := json.NewDecoder(data).Decode(request); err != nil {
		return fmt.Errorf("[Database] [InsertHandler] [Error] failed decode json request: %s", err)
	}

	response := database.InsertRequest(request)
	if err := json.NewEncoder(responseWriter).Encode(response); err != nil {
		return fmt.Errorf("[Database] [InsertHandler] [Error] failed encode json response: %s", err)
	}

	return nil
}
