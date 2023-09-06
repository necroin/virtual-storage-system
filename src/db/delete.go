package db

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func (database *Database) DeleteHandler(data io.Reader, responseWriter io.Writer) error {
	request := &Request{}
	if err := json.NewDecoder(data).Decode(request); err != nil {
		return fmt.Errorf("[Database] [Delete] [Error] failed decode json request: %s", err)
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
		return fmt.Errorf("[Database] [Delete] [Error] failed database request: %s", err)
	}

	responseWriter.Write([]byte(`{"result": true}`))

	return nil
}
