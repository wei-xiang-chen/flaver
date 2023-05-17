package tools

import "net/http"

type GatewayHandler interface {
	ProcessRequest(*http.Request)
}

type doNothing struct {
}

func (doNothing) ProcessRequest(*http.Request) {
}