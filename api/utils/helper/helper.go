package helper

import (
	"os"

	"github.com/labstack/echo/v4"
)

func ErrorMessage(msg string, errOrStr interface{}) echo.Map {
	var desc string
	switch v := errOrStr.(type) {
	case error:
		if os.Getenv("API_DEBUG") == "TRUE" && v != nil {
			desc = v.Error()
		} else {
			desc = "Internal Server Error"
		}
	case string:
		if os.Getenv("API_DEBUG") == "TRUE" && v != "" {
			desc = v
		} else {
			desc = "Internal Server Error"
		}
	default:
		desc = "Internal Server Error"
	}
	return echo.Map{
		"error":             msg,
		"error_description": desc,
	}
}
