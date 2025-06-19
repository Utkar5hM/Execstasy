package roles

import (
	"net/http"

	"github.com/Utkar5hM/Execstasy/api/controllers/authentication"
	"github.com/Utkar5hM/Execstasy/api/utils/helper"
	"github.com/doug-martin/goqu/v9"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type RoleEditAccessResult struct {
	HasAccess bool
	Message   string
}

func (h *roleHandler) hasRoleEditAccessFunc(c echo.Context) (*RoleEditAccessResult, error) {
	var roleID ParamsIDStruct

	err := (&echo.DefaultBinder{}).BindPathParams(c, &roleID)
	if err != nil {
		return &RoleEditAccessResult{
			HasAccess: false,
			Message:   "Invalid path parameters",
		}, err
	}
	if err := h.Validator.Struct(roleID); err != nil {
		return &RoleEditAccessResult{
			HasAccess: false,
			Message:   "Role ID is required",
		}, nil
	}
	// Get user from JWT token
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*authentication.JwtCustomClaims)
	// check if user is admin with query or is part of the role with role_users.role_users_role = 'admin'
	sql, _, _ := goqu.From("role_users").
		Join(goqu.T("users"), goqu.On(goqu.Ex{
			"role_users.user_id": goqu.I("users.id"),
		})).
		Where(goqu.Or(goqu.Ex{
			"role_users.role_id": roleID.ID,
			"role_users.role":    "admin",
			"users.id":           claims.Id,
		}, goqu.Ex{
			"users.role": "admin",
			"users.id":   claims.Id,
		})).Select(goqu.COUNT("*")).ToSQL()
	var count int
	if err := h.DB.QueryRow(c.Request().Context(), sql).Scan(&count); err != nil {
		return &RoleEditAccessResult{
			HasAccess: false,
			Message:   "Database error",
		}, err
	}
	if count == 0 {
		return &RoleEditAccessResult{
			HasAccess: false,
			Message:   "You do not have permission to edit this role",
		}, nil
	}
	return &RoleEditAccessResult{
		HasAccess: true,
		Message:   "You are allowed to update this role",
	}, nil
}

func (h *roleHandler) hasRoleEditAccessMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		result, err := h.hasRoleEditAccessFunc(c)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.ErrorMessage(result.Message, err))
		}
		if !result.HasAccess {
			return c.JSON(http.StatusForbidden, helper.ErrorMessage(result.Message, nil))
		}
		return next(c)
	}

}
