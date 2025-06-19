package roles

import "github.com/Utkar5hM/Execstasy/api/utils/config"

type (
	roleHandler struct {
		config.Handler
	}

	CreateRole struct {
		Name        string `json:"name" validate:"required,min=1,max=255,nameregex"`
		Description string `json:"description" validate:"ascii,max=500"`
	}

	ParamsIDStruct struct {
		ID uint64 `param:"id" validate:"required,gt=0"` // Match the path parameter name
	}
	Role struct {
		ID          int    `db:"id"`
		Name        string `db:"name"`
		Description string `db:"description"`
		CreatedBy   string `db:"createdBy"`
		CreatedAt   string `db:"createdAt"`
		UpdatedAt   string `db:"updatedAt"`
	}
	AddUserRole struct {
		Username string `json:"username" validate:"required,ascii,min=2,max=50"`
		Role     string `json:"role" validate:"required,oneof=admin standard"`
	}
)
