package roles

import (
	"context"
	"net/http"

	"github.com/Utkar5hM/Execstasy/api/controllers/authentication"
	"github.com/Utkar5hM/Execstasy/api/utils/helper"
	"github.com/doug-martin/goqu/v9"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type AddUserRole struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}

func (h *roleHandler) AddUserToRole(c echo.Context) error {
	var roleID ParamsIDStruct

	err := (&echo.DefaultBinder{}).BindPathParams(c, &roleID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Invalid path parameters", err))
	}
	if roleID.ID == 0 {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Role ID is required", nil))
	}
	addUserRoleStruct := new(AddUserRole)
	if err := c.Bind(addUserRoleStruct); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Invalid request format", err))
	}

	if addUserRoleStruct.Username == "" {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Username is required", nil))
	}
	if addUserRoleStruct.Role != "admin" && addUserRoleStruct.Role != "standard" {
		addUserRoleStruct.Role = "standard"
	}
	// Get user from JWT token
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*authentication.JwtCustomClaims)

	if claims.Role != "admin" {
		// check if the claims.user is present in role and has role_members_role = 'admin'
		sqlCheckAdmin, _, _ := goqu.From("role_users").
			Join(goqu.T("roles"), goqu.On(goqu.Ex{
				"role_users.role_id": goqu.I("roles.id"),
			})).
			Where(goqu.Ex{
				"role_users.user_id": claims.Id,
				"role_users.role":    "admin",
				"roles.id":           roleID.ID,
			}).
			Select(goqu.COUNT("role_users.user_id")).ToSQL()
		var isAdmin int
		if err := h.DB.QueryRow(context.Background(), sqlCheckAdmin).Scan(&isAdmin); err != nil {
			return c.JSON(http.StatusInternalServerError, helper.ErrorMessage("Database error", err))
		}
		if isAdmin == 0 {
			return c.JSON(http.StatusForbidden, helper.ErrorMessage("You do not have permission to add users to roles", nil))
		}
	}

	// Get user ID from username
	var userId uint64
	sqlUserCheck, _, _ := goqu.From("users").Where(goqu.Ex{"username": addUserRoleStruct.Username}).Select("id").ToSQL()
	if err := h.DB.QueryRow(context.Background(), sqlUserCheck).Scan(&userId); err != nil {
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
	if err := h.DB.QueryRow(context.Background(), sqlCheck).Scan(&roleCount, &memberCount); err != nil {
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

	if _, err := h.DB.Exec(context.Background(), sqlInsert); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Failed to add user to role", err))
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Successfully added user to role",
		"status": "success"})
}

func (h *roleHandler) DeleteUserFromRole(c echo.Context) error {
	var roleID ParamsIDStruct

	err := (&echo.DefaultBinder{}).BindPathParams(c, &roleID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Invalid path parameters", err))
	}
	if roleID.ID == 0 {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Role ID is required", nil))
	}
	deleteUserRoleStruct := new(AddUserRole)
	if err := c.Bind(deleteUserRoleStruct); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request format"})
	}

	if deleteUserRoleStruct.Username == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Role ID and Username are required"})
	}

	// Get user from JWT token
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*authentication.JwtCustomClaims)

	if claims.Role != "admin" {
		// check if the claims.user is present in role and has role_members_role = 'admin'
		sqlCheckAdmin, _, _ := goqu.From("role_users").
			Join(goqu.T("roles"), goqu.On(goqu.Ex{
				"role_users.role_id": goqu.I("roles.id"),
			})).
			Where(goqu.Ex{
				"role_users.user_id": claims.Id,
				"role_users.role":    "admin",
				"roles.id":           roleID.ID,
			}).
			Select(goqu.COUNT("role_users.user_id")).ToSQL()

		var isAdmin int
		if err := h.DB.QueryRow(context.Background(), sqlCheckAdmin).Scan(&isAdmin); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Database error"})
		}
		if isAdmin == 0 {
			return c.JSON(http.StatusForbidden, echo.Map{"error": "You do not have permission to remove users from roles"})
		}
	}

	// Get user ID from username
	var userId uint64
	sqlUserCheck, _, _ := goqu.From("users").Where(goqu.Ex{"username": deleteUserRoleStruct.Username}).Select("id").ToSQL()
	if err := h.DB.QueryRow(context.Background(), sqlUserCheck).Scan(&userId); err != nil {
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
	if err := h.DB.QueryRow(context.Background(), sqlCheck).Scan(&roleCount); err != nil {
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

	if _, err := h.DB.Exec(context.Background(), sqlDelete); err != nil {
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
	if roleID.ID == 0 {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Role ID is required", nil))
	}
	// Query to get users in the role
	sql, _, _ := goqu.From("role_users").
		Join(goqu.T("users"), goqu.On(goqu.Ex{
			"role_users.user_id": goqu.I("users.id"),
		})).
		Where(goqu.Ex{"role_users.role_id": roleID.ID}).
		Select("users.id", "users.username", "users.name", "role_users.role").
		ToSQL()
	rows, err := h.DB.Query(context.Background(), sql)
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
