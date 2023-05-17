package utils

import (
	"flaver/lib/constants"
	"net/http"
	"strings"
)

func GetHttpClientIP(request *http.Request) string {
	ip := strings.TrimSpace(strings.Split(request.Header.Get(constants.HeaderClientIP), ",")[0])
	if ip == "" {
		ip = request.RemoteAddr
	}
	return ip
}
