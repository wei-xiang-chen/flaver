package tools

import "net/http"

type whiteListHandler struct {
}

func (whiteListHandler) ProcessRequest(*http.Request) error {
	return nil
}
