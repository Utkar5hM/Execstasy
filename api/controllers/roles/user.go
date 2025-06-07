package roles

import (
	"context"
	"net/http"

	"github.com/Utkar5hM/Execstasy/api/controllers/authentication"
	"github.com/doug-martin/goqu/v9"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type AddUserRole struct {
	Id       uint64 `json:"id"`
	Username string `json:"username"`
}

func (h *roleHandler) AddUserToRole(c echo.Context) error {
	addUserRoleStruct := new(AddUserRole)
	if err := c.Bind(addUserRoleStruct); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request format"})
	}

	if addUserRoleStruct.Id == 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Role ID is required"})
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
				"roles.id":           addUserRoleStruct.Id,
			}).
			Select(goqu.COUNT("role_users.user_id")).ToSQL()
		var isAdmin int
		if err := h.DB.QueryRow(context.Background(), sqlCheckAdmin).Scan(&isAdmin); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Database error"})
		}
		if isAdmin == 0 {
			return c.JSON(http.StatusForbidden, echo.Map{"error": "You do not have permission to add users to roles"})
		}
	}

	// Get user ID from username
	var userId uint64
	sqlUserCheck, _, _ := goqu.From("users").Where(goqu.Ex{"username": addUserRoleStruct.Username}).Select("id").ToSQL()
	if err := h.DB.QueryRow(context.Background(), sqlUserCheck).Scan(&userId); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "User not found"})
	}

	// Check if role exists and user is not already in it
	sqlCheck, _, _ := goqu.From("roles").
		LeftJoin(goqu.T("role_users"), goqu.On(goqu.Ex{
			"roles.id":           goqu.I("role_users.role_id"),
			"role_users.user_id": userId,
		})).
		Where(goqu.Ex{"roles.id": addUserRoleStruct.Id}).
		Select(
			goqu.COUNT("roles.id").As("role_count"),
			goqu.COUNT("role_users.role_id").As("member_count"),
		).ToSQL()

	var roleCount, memberCount int
	if err := h.DB.QueryRow(context.Background(), sqlCheck).Scan(&roleCount, &memberCount); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Database error"})
	}

	if roleCount == 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Role not found"})
	}
	if memberCount > 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "User is already in the role"})
	}

	// Add user to role
	sqlInsert, _, _ := goqu.Insert("role_users").Rows(goqu.Record{
		"role_id":  addUserRoleStruct.Id,
		"user_id":  userId,
		"added_by": claims.Id,
	}).ToSQL()

	if _, err := h.DB.Exec(context.Background(), sqlInsert); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Failed to add user to role"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Successfully added user to role"})
}

func (h *roleHandler) DeleteUserFromRole(c echo.Context) error {
	deleteUserRoleStruct := new(AddUserRole)
	if err := c.Bind(deleteUserRoleStruct); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request format"})
	}

	if deleteUserRoleStruct.Id == 0 || deleteUserRoleStruct.Username == "" {
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
				"roles.id":           deleteUserRoleStruct.Id,
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
		Where(goqu.Ex{"roles.id": deleteUserRoleStruct.Id}).
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
			"role_id": deleteUserRoleStruct.Id,
			"user_id": userId,
		}).ToSQL()

	if _, err := h.DB.Exec(context.Background(), sqlDelete); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Failed to remove user from role"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Successfully removed user from role"})
}
