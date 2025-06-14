package instances

import (
	"context"
	"net/http"

	"github.com/Utkar5hM/Execstasy/api/utils/helper"
	"github.com/doug-martin/goqu/v9"
	"github.com/labstack/echo/v4"
)

func (h *instanceHandler) addInstanceHostUser(c echo.Context) error {
	var instance ParamsIDStruct
	err := (&echo.DefaultBinder{}).BindPathParams(c, &instance)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid path parameters",
		})
	}
	var hostUser hostUsernameStruct
	err = c.Bind(&hostUser)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid request body",
		})
	}
	if hostUser.Username == "" {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Host username is required", nil))
	}
	sql, _, _ := goqu.Insert("instance_host_users").Rows(
		goqu.Record{
			"instance_id": instance.ID,
			"username":    hostUser.Username,
		},
	).ToSQL()
	_, err = h.DB.Exec(context.Background(), sql)
	if err != nil {
		return c.JSON(400, helper.ErrorMessage("Failed to add host user: ", err))
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Successfully added host user",
		"status":  "success",
	})

}

type hostUsernameStruct struct {
	Username string `json:"host_username"` // Match the path parameter name
}

func (h *instanceHandler) deleteInstanceHostUser(c echo.Context) error {
	var instance ParamsIDStruct
	err := (&echo.DefaultBinder{}).BindPathParams(c, &instance)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid path parameters",
		})
	}
	var hostUser hostUsernameStruct
	err = c.Bind(&hostUser)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid request body",
		})
	}
	if hostUser.Username == "" {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Host username is required", nil))
	}
	sql, _, _ := goqu.From("instance_host_users").Where(goqu.C("instance_id").Eq(instance.ID), goqu.C("username").Eq(hostUser.Username)).Delete().ToSQL()
	_, err = h.DB.Exec(context.Background(), sql)
	if err != nil {
		return c.JSON(400, helper.ErrorMessage("Failed to delete host user: ", err))
	}
	return c.JSON(200, echo.Map{
		"status": "success",
	})
}

// func (h *instanceHandler) getInstanceHostUsers(c echo.Context) error {
// 	var instance ParamsIDStruct
// 	err := (&echo.DefaultBinder{}).BindPathParams(c, &instance)
// 	if err != nil || instance.ID <= 0 {
// 		return c.JSON(http.StatusBadRequest, echo.Map{
// 			"error": "Invalid path parameters",
// 		})
// 	}
// 	sql, _, _ := goqu.From("instance_host_users").Where(
// 		goqu.Ex{"instance_id": instance.ID},
// 	).Select("username").ToSQL()
// 	rows, err := h.DB.Query(context.Background(), sql)
// 	if err != nil {
// 		return c.JSON(400, helper.ErrorMessage("Failed to fetch instance host users", err))
// 	}
// 	defer rows.Close()

// 	var hostUsers []string
// 	for rows.Next() {
// 		var username string
// 		if err := rows.Scan(&username); err != nil {
// 			return c.JSON(500, helper.ErrorMessage("Failed to scan host user", err))
// 		}
// 		hostUsers = append(hostUsers, username)
// 	}

// 	return c.JSON(http.StatusOK, hostUsers)
// }
