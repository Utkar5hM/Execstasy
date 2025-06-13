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

func (h *authHandler) searchUsers(c echo.Context) error {
	// Get the search query from the request
	query := c.QueryParam("query")
	if query == "" {
		return c.JSON(400, echo.Map{
			"error": "Query parameter is required",
		})
	}

	// Build the SQL query with filtering
	sql, _, _ := goqu.From("users").
		Where(
			goqu.Or(
				goqu.I("username").ILike("%"+query+"%"),
				goqu.I("name").ILike("%"+query+"%"),
			),
		).
		Select("id", "username", "name", "email", "role").ToSQL()

	// Execute the query
	rows, err := h.DB.Query(c.Request().Context(), sql)
	if err != nil {
		return c.JSON(500, echo.Map{
			"error": "Failed to fetch users",
		})
	}
	defer rows.Close()

	// Parse the results
	var users []map[string]interface{} = make([]map[string]interface{}, 0)
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

	// Return the filtered users
	return c.JSON(200, users)
}
