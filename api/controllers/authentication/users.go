package authentication

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Utkar5hM/Execstasy/api/utils/helper"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func (h *AuthHandler) getUsers(c echo.Context) error {
	sql, _, _ := goqu.From("users").
		Select("id", "username", "name", "email", "role").ToSQL()
	rows, err := h.DB.Query(c.Request().Context(), sql)
	if err != nil {
		return c.JSON(500, echo.Map{
			"error": "Failed to fetch users",
		})
	}
	defer rows.Close()
	var users []map[string]interface{}
	for rows.Next() {
		var id uint64
		var username, name, email, role string
		if err := rows.Scan(&id, &username, &name, &email, &role); err != nil {
			return c.JSON(500, echo.Map{
				"error": "Failed to scan user",
			})
		}
		users = append(users, map[string]interface{}{
			"id":       id,
			"username": username,
			"name":     name,
			"email":    email,
			"role":     role,
		})
	}
	return c.JSON(200, users)
}

func (h *AuthHandler) searchUsers(c echo.Context) error {
	// Get the search query from the request
	query := c.QueryParam("query")
	if query == "" {
		return c.JSON(400, echo.Map{
			"error": "Query parameter is required",
		})
	}

	// Build the SQL query with filtering
	sql, _, _ := goqu.From("users").
		Where(
			goqu.Or(
				goqu.I("username").ILike("%"+query+"%"),
				goqu.I("name").ILike("%"+query+"%"),
			),
		).
		Select("id", "username", "name", "email", "role").ToSQL()

	// Execute the query
	rows, err := h.DB.Query(c.Request().Context(), sql)
	if err != nil {
		return c.JSON(500, echo.Map{
			"error": "Failed to fetch users",
		})
	}
	defer rows.Close()

	// Parse the results
	var users []map[string]interface{} = make([]map[string]interface{}, 0)
	for rows.Next() {
		var id uint64
		var username, name, email, role string
		if err := rows.Scan(&id, &username, &name, &email, &role); err != nil {
			return c.JSON(500, echo.Map{
				"error": "Failed to scan user",
			})
		}
		users = append(users, map[string]interface{}{
			"id":       id,
			"username": username,
			"name":     name,
			"email":    email,
			"role":     role,
		})
	}

	// Return the filtered users
	return c.JSON(200, users)
}

func (h *AuthHandler) getMe(c echo.Context) error {
	// Get the user ID from the context
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtCustomClaims)
	// Build the SQL query to fetch the user by ID
	sql, _, _ := goqu.From("users").
		Where(goqu.I("id").Eq(claims.Id)).
		Select("id", "username", "name", "email", "role").ToSQL()

	// Execute the query
	row := h.DB.QueryRow(c.Request().Context(), sql)
	var id uint64
	var username, name, email, role string
	if err := row.Scan(&id, &username, &name, &email, &role); err != nil {
		return c.JSON(500, helper.ErrorMessage("Failed to fetch user information", err))
	}

	// Return the user information
	return c.JSON(200, map[string]interface{}{
		"id":       id,
		"username": username,
		"name":     name,
		"email":    email,
		"role":     role,
	})
}

type updateMeRequest struct {
	Username string `json:"username" validate:"required,alphanumunicode,min=2,max=50"`
	Name     string `json:"name" validate:"required,min=1,max=100,usernameregex"`
}

func (h *AuthHandler) updateMe(c echo.Context) error {
	// Get the user ID from the context
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtCustomClaims)
	var req updateMeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, helper.ErrorMessage("Invalid request body", err))
	}
	if err := h.Validator.Struct(req); err != nil {
		return c.JSON(400, helper.ErrorMessage("Validation failed", err))
	}
	if req.Username == "" || req.Name == "" {
		return c.JSON(400, helper.ErrorMessage("Username and Name cannot be empty", nil))
	}
	// Build the SQL query to update the user
	dialect := goqu.Dialect("postgres")
	sql, args, _ := dialect.Update("users").
		Where(goqu.I("id").Eq(claims.Id)).
		Set(
			goqu.Record{
				"username": req.Username,
				"name":     req.Name,
			},
		).Returning(
		"id", "username", "name", "role", "email",
	).Prepared(true).ToSQL()
	// Execute the update query
	fmt.Println(sql)
	row := h.DB.QueryRow(c.Request().Context(), sql, args...)
	var id uint64
	var username, name, role, email string
	if err := row.Scan(&id, &username, &name, &role, &email); err != nil {
		return c.JSON(500, helper.ErrorMessage("Failed to update user information", err))
	}
	// Regenerate JWT Token
	newClaims := &JwtCustomClaims{
		name,
		username,
		role,
		int(id),
		email,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}
	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims).SignedString([]byte(h.Config.JWT_SECRET))
	if err != nil {
		return err
	}

	// Set JWT token as cookie
	cookie := new(http.Cookie)
	cookie.Name = "jwt"
	cookie.Value = tokenString
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.Path = "/"
	c.SetCookie(cookie)
	// Return success message
	return c.JSON(200, echo.Map{
		"message": "User information updated successfully",
		"status":  "success",
	})
}

type JSONIDStruct struct {
	ID uint64 `json:"id" validate:"required,gt=0"` // Match the path parameter name
}

func (h *AuthHandler) updateRole(c echo.Context) error {
	// Get the user ID from the context
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtCustomClaims)
	var userID JSONIDStruct
	if err := c.Bind(&userID); err != nil {
		return c.JSON(400, helper.ErrorMessage("Invalid request body", err))
	}
	if err := h.Validator.Struct(userID); err != nil {
		return c.JSON(400, helper.ErrorMessage("Validation failed", err))
	}
	if userID.ID == 0 {
		return c.JSON(400, helper.ErrorMessage("User ID cannot be empty", nil))
	}
	if userID.ID == uint64(claims.Id) {
		return c.JSON(400, helper.ErrorMessage("You cannot update your own role", nil))
	}
	//  Get the user's current role
	sql, _, _ := goqu.From("users").
		Where(goqu.I("id").Eq(userID.ID)).
		Select("role").ToSQL()
	row := h.DB.QueryRow(c.Request().Context(), sql)
	var currentRole string
	if err := row.Scan(&currentRole); err != nil {
		return c.JSON(500, helper.ErrorMessage("Failed to fetch user role", err))
	}
	var newRole string
	if currentRole == "admin" {
		newRole = "user" // If current role is admin, change to user
	} else {
		newRole = "admin" // Otherwise, change to admin
	}
	sql, _, _ = goqu.Update("users").
		Where(goqu.I("id").Eq(userID.ID)).
		Set(goqu.Record{"role": newRole}).ToSQL()
	// Execute the update query
	_, err := h.DB.Exec(c.Request().Context(), sql)
	if err != nil {
		return c.JSON(500, helper.ErrorMessage("Failed to update user role", err))
	}
	// Return success message
	return c.JSON(200, echo.Map{
		"message": "User role updated successfully to " + newRole,
		"status":  "success",
	})
}
