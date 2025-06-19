package instances

import (
	"fmt"
	"net/http"

	"github.com/Utkar5hM/Execstasy/api/utils/helper"
	"github.com/doug-martin/goqu/v9"
	"github.com/labstack/echo/v4"
)

func (h *instanceHandler) addUserInstanceAccess(c echo.Context) error {
	var instanceID ParamsIDStruct
	err := (&echo.DefaultBinder{}).BindPathParams(c, &instanceID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Invalid path parameters", err))
	}
	if err := h.Validator.Struct(instanceID); err != nil {
		return c.JSON(400, helper.ErrorMessage("Input Validation failed", err))
	}
	var addUserInstanceStruct InstanceUsernameStruct
	err = c.Bind(&addUserInstanceStruct)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid request body")
	}
	if err := h.Validator.Struct(addUserInstanceStruct); err != nil {
		return c.JSON(400, helper.ErrorMessage("Input Validation failed", err))
	}
	instanceId := instanceID.ID
	sql, _, _ := goqu.From("users").Where(goqu.C("username").Eq(addUserInstanceStruct.Username)).Select("id").ToSQL()
	row := h.DB.QueryRow(c.Request().Context(), sql)
	var userID string
	err = row.Scan(&userID)
	if err != nil {
		return c.JSON(400, helper.ErrorMessage("Failed to add access for user to Instance.", err.Error()))
	}
	sql, _, _ = goqu.Insert("instance_users").Rows(
		goqu.Record{
			"instance_id":            instanceId,
			"user_id":                userID,
			"instance_host_username": addUserInstanceStruct.HostUsername,
		},
	).ToSQL()
	_, err = h.DB.Exec(c.Request().Context(), sql)
	if err != nil {
		return c.JSON(400, helper.ErrorMessage("Failed to add access for user to Instance.", err.Error()))
	}
	return c.JSON(200, echo.Map{
		"status":  "success",
		"message": "User access added to instance successfully.",
	})
}

func (h *instanceHandler) deleteUserInstanceAccess(c echo.Context) error {
	var instanceID ParamsIDStruct
	err := (&echo.DefaultBinder{}).BindPathParams(c, &instanceID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Invalid path parameters", err))
	}
	if err := h.Validator.Struct(instanceID); err != nil {
		return c.JSON(400, helper.ErrorMessage("Input Validation failed", err))
	}
	var user InstanceUsernameStruct
	err = c.Bind(&user)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid request body")
	}
	if err := h.Validator.Struct(user); err != nil {
		return c.JSON(400, helper.ErrorMessage("Input Validation failed", err))
	}
	// Check if the user exists in the users table
	sql, _, _ := goqu.From("users").Where(
		goqu.Ex{"username": user.Username},
	).Select(goqu.COUNT("*")).ToSQL()
	var count int64
	err = h.DB.QueryRow(c.Request().Context(), sql).Scan(&count)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorMessage("Failed to check existing user: ", err))
	}
	if count == 0 {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("User does not exist", fmt.Sprintf("User with username %s does not exist", user.Username)))
	}
	// Check if the user is associated with the instance
	sql, _, _ = goqu.From("instance_users").Where(
		goqu.Ex{"instance_id": instanceID.ID, "user_id": goqu.I("users.id"), "instance_host_username": user.HostUsername},
	).Join(
		goqu.T("users"),
		goqu.On(goqu.Ex{"instance_users.user_id": goqu.I("users.id")}),
	).Select(goqu.COUNT("*")).ToSQL()
	var userCount int64
	err = h.DB.QueryRow(c.Request().Context(), sql).Scan(&userCount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorMessage("Failed to check existing instance user: ", err))
	}
	if userCount == 0 {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("User is not associated with the instance", nil))
	}
	sql, _, _ = goqu.From("users").
		Where(goqu.Ex{"username": user.Username}).
		Select("id").
		ToSQL()
	var userID int64
	err = h.DB.QueryRow(c.Request().Context(), sql).Scan(&userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorMessage("Failed to get user ID: ", err))
	}
	sql, _, _ = goqu.Delete("instance_users").
		Where(goqu.Ex{
			"instance_id":            instanceID.ID,
			"user_id":                userID,
			"instance_host_username": user.HostUsername,
		}).
		ToSQL()
	_, err = h.DB.Exec(c.Request().Context(), sql)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorMessage("Failed to delete instance user: ", err))
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Successfully deleted instance user.",
		"status":  "success",
	})
}
