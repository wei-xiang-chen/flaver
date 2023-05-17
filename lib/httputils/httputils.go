package httputils

import (
	"net/http"
	"path"
	"strings"
)

type ApiVerifyType string

const (
	UnknownVerifyType     ApiVerifyType = "unknown"
	WhiteListVerifyType   ApiVerifyType = "whiteList"
	AccessTokenVerifyType ApiVerifyType = "accessToken"
	ScriptVerifyType      ApiVerifyType = "script"
)

func (v ApiVerifyType) ToString() string {
	return string(v)
}

const MethodPathSeparator = " "

var apiVerifyMethodMap = map[string]ApiVerifyType{
	// access token
	strings.Join([]string{http.MethodGet, "/v1/users/*"}, MethodPathSeparator): AccessTokenVerifyType,
	strings.Join([]string{http.MethodPatch, "/v1/users"}, MethodPathSeparator): AccessTokenVerifyType,

	strings.Join([]string{http.MethodGet, "/v1/posts"}, MethodPathSeparator):         AccessTokenVerifyType,
	strings.Join([]string{http.MethodGet, "/v1/posts/*"}, MethodPathSeparator):       AccessTokenVerifyType,
	strings.Join([]string{http.MethodPost, "/v1/posts"}, MethodPathSeparator):        AccessTokenVerifyType,
	strings.Join([]string{http.MethodPatch, "/v1/posts/*"}, MethodPathSeparator):     AccessTokenVerifyType,
	strings.Join([]string{http.MethodPost, "/v1/posts/like/*"}, MethodPathSeparator): AccessTokenVerifyType,

	// whitelist
	strings.Join([]string{http.MethodPost, "/v1/users/register"}, MethodPathSeparator):      WhiteListVerifyType,
	strings.Join([]string{http.MethodPost, "/v1/users/login"}, MethodPathSeparator):         WhiteListVerifyType,
	strings.Join([]string{http.MethodPost, "/v1/users/refresh-token"}, MethodPathSeparator): WhiteListVerifyType,
}

func GetApiVerifyType(request *http.Request) ApiVerifyType {
	if request == nil || request.URL == nil {
		return UnknownVerifyType
	}
	for methodPatchString, verifyType := range apiVerifyMethodMap {
		slice := strings.Split(methodPatchString, MethodPathSeparator)
		method, pattern := slice[0], slice[1]
		if method != request.Method {
			continue
		} else if matched, err := path.Match(
			pattern, request.URL.Path,
		); err != nil {
			continue
		} else if matched {
			return verifyType
		}
	}
	return UnknownVerifyType
}
