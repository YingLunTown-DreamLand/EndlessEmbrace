package yorha_defines

import ProcessCenter "EndlessEmbrace/process_center"

const (
	ResponseTypeGetGroupMemberInfoSuccess = iota + 1
	ResponseTypeNotifyToAllMemberSuccess

	ResponseTypeBindOperationSuccess           // Bind or unbind
	ResponseTypeWhoAmISuccess                  // Who am I
	ResponseTypeBanOperationSuccess            // Ban or unban
	ResponseTypeForceBindOperationSuccess      // Force bind or force unbind
	ResponseTypeSetCouldAccessWithoutOPSuccess // Set/unset could access without op
	ResponseTypeSetPermissionSuccess           // Only for set permission
	ResponseTypeSearchOperationSuccess         // Search by name or search by qq id

	ResponseTypeInvalidRequest
	ResponseTypeUserNotFound
	ResponseTypeUserHasBinded
	ResponseTypeQQHasBinded
	ResponseTypeUserNotBinded
	ResponseTypeUserHasBanned
	ResponseTypeUserNotBanned
	ResponseTypeServerNotFound
	ResponseTypeUnknownPermissionLevel

	ResponseTypeCQHTTPFailed = 255
)

// Both Verify Server and CQ HTTP Server
// could send this.
type ServerResponse struct {
	ResponseType int `json:"response_type"`

	// GetGroupMemberInfo
	GroupMemberInfo ProcessCenter.GroupMemberInfo `json:"group_member_info,omitempty"`
	FailedReason    string                        `json:"failed_reason"`
	// UnbindEulogistUser, ForceUnbindEulogistUser
	OriginAccountNameLower string `json:"origin_account_name_lower,omitempty"`
	// Ban
	UnbanUnixTime int64 `json:"unban_unix_time,omitempty"`
	// SetPermission
	ResultPermissionName string `json:"result_permission_name,omitempty"`
	// WhoAmI, SearchUser
	UserData *User `json:"user_data,omitempty"`
}
