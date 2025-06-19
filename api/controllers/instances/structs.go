package instances

import "github.com/Utkar5hM/Execstasy/api/utils/config"

type (
	instanceHandler struct {
		config.Handler
	}
	InstanceUsernameStruct struct {
		HostUsername string `json:"host_username" validate:"required,min=1,max=32,linuxuser"`
		Username     string `json:"username" validate:"required,ascii,min=2,max=50"`
	}

	hostUsernameStruct struct {
		Username string `json:"host_username" validate:"required,min=1,max=32,linuxuser"`
	}

	InstanceRole struct {
		Id           int    `json:"id"`
		Name         string `json:"name"`
		HostUsername string `json:"host_username"`
	}

	IDUsernameStruct struct {
		ID       uint64 `json:"id" validate:"required,gt=0"`
		Username string `json:"host_username" validate:"required,min=1,max=32,linuxuser"`
	}

	Instances struct {
		ID          int    `db:"id"`
		Name        string `db:"name"`
		HostAddress string `db:"host_address"`
		Status      string `db:"status"`
		CreatedBy   string `db:"created_by"`
	}

	Instance struct {
		ID          int      `db:"id"`
		Name        string   `db:"name"`
		HostAddress string   `db:"host_address"`
		Status      string   `db:"status"`
		CreatedBy   string   `db:"created_by"`
		Description string   `db:"description"`
		HostUsers   []string `db:"host_users"`
		ClientID    string   `db:"client_id"`
	}

	ParamsIDStruct struct {
		ID uint64 `param:"id" validate:"required,gt=0"` // Match the path parameter name
	}

	createInstanceStruct struct {
		Name        string `json:"name" validate:"required,min=1,max=255,nameregex"`
		Description string `json:"description" validate:"required,ascii,min=1,max=500"`
		HostAddress string `json:"HostAddress" validate:"required,hostname_rfc1123|ip"`
		Status      string `json:"status" validate:"required,oneof=active inactive"`
	}

	InstanceUser struct {
		Id           int    `json:"id"`
		Username     string `json:"username"`
		Role         string `json:"role"`
		Name         string `json:"name"`
		HostUsername string `json:"host_username"`
	}

	deviceAuthorizationRequest struct {
		ClientId string `form:"client_id" validate:"required,uuid4"`
		Scope    string `form:"scope" validate:"required,scopelinuxuser"`
	}

	tokenRequest struct {
		GrantType  string `form:"grant_type" validate:"required,eq=urn:ietf:params:oauth:grant-type:device_code"`
		ClientId   string `form:"client_id" validate:"required,uuid4"`
		DeviceCode string `form:"device_code" validate:"required,uuid4"`
	}

	UserCode struct {
		Code string `json:"user_code" validate:"required,len=8,alphanumunicode"`
	}
)
