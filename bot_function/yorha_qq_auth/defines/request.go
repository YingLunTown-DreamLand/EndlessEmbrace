package yorha_defines

const (
	OperationTypeBindEulogistUser   = 1
	OperationTypeUnbindEulogistUser = 2

	OperationTypeBanEulogistUser   = 1
	OperationTypeUnbanEulogistUser = 2

	OperationTypeForceBindEulogistUser   = 1
	OperationTypeForceUnbindEulogistUser = 2

	OperationTypeSetServerCouldAccessWithoutOP   = 1
	OperationTypeUnsetServerCouldAccessWithoutOP = 2

	OperationTypeSearchUserByName = 1
	OperationTypeSearchUserByQQID = 2
)

// 管理类请求的通用字段
type AdminGeneralFields struct {
	AdminQQID     int64  `json:"admin_qq_id"`
	AdminUserName string `json:"admin_user_name"`
}

// Verify Server -> CQ HTTP Server
type GetGroupMemberInfo struct {
	QQID int64 `json:"qq_id"`
}

// Verify Server -> CQ HTTP Server
type NotifyToAllMember struct {
	Message string `json:"message"`
}

// CQ HTTP Server -> Verify Server
type BindEulogistUser struct {
	OperationType    uint8  `json:"operation_type"`
	QQID             int64  `json:"qq_id"`
	AccountNameLower string `json:"account_name_lower,omitempty"` // Only works when is bind
}

// CQ HTTP Server -> Verify Server
type WhoAmI struct {
	QQID int64 `json:"qq_id"`
}

// CQ HTTP Server -> Verify Server
type BanEulogistUser struct {
	AdminGeneralFields
	OperationType uint8 `json:"operation_type"`
	QQID          int64 `json:"qq_id"`
	Duration      int32 `json:"duration,omitempty"` // Only works when is ban
}

// CQ HTTP Server -> Verify Server
type ForceBindEulogistUser struct {
	AdminGeneralFields
	BindEulogistUser
}

// CQ HTTP Server -> Verify Server
type SetCouldAccessWithoutOP struct {
	AdminGeneralFields
	OperationType      uint8  `json:"operation_type"`
	QQID               int64  `json:"qq_id"`
	SetGeneral         bool   `json:"set_general"`
	RentalServerNumber string `json:"rental_server_number,omitempty"` // Only works when SetGeneral is false
}

// CQ HTTP Server -> Verify Server
type SetPermission struct {
	AdminGeneralFields
	QQID       int64 `json:"qq_id"`
	Permission uint8 `json:"permission_level"`
}

// CQ HTTP Server -> Verify Server
type SearchUser struct {
	AdminGeneralFields
	OperationType    uint8  `json:"operation_type"`
	AccountNameLower string `json:"account_name_lower,omitempty"` // Only works when search by name
	QQID             int64  `json:"qq_id,omitempty"`              // Only works when search by qq id
}
