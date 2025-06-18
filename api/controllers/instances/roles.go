package instances

import (
	"fmt"
	"net/http"

	"github.com/Utkar5hM/Execstasy/api/utils/helper"
	"github.com/doug-martin/goqu/v9"
	"github.com/labstack/echo/v4"
)

func (h *instanceHandler) getInstanceRoles(c echo.Context) error {
	var instance ParamsIDStruct

	// Bind path parameters using DefaultBinder
	err := (&echo.DefaultBinder{}).BindPathParams(c, &instance)
	if err != nil || instance.ID == 0 {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Invalid path parameters", err))
	}
	sql, _, _ := goqu.From("instance_roles").Join(
		goqu.T("roles"),
		goqu.On(goqu.Ex{"instance_roles.role_id": goqu.I("roles.id")}),
	).Select(
		goqu.I("roles.id"),
		goqu.I("roles.name"),
		goqu.I("instance_roles.instance_host_username"),
	).Where(
		goqu.I("instance_roles.instance_id").Eq(instance.ID),
	).ToSQL()

	rows, err := h.DB.Query(c.Request().Context(), sql)
	if err != nil {
		return c.JSON(400, helper.ErrorMessage("Failed to get instance roles", err))
	}
	defer rows.Close()

	var roles []InstanceRole = []InstanceRole{}
	for rows.Next() {
		var role InstanceRole
		if err := rows.Scan(&role.Id, &role.Name, &role.HostUsername); err != nil {
			return c.JSON(http.StatusInternalServerError, helper.ErrorMessage("Failed to scan instance role", err))
		}
		roles = append(roles, role)
	}

	return c.JSON(http.StatusOK, roles)
}

func (h *instanceHandler) addInstanceRoles(c echo.Context) error {
	var instance ParamsIDStruct
	err := (&echo.DefaultBinder{}).BindPathParams(c, &instance)
	if err != nil || instance.ID == 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":             "Invalid path parameters",
			"error_description": err.Error(),
		})
	}

	var role IDUsernameStruct
	err = c.Bind(&role)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid request body")
	}
	if role.ID == 0 || role.Username == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":             "Invalid role ID or username",
			"error_description": "Role ID and username must be provided",
		})
	}
	// check if given instance_role exists, error if do
	sql, _, err := goqu.From("instance_roles").Where(
		goqu.Ex{"instance_id": instance.ID, "role_id": role.ID, "instance_host_username": role.Username},
	).Select(goqu.COUNT("*")).ToSQL()
	var count int64
	err = h.DB.QueryRow(c.Request().Context(), sql).Scan(&count)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":             "Failed to check existing instance role: ",
			"error_description": err.Error(),
		})
	}
	if count > 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":             "Instance role already exists",
			"error_description": fmt.Sprintf("Instance role with ID %d and username %s already exists", role.ID, role.Username),
		})
	}
	sql, _, _ = goqu.Insert("instance_roles").Rows(
		goqu.Record{
			"instance_id":            instance.ID,
			"role_id":                role.ID,
			"instance_host_username": role.Username,
		},
	).ToSQL()
	_, err = h.DB.Exec(c.Request().Context(), sql)
	if err != nil {
		return c.JSON(400, echo.Map{
			"error":             "Failed to add instance role",
			"error_description": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Successfully added instance role",
		"status":  "success",
	})

}

func (h *instanceHandler) deleteInstanceRoles(c echo.Context) error {
	var instance ParamsIDStruct
	err := (&echo.DefaultBinder{}).BindPathParams(c, &instance)
	if err != nil || instance.ID == 0 {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Invalid path parameters", err))
	}
	var role IDUsernameStruct
	err = c.Bind(&role)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Invalid request body", err))
	}
	if role.ID == 0 || role.Username == "" {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Invalid request body", "Role ID and username must be provided"))
	}
	// check if instance role exists, error if not
	sql, _, _ := goqu.From("instance_roles").Where(
		goqu.Ex{"instance_id": instance.ID, "role_id": role.ID, "instance_host_username": role.Username},
	).Select(goqu.COUNT("*")).ToSQL()
	var count int64
	err = h.DB.QueryRow(c.Request().Context(), sql).Scan(&count)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorMessage("Failed to check existing instance role", err))
	}
	if count == 0 {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Instance role does not exist", fmt.Sprintf("Instance role with ID %d and username %s does not exist", role.ID, role.Username)))
	}
	sql, _, _ = goqu.Delete("instance_roles").Where(
		goqu.Ex{"instance_id": instance.ID, "role_id": role.ID, "instance_host_username": role.Username},
	).ToSQL()
	_, err = h.DB.Exec(c.Request().Context(), sql)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorMessage("Failed to delete instance role", err))
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Successfully deleted instance role with ID",
		"status":  "success",
	})
}
