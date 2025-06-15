package instances

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Utkar5hM/Execstasy/api/controllers/authentication"
	"github.com/Utkar5hM/Execstasy/api/utils/helper"
	"github.com/doug-martin/goqu/v9"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Instances struct {
	ID          int    `db:"id"`
	Name        string `db:"name"`
	HostAddress string `db:"host_address"`
	Status      string `db:"status"`
	CreatedBy   string `db:"created_by"`
}

type Instance struct {
	ID          int      `db:"id"`
	Name        string   `db:"name"`
	HostAddress string   `db:"host_address"`
	Status      string   `db:"status"`
	CreatedBy   string   `db:"created_by"`
	Description string   `db:"description"`
	HostUsers   []string `db:"host_users"` // Assuming this is a slice of usernames
	ClientID    string   `db:"client_id"`  // Assuming this is a string, adjust as necessary
}

type ParamsIDStruct struct {
	ID uint64 `param:"id"` // Match the path parameter name
}

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
	rows, err := h.DB.Query(context.Background(), sql)
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

type createInstanceStruct struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	HostAddress string `json:"HostAddress"`
	Status      string `json:"status"`
}

func (h *instanceHandler) createInstance(c echo.Context) error {

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*authentication.JwtCustomClaims)
	var instance createInstanceStruct
	err := c.Bind(&instance)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Invalid request body", err))
	}
	if instance.Name == "" ||
		(instance.Status != "" && instance.Status != "active" && instance.Status != "disabled") {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Name and description are required fields", nil))
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
	row := h.DB.QueryRow(context.Background(), sql)
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
	id := c.Param("id")
	sql, _, _ := goqu.From("instances").Where(goqu.C("id").Eq(id)).Delete().ToSQL()
	_, err := h.DB.Exec(context.Background(), sql)
	if err != nil {
		return c.JSON(400, helper.ErrorMessage("Failed to delete instance", err))
	}
	return c.JSON(200, echo.Map{
		"status": "success",
	})
}

func (h *instanceHandler) setStatusInstance(c echo.Context) error {
	id := c.Param("id")
	status := c.FormValue("status")
	sql, _, _ := goqu.From("instances").Where(goqu.C("id").Eq(id)).Update().Set(goqu.Record{"status": status}).ToSQL()
	_, err := h.DB.Exec(context.Background(), sql)
	if err != nil {
		return c.JSON(400, helper.ErrorMessage(fmt.Sprintf("Failed to set instance status to %s.", status), err))
	}
	return c.JSON(200, echo.Map{
		"status": "success",
	})
}

func (h *instanceHandler) isAdminOrCreatorMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*authentication.JwtCustomClaims)
		if claims.Role != "admin" {
			id := c.Param("id")
			sql, _, _ := goqu.From("instances").Where(goqu.C("id").Eq(id)).Select("created_by").ToSQL()
			row := h.DB.QueryRow(context.Background(), sql)
			var instance_created_by int
			err := row.Scan(&instance_created_by)
			if err != nil {
				return c.JSON(400, echo.Map{
					"message": "Failed to Authroize your action.",
					"error":   err.Error(),
				})
			}
			if claims.Id != instance_created_by {
				return c.JSON(403, echo.Map{
					"message": "You are not authorized to disable this instance",
				})
			}
		}
		return next(c)
	}
}

func CensorClientID(clientID string) string {
	if len(clientID) < 19 {
		return "XXXX-HIDDENXXX"
	}
	return clientID[:14] + "XXXX-HIDDENXXX"
}

func (h *instanceHandler) getInstance(c echo.Context) error {
	var instanceID ParamsIDStruct

	// Bind path parameters using DefaultBinder
	err := (&echo.DefaultBinder{}).BindPathParams(c, &instanceID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Invalid path parameters", err))
	}
	if instanceID.ID == 0 {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Instance ID is required", nil))
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

	row := h.DB.QueryRow(context.Background(), sql)
	var instance Instance
	if err := row.Scan(&instance.ID, &instance.Name, &instance.HostAddress, &instance.Status, &instance.Description, &instance.CreatedBy, &instance.ClientID); err != nil {
		return c.JSON(500, helper.ErrorMessage("Failed to scan host user", err))
	}
	instance.ClientID = CensorClientID(instance.ClientID)
	hostUsersSql, _, _ := goqu.From("instance_host_users").
		Select("username").
		Where(goqu.Ex{"instance_id": instanceID.ID}).
		ToSQL()
	hostUsersRows, err := h.DB.Query(context.Background(), hostUsersSql)
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

type InstanceUser struct {
	Id           int    `json:"id"`
	Username     string `json:"username"`
	Role         string `json:"role"`
	Name         string `json:"name"`
	HostUsername string `json:"host_username"`
}

func (h *instanceHandler) getInstanceUsers(c echo.Context) error {
	var instance ParamsIDStruct

	// Bind path parameters using DefaultBinder
	err := (&echo.DefaultBinder{}).BindPathParams(c, &instance)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid path parameters",
		})
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
		goqu.I("instance_users.instance_id").Eq(instance.ID),
	).ToSQL()

	rows, err := h.DB.Query(context.Background(), sql)
	if err != nil {
		return c.JSON(400, echo.Map{
			"error": "Failed to fetch instance users: " + err.Error(),
		})
	}
	defer rows.Close()

	var users []InstanceUser = []InstanceUser{}
	for rows.Next() {
		var user InstanceUser
		if err := rows.Scan(&user.Id, &user.Name, &user.Username, &user.Role, &user.HostUsername); err != nil {
			log.Fatalf("Failed to scan user: %v", err)
		}
		users = append(users, user)
	}

	return c.JSON(http.StatusOK, users)
}

type InstanceUsernameStruct struct {
	Username     string `json:"username"`      // Match the path parameter name
	HostUsername string `json:"host_username"` // Match the body parameter name
}

func (h *instanceHandler) deleteInstanceUsers(c echo.Context) error {
	var instance ParamsIDStruct
	var user InstanceUsernameStruct
	err := (&echo.DefaultBinder{}).BindPathParams(c, &instance)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid path parameters",
		})
	}

	err = (&echo.DefaultBinder{}).BindBody(c, &user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid body parameters",
		})
	}
	if user.Username == "" || user.HostUsername == "" {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Username and host username are required", nil))
	}
	// Check if the user exists in the users table
	sql, _, _ := goqu.From("users").Where(
		goqu.Ex{"username": user.Username},
	).Select(goqu.COUNT("*")).ToSQL()
	var count int64
	err = h.DB.QueryRow(context.Background(), sql).Scan(&count)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorMessage("Failed to check existing user: ", err))
	}
	if count == 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":             "User does not exist",
			"error_description": fmt.Sprintf("User with username %s does not exist", user.Username),
		})
	}
	// Check if the user is associated with the instance
	sql, _, _ = goqu.From("instance_users").Where(
		goqu.Ex{"instance_id": instance.ID, "user_id": goqu.I("users.id"), "instance_host_username": user.HostUsername},
	).Join(
		goqu.T("users"),
		goqu.On(goqu.Ex{"instance_users.user_id": goqu.I("users.id")}),
	).Select(goqu.COUNT("*")).ToSQL()
	var userCount int64
	err = h.DB.QueryRow(context.Background(), sql).Scan(&userCount)
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
	err = h.DB.QueryRow(context.Background(), sql).Scan(&userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorMessage("Failed to get user ID: ", err))
	}
	sql, _, _ = goqu.Delete("instance_users").
		Where(goqu.Ex{
			"instance_id":            instance.ID,
			"user_id":                userID,
			"instance_host_username": user.HostUsername,
		}).
		ToSQL()
	_, err = h.DB.Exec(context.Background(), sql)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorMessage("Failed to delete instance user: ", err))
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Successfully deleted instance user.",
		"status":  "success",
	})
}

type InstanceRole struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	HostUsername string `json:"host_username"`
}

type IDUsernameStruct struct {
	ID       uint64 `json:"id"`            // Match the path parameter name
	Username string `json:"host_username"` // Match the body parameter name
}

func (h *instanceHandler) editInstance(c echo.Context) error {
	var instance ParamsIDStruct
	err := (&echo.DefaultBinder{}).BindPathParams(c, &instance)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid path parameters",
		})
	}

	var instanceData Instance
	err = c.Bind(&instanceData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid request body",
		})
	}
	if instanceData.Name == "" ||
		(instanceData.Status != "" && instanceData.Status != "active" && instanceData.Status != "disabled") {
		return c.JSON(http.StatusBadRequest, helper.ErrorMessage("Name and description are required fields", nil))
	}
	sql, _, _ := goqu.From("instances").
		Where(goqu.C("id").Eq(instance.ID)).Update().Set(goqu.Record{
		"name":         instanceData.Name,
		"description":  instanceData.Description,
		"host_address": instanceData.HostAddress,
		"status":       instanceData.Status,
	}).Returning(goqu.I("id"), goqu.I("client_id")).ToSQL()
	row := h.DB.QueryRow(context.Background(), sql)
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
