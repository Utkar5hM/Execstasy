package instances

import (
	"crypto/rand"
	"errors"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Utkar5hM/Execstasy/api/controllers/authentication"
	"github.com/Utkar5hM/Execstasy/api/utils/helper"
	"github.com/doug-martin/goqu/v9"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	userCodeLength = 8
	charset        = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func generateUserCode() (string, error) {
	code := make([]byte, userCodeLength)
	for i := range code {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		code[i] = charset[num.Int64()]
	}
	// Optionally, insert a hyphen in the middle for readability
	return string(code), nil
}

func (h *instanceHandler) deviceAuthorization(c echo.Context) error {
	var authReq deviceAuthorizationRequest
	if err := c.Bind(&authReq); err != nil {
		return c.JSON(http.StatusBadRequest, helper.AuthErrorMessage("invalid_request", "Bad request format"))
	}
	if err := h.Validator.Struct(authReq); err != nil {
		return c.JSON(400, helper.AuthErrorMessage("invalid_request", err))
	}
	clientIP := c.RealIP()

	sql, _, _ := goqu.From("instances").Where(goqu.C("client_id").Eq(authReq.ClientId)).Select("id").ToSQL()

	row := h.DB.QueryRow(c.Request().Context(), sql)
	var instanceId int
	err := row.Scan(&instanceId)
	if err != nil {
		return c.JSON(400, helper.AuthErrorMessage("Invalid client_id", "Instance with the specified client_id does not exist"))
	}

	deviceCode := uuid.New().String()
	userCode, err := generateUserCode()

	if err != nil {
		return c.JSON(400, helper.AuthErrorMessage("Failed to generate user code", err))
	}

	expiration := time.Now().Add(10 * time.Minute).Unix()

	data := map[string]interface{}{
		"client_id":   authReq.ClientId,
		"scope":       authReq.Scope,
		"user_code":   userCode,
		"device_code": deviceCode,
		"expires_at":  expiration,
		"status":      "pending",
		"clientIP":    clientIP,
		"instance_id": instanceId,
	}
	_, err = h.RDB.HSet(c.Request().Context(), deviceCode, data).Result()
	if err != nil {
		return c.JSON(400, helper.AuthErrorMessage("Failed to store device code", err))
	}
	_, err = h.RDB.HSet(c.Request().Context(), userCode, data).Result()
	if err != nil {
		return c.JSON(400, helper.AuthErrorMessage("Failed to store user code", err))
	}
	_, err = h.RDB.Expire(c.Request().Context(), deviceCode, 20*time.Minute).Result()
	if err != nil {
		return c.JSON(400, helper.AuthErrorMessage("Failed to set expiration for device code", err))
	}
	_, err = h.RDB.Expire(c.Request().Context(), userCode, 20*time.Minute).Result()
	if err != nil {
		return c.JSON(400, echo.Map{
			"error":   "Failed to set expiration for device code",
			"message": err.Error(),
		})
	}

	return c.JSON(200, echo.Map{
		"device_code":               deviceCode,
		"user_code":                 string(userCode[:4]) + "-" + string(userCode[4:]),
		"verification_uri":          c.Scheme() + "://" + c.Request().Host + "/oauth",
		"verification_uri_complete": c.Scheme() + "://" + c.Request().Host + "/oauth?user_code=" + userCode,
		"expires_in":                (10 * time.Minute) / time.Second,
		"interval":                  5,
	})
}

func (h *instanceHandler) token(c echo.Context) error {
	var tokenReq tokenRequest
	if err := c.Bind(&tokenReq); err != nil {
		return c.JSON(http.StatusBadRequest, helper.AuthErrorMessage("invalid_request", "Bad request format"))
	}
	if err := h.Validator.Struct(tokenReq); err != nil {
		return c.JSON(400, helper.AuthErrorMessage("invalid_request", err))
	}

	value, err := h.RDB.HGetAll(c.Request().Context(), tokenReq.DeviceCode).Result()
	if err != nil || len(value) == 0 {
		return c.JSON(http.StatusBadRequest, helper.AuthErrorMessage("invalid_grant", "Device code does not exist or has expired"))
	}

	if value["client_id"] != tokenReq.ClientId {
		return c.JSON(http.StatusUnauthorized, helper.AuthErrorMessage("invalid_client", "Client ID does not match the device code"))
	}

	switch value["status"] {
	case "pending":
		return c.JSON(http.StatusBadRequest, helper.AuthErrorMessage("authorization_pending", "User has not yet completed the authorization"))
	case "denied":
		return c.JSON(http.StatusBadRequest, helper.AuthErrorMessage("access_denied", "User has denied the authorization"))
	case "approved":
		expiresAt, err := strconv.ParseInt(value["expires_at"], 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.AuthErrorMessage("invalid_grant", "Invalid expiration time format"))
		}
		if expiresAt < time.Now().Unix() {
			return c.JSON(http.StatusBadRequest, helper.AuthErrorMessage("expired_token", "Device code has expired"))
		}

		accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub":   value["client_id"],
			"scope": value["scope"],
		})
		token, err := accessToken.SignedString([]byte(h.Config.JWT_SECRET))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.AuthErrorMessage("server_error", "Failed to generate access token"))
		}

		c.Response().Header().Set("Cache-Control", "no-store")
		c.Response().Header().Set("Pragma", "no-cache")

		return c.JSON(http.StatusOK, echo.Map{
			"access_token": token,
			"token_type":   "Bearer",
			"expires_in":   3600, // typically 3600 seconds (1 hour)
			"scope":        value["scope"],
			"client_id":    value["client_id"],
		})
	default:
		return c.JSON(http.StatusBadRequest, helper.AuthErrorMessage("invalid_grant", "Device code is not in a valid state for token issuance"))
	}
}

func (h *instanceHandler) VerifyUserCode(c echo.Context) error {
	userCodeStruct := new(UserCode)
	if err := c.Bind(userCodeStruct); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request body")
	}
	if err := h.Validator.Struct(userCodeStruct); err != nil {
		return c.JSON(400, helper.ErrorMessage("Invalid request body", err))
	}
	userCode := userCodeStruct.Code
	value, err := h.RDB.HGetAll(c.Request().Context(), userCode).Result()
	if err != nil {
		return c.JSON(400, helper.ErrorMessage("Invalid user_code", "User code does not exist or has expired"))
	}
	if value["status"] != "pending" {
		return c.JSON(400, helper.ErrorMessage("Invalid user_code", "User code has already been used or is not pending"))
	}
	scopes := strings.Split(value["scope"], ",")
	firstScope := scopes[0]
	if strings.Split(firstScope, ":")[0] != "user" {
		return c.JSON(400, helper.ErrorMessage("invalid_scope", "User code does not have the correct scope"))
	}
	hostUsername := strings.Split(firstScope, ":")[1]
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*authentication.JwtCustomClaims)
	userId := claims.Id
	instanceId := value["instance_id"]

	// check if instance exists and is active
	sql, _, err := goqu.From("instances").
		Where(goqu.C("id").Eq(instanceId), goqu.C("status").Eq("active")).Select(goqu.COUNT("*")).ToSQL()
	row := h.DB.QueryRow(c.Request().Context(), sql)
	var instanceExists int
	err = row.Scan(&instanceExists)
	if err != nil {
		_, redisErr := h.RDB.HSet(c.Request().Context(), userCode, "status", "denied").Result()
		return c.JSON(400, helper.ErrorMessage("Failed to verify user code", errors.Join(err, redisErr)))
	}
	if instanceExists == 0 {
		_, err = h.RDB.HSet(c.Request().Context(), userCode, "status", "denied").Result()
		if err != nil {
			return c.JSON(400, helper.ErrorMessage("Failed to verify user code", err))
		}
		return c.JSON(400, helper.ErrorMessage("User code verification failed", "Instance does not exist or is not active"))
	}

	sql, _, err = goqu.From("instance_users").
		Where(
			goqu.C("instance_id").Eq(instanceId),
			goqu.C("user_id").Eq(userId),
			goqu.Or(goqu.C("instance_host_username").Eq(hostUsername), goqu.C("instance_host_username").Eq("*")),
		).Select(goqu.COUNT("*")).ToSQL()

	row = h.DB.QueryRow(c.Request().Context(), sql)
	var count int
	err = row.Scan(&count)
	if err != nil {
		_, redisErr := h.RDB.HSet(c.Request().Context(), userCode, "status", "denied").Result()
		return c.JSON(400, helper.ErrorMessage("Failed to verify user code", errors.Join(err, redisErr)))
	}
	if count == 0 {
		sql, _, err = goqu.From("instance_roles").
			Join(goqu.T("role_users"), goqu.On(
				goqu.I("role_users.role_id").Eq(goqu.I("instance_roles.role_id")),
				goqu.I("role_users.user_id").Eq(userId),
			)).
			Where(
				goqu.I("instance_roles.instance_id").Eq(instanceId),
				goqu.Or(
					goqu.I("instance_roles.instance_host_username").Eq(hostUsername),
					goqu.I("instance_roles.instance_host_username").Eq("*"),
				),
			).
			Select(goqu.COUNT("*")).ToSQL()
		row = h.DB.QueryRow(c.Request().Context(), sql)
		var roleCount int
		err = row.Scan(&roleCount)
		if err != nil {
			_, redisErr := h.RDB.HSet(c.Request().Context(), userCode, "status", "denied").Result()
			return c.JSON(400, helper.ErrorMessage("Failed to verify user code", errors.Join(err, redisErr)))
		}
		if roleCount == 0 {
			_, err = h.RDB.HSet(c.Request().Context(), userCode, "status", "denied").Result()
			if err != nil {
				return c.JSON(400, helper.ErrorMessage("Failed to verify user code", err))
			}
			return c.JSON(400, helper.ErrorMessage("User code verification failed", "You do not have permission to access this instance"))
		}
	}
	_, err = h.RDB.HSet(c.Request().Context(), userCode, "status", "approved").Result()
	if err != nil {
		return c.JSON(400, helper.ErrorMessage("Verified user code but failed to update status", err))
	}
	_, err = h.RDB.HSet(c.Request().Context(), value["device_code"], "status", "approved").Result()
	if err != nil {
		return c.JSON(400, helper.ErrorMessage("Verified user code but failed to update device code status", err))
	}
	return c.JSON(200, echo.Map{
		"message": "User code verified successfully",
		"status":  "success",
	})
}
