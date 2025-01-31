package prettier

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	// PlaceholderDollar отвечает за подстановку значения в PG
	PlaceholderDollar = "$"

	// PlaceholderQuestion ...
	PlaceholderQuestion = "?"
)

// Pretty форматирует результат запроса в бд
func Pretty(query string, placeholder string, args ...any) string {
	for i, param := range args {
		var value string
		switch v := param.(type) {
		case string:
			value = fmt.Sprintf("%q", v)
		case []byte:
			value = fmt.Sprintf("%q", string(v))
		default:
			value = fmt.Sprintf("%v", v)
		}

		query = strings.Replace(query, fmt.Sprintf("%s%s", placeholder, strconv.Itoa(i+1)), value, -1)
	}

	query = strings.ReplaceAll(query, "\t", "")
	query = strings.ReplaceAll(query, "\n", " ")

	return strings.TrimSpace(query)
}
