package roles

import (
	"context"
	"net/http"

	"github.com/Utkar5hM/Execstasy/api/controllers/authentication"
	"github.com/Utkar5hM/Execstasy/api/utils/config"
	"github.com/doug-martin/goqu/v9"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type roleHandler struct {
	config.Handler
}

type CreateRole struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (h *roleHandler) createRole(c echo.Context) error {

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*authentication.JwtCustomClaims)
	roleStruct := new(CreateRole)
	if err := c.Bind(roleStruct); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	if roleStruct.Name == "" {
		return c.JSON(400, echo.Map{
			"error": "Role name is required",
		})
	}

	sql, _, _ := goqu.Insert("roles").Rows(
		goqu.Record{
			"name":        roleStruct.Name,
			"description": roleStruct.Description,
			"created_by":  claims.Id,
		},
	).ToSQL()
	_, err := h.DB.Exec(context.Background(), sql)
	if err != nil {
		return c.JSON(400, echo.Map{
			"error": "Failed to create role: ",
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Successfully created role",
	})
}

func (h *roleHandler) getRoles(c echo.Context) error {
	sql, _, _ := goqu.From("roles").
		Join(
			goqu.T("users"),
			goqu.On(goqu.Ex{"roles.created_by": goqu.I("users.id")}),
		).
		Select(
			goqu.I("roles.id"), // Correctly quote table and column separately
			goqu.I("roles.name"),
			goqu.I("roles.description"),
			goqu.I("users.username").As("created_by_username"),
			goqu.L("to_char(roles.updated_at, 'YYYY-MM-DD\"T\"HH24:MI:SS\"Z\"')").As("updated_at"),
			goqu.L("to_char(roles.created_at, 'YYYY-MM-DD\"T\"HH24:MI:SS\"Z\"')").As("created_at"),
		).
		Order(goqu.I("roles.id").Asc()). // Correctly quote table and column separately
		ToSQL()
	rows, err := h.DB.Query(context.Background(), sql)
	if err != nil {
		return c.JSON(400, echo.Map{
			"error": "Failed to fetch roles: " + err.Error(),
		})
	}
	defer rows.Close()

	var roles []map[string]interface{}

	// Iterate over rows and populate the slice
	for rows.Next() {
		var id int
		var name, description, createdByUsername, updatedAt, createdAt string

		// Scan the row into variables
		err := rows.Scan(&id, &name, &description, &createdByUsername, &updatedAt, &createdAt)
		if err != nil {
			return c.JSON(400, echo.Map{
				"error": "Failed to scan row: " + err.Error(),
			})
		}

		// Create a map for the current row
		role := map[string]interface{}{
			"id":          id,
			"name":        name,
			"description": description,
			"created_by":  createdByUsername,
			"updated_at":  updatedAt,
			"created_at":  createdAt,
		}

		// Append the map to the slice
		roles = append(roles, role)
	}

	// Return the roles as JSON
	return c.JSON(http.StatusOK, echo.Map{
		"roles": roles,
	})
}

type DeleteRole struct {
	Id uint64 `json:"id"`
}

func (h *roleHandler) deleteRole(c echo.Context) error {
	deleteRoleStruct := new(DeleteRole)
	if err := c.Bind(deleteRoleStruct); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	if deleteRoleStruct.Id == 0 {
		return c.JSON(400, echo.Map{
			"error": "Role ID is required",
		})
	}
	// Check if the role exists before attempting to delete
	sqlCheck, _, _ := goqu.From("roles").Where(goqu.Ex{"id": deleteRoleStruct.Id}).Select(goqu.COUNT("*")).ToSQL()
	row := h.DB.QueryRow(context.Background(), sqlCheck)
	var count int
	err := row.Scan(&count)
	if err != nil || count == 0 {
		return c.JSON(400, echo.Map{
			"error": "Role not found",
		})
	}
	// Proceed with deletion if the role exists
	sql, _, _ := goqu.Delete("roles").Where(goqu.Ex{"id": deleteRoleStruct.Id}).ToSQL()
	_, err = h.DB.Exec(context.Background(), sql)
	if err != nil {
		return c.JSON(400, echo.Map{
			"error": "Failed to delete role: " + err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Successfully deleted role",
	})
}

func (h *roleHandler) isAdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*authentication.JwtCustomClaims)
		if claims.Role != "admin" {
			return c.JSON(403, echo.Map{
				"message": "You are not authorized for this action",
			})
		}
		return next(c)
	}
}
