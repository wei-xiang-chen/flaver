package main

import (
	"context"
	"errors"
	"flaver/gateway"
	"flaver/globals"
	"flaver/initialize"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	serverAddr := 50001
	proxyURL := fmt.Sprintf("http://localhost:%d", serverAddr)
	httpServers := []http.Server{
		{
			Addr:    fmt.Sprintf(":%d", globals.GetConfig().GetServer().GetAddr()),
			Handler: gateway.InitServer(proxyURL),
		},
		{
			Addr:    fmt.Sprintf(":%d", serverAddr),
			Handler: initialize.Routers(),
		},
	}

	for _, server := range httpServers {
		go func(server http.Server) {
			if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				globals.GetLogger().Warnf("[ListenAndServe] error: %s", err.Error())
			}
		}(server)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(len(httpServers))
	for _, server := range httpServers {
		go func(server http.Server, wg *sync.WaitGroup) {
			if err := server.Shutdown(ctx); err != nil {
				globals.GetLogger().Warnf("[Shutdown] error: %s", err.Error())
			}
			wg.Done()
		}(server, &wg)
	}
	wg.Wait()
	globals.GetLogger().Info("[Graceful Shut Down] ...")
}