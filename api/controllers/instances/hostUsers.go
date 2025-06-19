package instances

import (
	"net/http"

	"github.com/Utkar5hM/Execstasy/api/utils/helper"
	"github.com/doug-martin/goqu/v9"
	"github.com/labstack/echo/v4"
)

func (h *instanceHandler) addInstanceHostUser(c echo.Context) error {
	var instanceID ParamsIDStruct
	err := (&echo.DefaultBinder{}).BindPathParams(c, &instanceID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Invalid path parameters", err))
	}
	if err := h.Validator.Struct(instanceID); err != nil {
		return c.JSON(400, helper.ErrorMessage("Input path parameters Validation failed", err))
	}
	var hostUser hostUsernameStruct
	err = c.Bind(&hostUser)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Invalid request body", err))
	}
	if err := h.Validator.Struct(hostUser); err != nil {
		return c.JSON(400, helper.ErrorMessage("Input Validation failed", err))
	}
	sql, _, _ := goqu.Insert("instance_host_users").Rows(
		goqu.Record{
			"instance_id": instanceID.ID,
			"username":    hostUser.Username,
		},
	).ToSQL()
	_, err = h.DB.Exec(c.Request().Context(), sql)
	if err != nil {
		return c.JSON(400, helper.ErrorMessage("Failed to add host user", err))
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Successfully added host user",
		"status":  "success",
	})

}

func (h *instanceHandler) deleteInstanceHostUser(c echo.Context) error {
	var instanceID ParamsIDStruct
	err := (&echo.DefaultBinder{}).BindPathParams(c, &instanceID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Invalid path parameters", err))
	}
	if err := h.Validator.Struct(instanceID); err != nil {
		return c.JSON(400, helper.ErrorMessage("Input path parameters Validation failed", err))
	}
	var hostUser hostUsernameStruct
	err = c.Bind(&hostUser)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Invalid request body", err))
	}
	if err := h.Validator.Struct(hostUser); err != nil {
		return c.JSON(400, helper.ErrorMessage("Input Validation failed", err))
	}
	sql, _, _ := goqu.From("instance_host_users").Where(goqu.C("instance_id").Eq(instanceID.ID), goqu.C("username").Eq(hostUser.Username)).Delete().ToSQL()
	_, err = h.DB.Exec(c.Request().Context(), sql)
	if err != nil {
		return c.JSON(400, helper.ErrorMessage("Failed to delete host user", err))
	}
	return c.JSON(200, echo.Map{
		"status":  "success",
		"message": "Successfully deleted host user",
	})
}
