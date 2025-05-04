// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/luraproject/lura/v2/config"
	"github.com/zero-shubham/surveyx-apigw/client"
	"github.com/zero-shubham/surveyx-apigw/wrapper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const pluginName = "krakend-grpc-proxy"

// ClientRegisterer is the symbol the plugin loader will try to load. It must implement the RegisterClient interface
var ClientRegisterer = registerer(pluginName)

type registerer string

var logger Logger = nil

func (registerer) RegisterLogger(v interface{}) {
	l, ok := v.(Logger)
	if !ok {
		return
	}
	logger = l
	logger.Debug(fmt.Sprintf("[PLUGIN: %s] Logger loaded", ClientRegisterer))
}

func (r registerer) RegisterClients(f func(
	name string,
	handler func(context.Context, map[string]interface{}) (http.Handler, error),
)) {
	f(string(r), r.registerClients)
}

// go build -buildmode=plugin -o krakend-client-example.so .
func (r registerer) registerClients(ctx context.Context, extra map[string]interface{}) (http.Handler, error) {
	// Access global configuration from context
	cfg, err := config.NewParser().Parse("/etc/krakend/krakend.json")
	if err != nil {
		return nil, errors.New("unable to parse the configuration")
	}
	var host string
	if proxyCfg, ok := cfg.ExtraConfig[pluginName]; ok {
		if proxyCfgParsed, parsed := proxyCfg.(map[string]interface{}); parsed {
			if h, done := proxyCfgParsed["host"].(string); done {
				host = h
			}
		}
	}

	logger.Info("host: ", host)
	// Set up a connection to the server.
	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	grpcClient := client.NewAuthServiceClient(conn)

	client := wrapper.NewWrapperClient(grpcClient, logger)
	wrapper := wrapper.NewGRPCwrapper(logger, wrapper.WrapperParam{
		Endpoint: "/v1/users/token",
		Handler:  client.HandleUserToken,
		Method:   "POST",
	}, wrapper.WrapperParam{
		Endpoint: "/v1/users",
		Handler:  client.HandleCreateUser,
		Method:   "POST",
	}, wrapper.WrapperParam{
		Endpoint: "/v1/apps",
		Handler:  client.HandleCreateApp,
		Method:   "POST",
	}, wrapper.WrapperParam{
		Endpoint: "/v1/app-groups",
		Handler:  client.HandleCreateAppGroup,
		Method:   "POST",
	}, wrapper.WrapperParam{
		Endpoint: "/v1/app-groups",
		Handler:  client.HandleGetAppGroup,
		Method:   "GET",
	})

	// return the actual handler wrapping or your custom logic so it can be used as a replacement for the default http handler
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		// grpc.NewAuthServiceClient().UserToken()
		logger.Debug(fmt.Sprintf("req path: %v %s", req.URL, req.RequestURI))

		wrapper.GetHandler(req.URL.Path, req.Method)(w, req)

	}), nil
}

func main() {}

type Logger interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Warning(v ...interface{})
	Error(v ...interface{})
	Critical(v ...interface{})
	Fatal(v ...interface{})
}

// docker run -it -v "$PWD:/app" -w /app krakend/builder:2.9.3 go build -buildmode=plugin -o krakend-grpc-proxy.so .
