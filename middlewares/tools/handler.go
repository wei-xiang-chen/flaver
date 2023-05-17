package tools

import (
	"flaver/lib/httputils"
	"net/http"
)

func GetApiVerifyHandler(request *http.Request) ApiVerifyHandler {
	switch httputils.GetApiVerifyType(request) {
	case httputils.WhiteListVerifyType:
		return whiteListHandler{}
	case httputils.AccessTokenVerifyType:
		return accessTokenHandler{}
	// case httputils.ScriptVerifyType:
	// 	return scriptHandler{}
	default:
		return unsupportHandler{}
	}
}
