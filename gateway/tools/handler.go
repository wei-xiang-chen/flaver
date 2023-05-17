package tools

import (
	"flaver/lib/httputils"
	"net/http"
)

func GetGatewayHandler(request *http.Request) GatewayHandler {
	switch httputils.GetApiVerifyType(request) {
	case httputils.AccessTokenVerifyType:
		return accessTokenHandler{
		}
	case httputils.WhiteListVerifyType,
		httputils.ScriptVerifyType:
		return doNothing{}
	default:
		return doNothing{}
	}
}
