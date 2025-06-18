package instances

import (
	"context"
	"net/http"

	"github.com/Utkar5hM/Execstasy/api/utils/helper"
	"github.com/doug-martin/goqu/v9"
	"github.com/labstack/echo/v4"
)

type AddUserInstanceStruct struct {
	HostUsername string `json:"host_username"` // Match the path parameter name
	Username     string `json:"username"`      // Match the body parameter name
}

func (h *instanceHandler) addUserInstanceAccess(c echo.Context) error {
	var addUserInstanceStruct AddUserInstanceStruct
	err := c.Bind(&addUserInstanceStruct)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid request body")
	}
	if addUserInstanceStruct.HostUsername == "" || addUserInstanceStruct.Username == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid role ID or username",
		})
	}
	instanceId := c.Param("id")
	sql, _, _ := goqu.From("users").Where(goqu.C("username").Eq(addUserInstanceStruct.Username)).Select("id").ToSQL()
	row := h.DB.QueryRow(context.Background(), sql)
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
	_, err = h.DB.Exec(context.Background(), sql)
	if err != nil {
		return c.JSON(400, helper.ErrorMessage("Failed to add access for user to Instance.", err.Error()))
	}
	return c.JSON(200, echo.Map{
		"status":  "success",
		"message": "User access added to instance successfully.",
	})
}

func (h *instanceHandler) deleteUserInstanceAccess(c echo.Context) error {
	instanceId := c.Param("id")
	username := c.FormValue("username")
	hostUsername := c.FormValue("host_username")
	sql, _, _ := goqu.From("users").Where(goqu.C("username").Eq(username)).Select("id").ToSQL()
	row := h.DB.QueryRow(context.Background(), sql)
	var userID string
	err := row.Scan(&userID)
	if err != nil {
		return c.JSON(400, echo.Map{
			"error":   "Failed to delete instance access to the user.",
			"message": err.Error(),
			"status":  "error",
		})
	}
	sql, _, _ = goqu.From("instance_users").Where(goqu.C("instance_id").Eq(instanceId), goqu.C("user_id").Eq(userID), goqu.C("instance_host_username").Eq(hostUsername)).Delete().ToSQL()
	_, err = h.DB.Exec(context.Background(), sql)
	if err != nil {
		return c.JSON(400, echo.Map{
			"error":   "Failed to delete instance access to the user.",
			"message": err.Error(),
			"status":  "error",
		})
	}
	return c.JSON(200, echo.Map{
		"status": "success",
	})
}
