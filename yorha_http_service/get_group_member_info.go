package yorha_http_service

import (
	APIStruct "EndlessEmbrace/api_struct"
	BotFunction "EndlessEmbrace/bot_function"
	"EndlessEmbrace/bot_function/yorha_qq_auth"
	yorha_defines "EndlessEmbrace/bot_function/yorha_qq_auth/defines"
	ProcessCenter "EndlessEmbrace/process_center"
	RequestCenter "EndlessEmbrace/request_center"
	"fmt"
	yorha_qq_auth_key "yorha/qq_auth_key"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

// ...
func (r *Router) GetGroupMemberInfo(c *gin.Context) {
	var groupMemberInfo ProcessCenter.GroupMemberInfo

	pk, err := yorha_qq_auth.UnmarshalClientRequest[yorha_defines.GetGroupMemberInfo](
		c.Request.Body, yorha_qq_auth_key.QQAuthKey,
	)
	if err != nil {
		yorha_qq_auth.WriteResponse(c, &yorha_qq_auth_key.QQAuthKey.PublicKey, yorha_defines.ServerResponse{
			ResponseType: yorha_defines.ResponseTypeInvalidRequest,
		})
		return
	}

	resp, err := r.BotClient.Resources.SendRequestWithResponce(
		r.BotClient.Conn,
		RequestCenter.Request{
			Action: APIStruct.GetGroupMemberInfoAction,
			Params: APIStruct.GetGroupMemberInfo{
				GroupId: BotFunction.EulogistGroupID,
				UserID:  pk.QQID,
				NoCache: true,
			},
			RequestId: fmt.Sprintf("%d", r.BotClient.Resources.GetNewRequestId()),
		},
	)
	if err != nil {
		yorha_qq_auth.WriteResponse(c, &yorha_qq_auth_key.QQAuthKey.PublicKey, yorha_defines.ServerResponse{
			ResponseType: yorha_defines.ResponseTypeCQHTTPFailed,
			FailedReason: fmt.Sprintf("GetGroupMemberInfo: %v", err),
		})
		return
	}

	if resp.Status != "OK" {
		if resp.Wording != nil {
			yorha_qq_auth.WriteResponse(c, &yorha_qq_auth_key.QQAuthKey.PublicKey, yorha_defines.ServerResponse{
				ResponseType: yorha_defines.ResponseTypeCQHTTPFailed,
				FailedReason: fmt.Sprintf("GetGroupMemberInfo: Request failed with error %s (%d)", *resp.Wording, resp.Retcode),
			})
			return
		}
		yorha_qq_auth.WriteResponse(c, &yorha_qq_auth_key.QQAuthKey.PublicKey, yorha_defines.ServerResponse{
			ResponseType: yorha_defines.ResponseTypeCQHTTPFailed,
			FailedReason: fmt.Sprintf("GetGroupMemberInfo: Request failed with error code %d", resp.Retcode),
		})
		return
	}

	err = mapstructure.Decode(resp.Data, &groupMemberInfo)
	if err != nil {
		yorha_qq_auth.WriteResponse(c, &yorha_qq_auth_key.QQAuthKey.PublicKey, yorha_defines.ServerResponse{
			ResponseType: yorha_defines.ResponseTypeCQHTTPFailed,
			FailedReason: fmt.Sprintf("GetGroupMemberInfo: %v", err),
		})
		return
	}

	yorha_qq_auth.WriteResponse(c, &yorha_qq_auth_key.QQAuthKey.PublicKey, yorha_defines.ServerResponse{
		ResponseType:    yorha_defines.ResponseTypeGetGroupMemberInfoSuccess,
		GroupMemberInfo: groupMemberInfo,
	})
}
