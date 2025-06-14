package instances

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

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
		"message":   "Successfully created instance",
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
	id := c.Param("id")
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
		goqu.I("instances.id").Eq(id), // Correctly quote table and column separately
	).ToSQL()

	row := h.DB.QueryRow(context.Background(), sql)
	var instance Instance
	if err := row.Scan(&instance.ID, &instance.Name, &instance.HostAddress, &instance.Status, &instance.Description, &instance.CreatedBy, &instance.ClientID); err != nil {
		return c.JSON(500, helper.ErrorMessage("Failed to scan host user", err))
	}
	instance.ClientID = CensorClientID(instance.ClientID)
	hostUsersSql, _, _ := goqu.From("instance_host_users").
		Select("username").
		Where(goqu.Ex{"instance_id": id}).
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

func (h *instanceHandler) deleteInstanceUsers(c echo.Context) error {
	var instance ParamsIDStruct
	var user ParamsIDStruct
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
	// sql, _, _ := goqu.From("instance_users").Where(
	// 	goqu.Ex{"instance_id": instance.ID, "user_id": user.ID},
	// ).Delete().ToSQL()
	return c.JSON(http.StatusOK, echo.Map{
		"message":  "Successfully deleted instance user: " + strconv.FormatUint(user.ID, 10),
		"instance": instance.ID,
	})
}

type InstanceRole struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	HostUsername string `json:"host_username"`
}

func (h *instanceHandler) getInstanceRoles(c echo.Context) error {
	var instance ParamsIDStruct

	// Bind path parameters using DefaultBinder
	err := (&echo.DefaultBinder{}).BindPathParams(c, &instance)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid path parameters",
		})
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

	rows, err := h.DB.Query(context.Background(), sql)
	if err != nil {
		return c.JSON(400, echo.Map{
			"error": "Failed to fetch instance roles: " + err.Error(),
		})
	}
	defer rows.Close()

	var roles []InstanceRole = []InstanceRole{}
	for rows.Next() {
		var role InstanceRole
		if err := rows.Scan(&role.Id, &role.Name, &role.HostUsername); err != nil {
			log.Fatalf("Failed to scan role: %v", err)
		}
		roles = append(roles, role)
	}

	return c.JSON(http.StatusOK, roles)
}

type IDUsernameStruct struct {
	ID       uint64 `json:"id"`            // Match the path parameter name
	Username string `json:"host_username"` // Match the body parameter name
}

func (h *instanceHandler) addInstanceRoles(c echo.Context) error {
	var instance ParamsIDStruct
	err := (&echo.DefaultBinder{}).BindPathParams(c, &instance)
	if err != nil {
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
	err = h.DB.QueryRow(context.Background(), sql).Scan(&count)
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
	_, err = h.DB.Exec(context.Background(), sql)
	if err != nil {
		return c.JSON(400, echo.Map{
			"error":             "Failed to add instance role",
			"error_description": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": fmt.Sprintf("Successfully added instance role with ID %d and username %s", role.ID, role.Username),
		"status":  "success",
	})

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
