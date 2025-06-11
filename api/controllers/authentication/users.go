package authentication

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/labstack/echo/v4"
)

func (h *authHandler) getUsers(c echo.Context) error {
	sql, _, _ := goqu.From("users").
		Select("id", "username", "name", "email", "role").ToSQL()
	rows, err := h.DB.Query(c.Request().Context(), sql)
	if err != nil {
		return c.JSON(500, echo.Map{
			"error": "Failed to fetch users",
		})
	}
	defer rows.Close()
	var users []map[string]interface{}
	for rows.Next() {
		var id uint64
		var username, name, email, role string
		if err := rows.Scan(&id, &username, &name, &email, &role); err != nil {
			return c.JSON(500, echo.Map{
				"error": "Failed to scan user",
			})
		}
		users = append(users, map[string]interface{}{
			"id":       id,
			"username": username,
			"name":     name,
			"email":    email,
			"role":     role,
		})
	}
	return c.JSON(200, users)
}
