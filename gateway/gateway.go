package gateway

import (
	"flaver/gateway/tools"
	"flaver/lib/constants"
	"flaver/lib/utils"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// Reverse Proxy Server.
// Proxy To "toURL".
// Only Handle Add Header.
func InitServer(toURL string) *httputil.ReverseProxy {
	url, err := url.Parse(toURL)
	if err != nil {
		panic(err)
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(url)
	reverseProxy.Director = func(request *http.Request) {
		request.URL.Scheme = url.Scheme
		request.URL.Host = url.Host
		request.Host = url.Host
		handler := tools.GetGatewayHandler(request)
		handler.ProcessRequest(request)
		request.Header.Set(constants.HeaderClientIP, utils.GetHttpClientIP(request))
	}
	return reverseProxy
}