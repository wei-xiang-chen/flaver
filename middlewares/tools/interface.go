package tools

import (
	"flaver/api"
	"net/http"
)

type ApiVerifyHandler interface {
	ProcessRequest(*http.Request) error
}

type unsupportHandler struct {
}

func (unsupportHandler) ProcessRequest(*http.Request) error {
	return api.Unimplemented
}
