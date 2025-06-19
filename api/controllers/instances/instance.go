package instances

import (
	"net/http"

	"github.com/Utkar5hM/Execstasy/api/controllers/authentication"
	"github.com/Utkar5hM/Execstasy/api/utils/helper"
	"github.com/doug-martin/goqu/v9"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// id, name, host_address, status, created_by
func (h *instanceHandler) getInstances(c echo.Context) error {
	sql, _, _ := goqu.From("instances").Join(
		goqu.T("users"),
		goqu.On(goqu.Ex{"instances.created_by": goqu.I("users.id")}),
	).Select(
		"instances.id",
		"instances.name",
		"instances.host_address",
		"instances.status",
		goqu.I("users.username").As("created_by"),
	).ToSQL()
	rows, err := h.DB.Query(c.Request().Context(), sql)
	if err != nil {
		return c.JSON(400, helper.ErrorMessage("Failed to fetch instances: ", err))
	}
	defer rows.Close()

	var instances []Instances
	for rows.Next() {
		var instance Instances
		if err := rows.Scan(&instance.ID, &instance.Name, &instance.HostAddress, &instance.Status, &instance.CreatedBy); err != nil {
			return c.JSON(http.StatusInternalServerError, helper.ErrorMessage("Failed to fetch instance", err))
		}
		instances = append(instances, instance)
	}

	return c.JSON(http.StatusOK, instances)
}

func (h *instanceHandler) getMyInstances(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*authentication.JwtCustomClaims)
	userID := claims.Id
	// Instances the user has access to
	// instance_users and also instance_roles as in
	// for a particular instance, he has access to as an direct allowed user
	// or as a role user
	sql, _, _ := goqu.From("instances").
		Where(
			goqu.Or(
				// Condition 1: User is in instance_users
				goqu.Ex{
					"instances.id":           goqu.I("instance_users.instance_id"),
					"instance_users.user_id": userID,
				},
				// Condition 2: User is part of a role that the instance allows access to
				goqu.Ex{
					"instances.id":           goqu.I("instance_roles.instance_id"),
					"instance_roles.role_id": goqu.I("role_users.role_id"),
					"role_users.user_id":     userID,
				},
			),
		).
		LeftJoin(
			goqu.T("instance_users"),
			goqu.On(goqu.Ex{"instances.id": goqu.I("instance_users.instance_id")}),
		).
		LeftJoin(
			goqu.T("instance_roles"),
			goqu.On(goqu.Ex{"instances.id": goqu.I("instance_roles.instance_id")}),
		).
		LeftJoin(
			goqu.T("role_users"),
			goqu.On(goqu.Ex{"instance_roles.role_id": goqu.I("role_users.role_id")}),
		).
		Select(
			"instances.id",
			"instances.name",
			"instances.host_address",
			"instances.status",
			"instances.description",
		).
		ToSQL()
	rows, err := h.DB.Query(c.Request().Context(), sql)
	if err != nil {
		return c.JSON(400, helper.ErrorMessage("Failed to fetch instances: ", err))
	}
	defer rows.Close()

	var instances []Instances = []Instances{}
	for rows.Next() {
		var instance Instances
		if err := rows.Scan(&instance.ID, &instance.Name, &instance.HostAddress, &instance.Status, &instance.CreatedBy); err != nil {
			return c.JSON(http.StatusInternalServerError, helper.ErrorMessage("Failed to fetch instance", err))
		}
		instances = append(instances, instance)
	}

	return c.JSON(http.StatusOK, instances)
}

func (h *instanceHandler) createInstance(c echo.Context) error {

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*authentication.JwtCustomClaims)
	var instance createInstanceStruct
	err := c.Bind(&instance)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Invalid request body", err))
	}

	if err := h.Validator.Struct(instance); err != nil {
		return c.JSON(400, helper.ErrorMessage("Input Validation failed", err))
	}
	sql, _, _ := goqu.Insert("instances").Rows(
		goqu.Record{
			"name":         instance.Name,
			"description":  instance.Description,
			"status":       instance.Status,
			"host_address": instance.HostAddress,
			"created_by":   claims.Id,
		},
	).Returning("id", "client_id").ToSQL()
	row := h.DB.QueryRow(c.Request().Context(), sql)
	var id int64
	var clientID string
	if err := row.Scan(&id, &clientID); err != nil {
		return c.JSON(500, helper.ErrorMessage("Failed to insert instance and retrieve ID", err))
	}

	if err != nil {
		return c.JSON(400, helper.ErrorMessage("Failed to create instance", err))
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message":   "Successfully created instance, Please save the client_id somewhere to configure the instance.",
		"status":    "success",
		"id":        id,
		"client_id": clientID,
	})
}

func (h *instanceHandler) deleteInstance(c echo.Context) error {
	var instanceID ParamsIDStruct
	err := (&echo.DefaultBinder{}).BindPathParams(c, &instanceID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Invalid path parameters", err))
	}
	if err := h.Validator.Struct(instanceID); err != nil {
		return c.JSON(400, helper.ErrorMessage("Input path parameters Validation failed", err))
	}
	sql, _, _ := goqu.From("instances").Where(goqu.C("id").Eq(instanceID.ID)).Delete().ToSQL()
	_, err = h.DB.Exec(c.Request().Context(), sql)
	if err != nil {
		return c.JSON(400, helper.ErrorMessage("Failed to delete instance", err))
	}
	return c.JSON(200, echo.Map{
		"status":  "success",
		"message": "Instance deleted successfully",
	})
}

func CensorClientID(clientID string) string {
	if len(clientID) < 19 {
		return "XXXX-HIDDENXXX"
	}
	return clientID[:14] + "XXXX-HIDDENXXX"
}

func (h *instanceHandler) getInstance(c echo.Context) error {
	var instanceID ParamsIDStruct
	err := (&echo.DefaultBinder{}).BindPathParams(c, &instanceID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Invalid path parameters", err))
	}
	if err := h.Validator.Struct(instanceID); err != nil {
		return c.JSON(400, helper.ErrorMessage("Input path parameters Validation failed", err))
	}
	sql, _, _ := goqu.From("instances").Join(
		goqu.T("users"),
		goqu.On(goqu.Ex{"instances.created_by": goqu.I("users.id")}),
	).Select(
		goqu.I("instances.id"), // Correctly quote table and column separately
		goqu.I("instances.name"),
		goqu.I("instances.host_address"),
		goqu.I("instances.status"),
		goqu.I("instances.description"),
		goqu.I("users.username").As("created_by"),
		goqu.I("instances.client_id"),
	).Where(
		goqu.I("instances.id").Eq(instanceID.ID), // Correctly quote table and column separately
	).ToSQL()

	row := h.DB.QueryRow(c.Request().Context(), sql)
	var instance Instance
	if err := row.Scan(&instance.ID, &instance.Name, &instance.HostAddress, &instance.Status, &instance.Description, &instance.CreatedBy, &instance.ClientID); err != nil {
		return c.JSON(500, helper.ErrorMessage("Failed to scan host user", err))
	}
	instance.ClientID = CensorClientID(instance.ClientID)
	hostUsersSql, _, _ := goqu.From("instance_host_users").
		Select("username").
		Where(goqu.Ex{"instance_id": instanceID.ID}).
		ToSQL()
	hostUsersRows, err := h.DB.Query(c.Request().Context(), hostUsersSql)
	if err != nil {
		return c.JSON(500, helper.ErrorMessage("Failed to fetch instance host users", err))
	}
	defer hostUsersRows.Close()
	var hostUsers []string
	for hostUsersRows.Next() {
		var username string
		if err := hostUsersRows.Scan(&username); err != nil {
			return c.JSON(500, helper.ErrorMessage("Failed to scan host user", err))
		}
		hostUsers = append(hostUsers, username)
	}
	instance.HostUsers = hostUsers

	return c.JSON(http.StatusOK, instance)
}

func (h *instanceHandler) getInstanceUsers(c echo.Context) error {
	var instanceID ParamsIDStruct
	err := (&echo.DefaultBinder{}).BindPathParams(c, &instanceID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Invalid path parameters", err))
	}
	if err := h.Validator.Struct(instanceID); err != nil {
		return c.JSON(400, helper.ErrorMessage("Input path parameters Validation failed", err))
	}
	sql, _, _ := goqu.From("instance_users").Join(
		goqu.T("users"),
		goqu.On(goqu.Ex{"instance_users.user_id": goqu.I("users.id")}),
	).Select(
		goqu.I("users.id"),
		goqu.I("users.name"),
		goqu.I("users.username"),
		goqu.I("users.role"),
		goqu.I("instance_users.instance_host_username"),
	).Where(
		goqu.I("instance_users.instance_id").Eq(instanceID.ID),
	).ToSQL()

	rows, err := h.DB.Query(c.Request().Context(), sql)
	if err != nil {
		return c.JSON(400, helper.ErrorMessage("Failed to fetch instance users", err))
	}
	defer rows.Close()

	var users []InstanceUser = []InstanceUser{}
	for rows.Next() {
		var user InstanceUser
		if err := rows.Scan(&user.Id, &user.Name, &user.Username, &user.Role, &user.HostUsername); err != nil {
			return c.JSON(http.StatusInternalServerError, helper.ErrorMessage("Failed to scan user", err))
		}
		users = append(users, user)
	}

	return c.JSON(http.StatusOK, users)
}

func (h *instanceHandler) editInstance(c echo.Context) error {
	var instanceID ParamsIDStruct
	err := (&echo.DefaultBinder{}).BindPathParams(c, &instanceID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Invalid path parameters", err))
	}
	if err := h.Validator.Struct(instanceID); err != nil {
		return c.JSON(400, helper.ErrorMessage("Input path parameters Validation failed", err))
	}

	var instanceData createInstanceStruct
	err = c.Bind(&instanceData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Invalid request body", err))
	}
	if err := h.Validator.Struct(instanceData); err != nil {
		return c.JSON(400, helper.ErrorMessage("Invalid request body", err))
	}
	sql, _, _ := goqu.From("instances").
		Where(goqu.C("id").Eq(instanceID.ID)).Update().Set(goqu.Record{
		"name":         instanceData.Name,
		"description":  instanceData.Description,
		"host_address": instanceData.HostAddress,
		"status":       instanceData.Status,
	}).Returning(goqu.I("id"), goqu.I("client_id")).ToSQL()
	row := h.DB.QueryRow(c.Request().Context(), sql)
	var id int64
	var clientID string
	if err := row.Scan(&id, &clientID); err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorMessage("Failed to update instance and retrieve ID", err))
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message":   "Instance updated successfully",
		"status":    "success",
		"id":        id,
		"client_id": clientID,
	})
}
