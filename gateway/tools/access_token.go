package tools

import (
	"flaver/lib/constants"
	"flaver/lib/utils"
	"net/http"
	"strings"
)

type accessTokenHandler struct {
}

func (this accessTokenHandler) ProcessRequest(request *http.Request) {
	authorization := request.Header.Get("Authorization")

	if jwtToken, err := utils.ValidateJWTToken(strings.Replace(authorization, "Bearer ", "", -1)); err != nil {
		request.Header.Set(constants.HeaderUserUID, "")
		request.Header.Set(constants.HeaderUserRole, "")
	} else if claims, ok := jwtToken.Claims.(*utils.FlaverCliaims); !ok {
		request.Header.Set(constants.HeaderUserUID, "")
		request.Header.Set(constants.HeaderUserRole, "")
	} else {
		request.Header.Set(constants.HeaderUserUID, utils.GetUIDFromToken(jwtToken))
		request.Header.Set(constants.HeaderUserRole, claims.Role)
	}
}
