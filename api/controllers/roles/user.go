package roles

import (
	"net/http"

	"github.com/Utkar5hM/Execstasy/api/controllers/authentication"
	"github.com/Utkar5hM/Execstasy/api/utils/helper"
	"github.com/doug-martin/goqu/v9"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func (h *roleHandler) AddUserToRole(c echo.Context) error {
	var roleID ParamsIDStruct

	err := (&echo.DefaultBinder{}).BindPathParams(c, &roleID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Invalid path parameters", err))
	}
	if err := h.Validator.Struct(roleID); err != nil {
		return c.JSON(400, helper.ErrorMessage("Input path parameters Validation failed", err))
	}
	addUserRoleStruct := new(AddUserRole)
	if err := c.Bind(addUserRoleStruct); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Invalid request format", err))
	}
	if err := h.Validator.Struct(addUserRoleStruct); err != nil {
		return c.JSON(400, helper.ErrorMessage("Invalid request body", err))
	}
	// Get user from JWT token
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*authentication.JwtCustomClaims)

	// Get user ID from username
	var userId uint64
	var name string
	sqlUserCheck, _, _ := goqu.From("users").Where(goqu.Ex{"username": addUserRoleStruct.Username}).Select("id", "name").ToSQL()
	if err := h.DB.QueryRow(c.Request().Context(), sqlUserCheck).Scan(&userId, &name); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("User not found", err))
	}

	// Check if role exists and user is not already in it
	sqlCheck, _, _ := goqu.From("roles").
		LeftJoin(goqu.T("role_users"), goqu.On(goqu.Ex{
			"roles.id":           goqu.I("role_users.role_id"),
			"role_users.user_id": userId,
		})).
		Where(goqu.Ex{"roles.id": roleID.ID}).
		Select(
			goqu.COUNT("roles.id").As("role_count"),
			goqu.COUNT("role_users.role_id").As("member_count"),
		).ToSQL()

	var roleCount, memberCount int
	if err := h.DB.QueryRow(c.Request().Context(), sqlCheck).Scan(&roleCount, &memberCount); err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorMessage("Database error", err))
	}

	if roleCount == 0 {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Role does not exist", nil))
	}
	if memberCount > 0 {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("User is already in this role", nil))
	}

	// Add user to role
	sqlInsert, _, _ := goqu.Insert("role_users").Rows(goqu.Record{
		"role_id":  roleID.ID,
		"user_id":  userId,
		"added_by": claims.Id,
		"role":     addUserRoleStruct.Role,
	}).ToSQL()

	if _, err := h.DB.Exec(c.Request().Context(), sqlInsert); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Failed to add user to role", err))
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Successfully added user to role",
		"status":  "success",
		"data": echo.Map{
			"user_id": userId,
			"name":    name,
		},
	})
}

func (h *roleHandler) DeleteUserFromRole(c echo.Context) error {

	var roleID ParamsIDStruct

	err := (&echo.DefaultBinder{}).BindPathParams(c, &roleID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Invalid path parameters", err))
	}
	if err := h.Validator.Struct(roleID); err != nil {
		return c.JSON(400, helper.ErrorMessage("Input path parameters Validation failed", err))
	}
	deleteUserRoleStruct := new(DeleteUserRole)
	if err := c.Bind(deleteUserRoleStruct); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Invalid request format", err))
	}
	if err := h.Validator.Struct(deleteUserRoleStruct); err != nil {
		return c.JSON(400, helper.ErrorMessage("Invalid request body", err))
	}

	// Get user ID from username
	var userId uint64
	sqlUserCheck, _, _ := goqu.From("users").Where(goqu.Ex{"username": deleteUserRoleStruct.Username}).Select("id").ToSQL()
	if err := h.DB.QueryRow(c.Request().Context(), sqlUserCheck).Scan(&userId); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "User not found"})
	}

	// Check if role exists and user is in it
	sqlCheck, _, _ := goqu.From("roles").
		Join(goqu.T("role_users"), goqu.On(goqu.Ex{
			"roles.id":           goqu.I("role_users.role_id"),
			"role_users.user_id": userId,
		})).
		Where(goqu.Ex{"roles.id": roleID.ID}).
		Select(goqu.COUNT("roles.id")).ToSQL()

	var roleCount int
	if err := h.DB.QueryRow(c.Request().Context(), sqlCheck).Scan(&roleCount); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Database error"})
	}

	if roleCount == 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "User is not in this role"})
	}

	// Delete user from role
	sqlDelete, _, _ := goqu.Delete("role_users").
		Where(goqu.Ex{
			"role_id": roleID.ID,
			"user_id": userId,
		}).ToSQL()

	if _, err := h.DB.Exec(c.Request().Context(), sqlDelete); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Failed to remove user from role"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Successfully removed user from role",
		"status":  "success",
	})
}

func (h *roleHandler) getRoleUsers(c echo.Context) error {

	var roleID ParamsIDStruct

	err := (&echo.DefaultBinder{}).BindPathParams(c, &roleID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Invalid path parameters", err))
	}
	if err := h.Validator.Struct(roleID); err != nil {
		return c.JSON(400, helper.ErrorMessage("Input path parameters Validation failed", err))
	}
	// Query to get users in the role
	sql, _, _ := goqu.From("role_users").
		Join(goqu.T("users"), goqu.On(goqu.Ex{
			"role_users.user_id": goqu.I("users.id"),
		})).
		Where(goqu.Ex{"role_users.role_id": roleID.ID}).
		Select("users.id", "users.username", "users.name", "role_users.role").
		ToSQL()
	rows, err := h.DB.Query(c.Request().Context(), sql)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorMessage("Database error", err))
	}
	defer rows.Close()
	var users []map[string]interface{} = make([]map[string]interface{}, 0)
	for rows.Next() {
		var userId uint64
		var username, name, role string
		if err := rows.Scan(&userId, &username, &name, &role); err != nil {
			return c.JSON(http.StatusInternalServerError, helper.ErrorMessage("Error scanning row", err))
		}
		user := map[string]interface{}{
			"id":       userId,
			"username": username,
			"name":     name,
			"role":     role, // Assuming role ID is the same as the path parameter
		}
		users = append(users, user)
	}
	return c.JSON(http.StatusOK, users)
}

func (h *roleHandler) hasRoleEditAccess(c echo.Context) error {
	result, err := h.hasRoleEditAccessFunc(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorMessage(result.Message, err))
	}
	return c.JSON(http.StatusOK, echo.Map{
		"access":  result.HasAccess,
		"message": result.Message,
		"status":  "success",
	})
}
