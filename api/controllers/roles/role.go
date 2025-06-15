package roles

import (
	"context"
	"net/http"

	"github.com/Utkar5hM/Execstasy/api/controllers/authentication"
	"github.com/Utkar5hM/Execstasy/api/utils/config"
	"github.com/Utkar5hM/Execstasy/api/utils/helper"
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
		return c.JSON(400, helper.ErrorMessage("Role name is required", nil))
	}

	sql, _, _ := goqu.Insert("roles").Rows(
		goqu.Record{
			"name":        roleStruct.Name,
			"description": roleStruct.Description,
			"created_by":  claims.Id,
		},
	).Returning(goqu.I("id")).ToSQL()
	row := h.DB.QueryRow(context.Background(), sql)
	var roleID int
	err := row.Scan(&roleID)
	if err != nil {
		return c.JSON(400, helper.ErrorMessage("Failed to create role", err))
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Successfully created role",
		"id":      roleID,
	})
}

type ParamsIDStruct struct {
	ID uint64 `param:"id"` // Match the path parameter name
}
type Role struct {
	ID          int    `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	CreatedBy   string `db:"createdBy"`
	CreatedAt   string `db:"createdAt"`
	UpdatedAt   string `db:"updatedAt"`
}

func (h *roleHandler) getRole(c echo.Context) error {
	var roleID ParamsIDStruct

	err := (&echo.DefaultBinder{}).BindPathParams(c, &roleID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Invalid path parameters", err))
	}
	if roleID.ID == 0 {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Role ID is required", nil))
	}
	sql, _, _ := goqu.From("roles").
		Join(
			goqu.T("users"),
			goqu.On(goqu.Ex{"roles.created_by": goqu.I("users.id")}),
		).Where(goqu.Ex{"roles.id": roleID.ID}).
		Select(
			goqu.I("roles.id"),
			goqu.I("roles.name"),
			goqu.I("roles.description"),
			goqu.I("users.username").As("createdBy"),
			goqu.L("to_char(roles.created_at, 'YYYY-MM-DD\"T\"HH24:MI:SS\"Z\"')").As("createdAt"),
			goqu.L("to_char(roles.updated_at, 'YYYY-MM-DD\"T\"HH24:MI:SS\"Z\"')").As("updatedAt"),
		).ToSQL()
	row := h.DB.QueryRow(context.Background(), sql)
	var role Role
	err = row.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedBy, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Failed to scan role", err))
	}
	return c.JSON(http.StatusOK, role)
}

func (h *roleHandler) getRoles(c echo.Context) error {
	sql, _, _ := goqu.From("roles").
		Join(
			goqu.T("users"),
			goqu.On(goqu.Ex{"roles.created_by": goqu.I("users.id")}),
		).
		LeftJoin(
			goqu.T("role_users"), // Assuming a `role_users` table exists to map roles to users
			goqu.On(goqu.Ex{"roles.id": goqu.I("role_users.role_id")}),
		).
		Select(
			goqu.I("roles.id"), // Correctly quote table and column separately
			goqu.I("roles.name"),
			goqu.I("roles.description"),
			goqu.I("users.username").As("created_by_username"),
			goqu.L("to_char(roles.updated_at, 'YYYY-MM-DD\"T\"HH24:MI:SS\"Z\"')").As("updated_at"),
			goqu.L("to_char(roles.created_at, 'YYYY-MM-DD\"T\"HH24:MI:SS\"Z\"')").As("created_at"),
			goqu.L("COUNT(role_users.user_id)").As("role_users"), // Count the number of users in each role
		).
		GroupBy(
			goqu.I("roles.id"),
			goqu.I("roles.name"),
			goqu.I("roles.description"),
			goqu.I("users.username"),
			goqu.L("to_char(roles.updated_at, 'YYYY-MM-DD\"T\"HH24:MI:SS\"Z\"')"),
			goqu.L("to_char(roles.created_at, 'YYYY-MM-DD\"T\"HH24:MI:SS\"Z\"')"),
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

	var roles []map[string]interface{} = make([]map[string]interface{}, 0)

	// Iterate over rows and populate the slice
	for rows.Next() {
		var id int
		var name, description, createdByUsername, updatedAt, createdAt, role_users string

		// Scan the row into variables
		err := rows.Scan(&id, &name, &description, &createdByUsername, &updatedAt, &createdAt, &role_users)
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
			"users":       role_users,
		}

		// Append the map to the slice
		roles = append(roles, role)
	}

	// Return the roles as JSON
	return c.JSON(http.StatusOK, roles)
}

func (h *roleHandler) deleteRole(c echo.Context) error {
	var roleID ParamsIDStruct

	err := (&echo.DefaultBinder{}).BindPathParams(c, &roleID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Invalid path parameters", err))
	}
	if roleID.ID == 0 {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Role ID is required", nil))
	}
	// Check if the role exists before attempting to delete
	sqlCheck, _, _ := goqu.From("roles").Where(goqu.Ex{"id": roleID.ID}).Select(goqu.COUNT("*")).ToSQL()
	row := h.DB.QueryRow(context.Background(), sqlCheck)
	var count int
	err = row.Scan(&count)
	if err != nil || count == 0 {
		return c.JSON(400, echo.Map{
			"error": "Role not found",
		})
	}
	// Proceed with deletion if the role exists
	sql, _, _ := goqu.Delete("roles").Where(goqu.Ex{"id": roleID.ID}).ToSQL()
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

func (h *roleHandler) editRole(c echo.Context) error {
	var roleID ParamsIDStruct

	err := (&echo.DefaultBinder{}).BindPathParams(c, &roleID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Invalid path parameters", err))
	}
	if roleID.ID == 0 {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Role ID is required", nil))
	}
	var roleStruct CreateRole
	if err := c.Bind(&roleStruct); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Failed to bind request body", err))
	}
	if roleStruct.Name == "" {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Role name is required", nil))
	}
	// check if role exists
	sqlCheck, _, _ := goqu.From("roles").Where(goqu.Ex{"id": roleID.ID}).Select(goqu.COUNT("*")).ToSQL()
	row := h.DB.QueryRow(context.Background(), sqlCheck)
	var count int
	err = row.Scan(&count)
	if err != nil {
		return c.JSON(400, helper.ErrorMessage("Failed to check role existence", err))
	}
	if count == 0 {
		return c.JSON(400, helper.ErrorMessage("Role not found", nil))
	}
	sql, _, _ := goqu.Update("roles").Set(
		goqu.Record{
			"name":        roleStruct.Name,
			"description": roleStruct.Description,
		},
	).Where(goqu.Ex{"id": roleID.ID}).ToSQL()
	_, err = h.DB.Exec(context.Background(), sql)
	if err != nil {
		return c.JSON(400, helper.ErrorMessage("Failed to update role", err))
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Successfully updated role",
		"status":  "success",
	})
}
