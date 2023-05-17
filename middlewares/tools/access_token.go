package tools

import (
	"flaver/api"
	"flaver/lib/constants"
	"net/http"
)

type accessTokenHandler struct{}

func (accessTokenHandler) ProcessRequest(request *http.Request) error {
	if request.Header.Get(constants.HeaderUserUID) == "" ||
		request.Header.Get(constants.HeaderUserRole) == "" {
		return api.PermissionDenied
	}
	return nil
}
