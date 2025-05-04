package wrapper

import (
	"net/http"
	"slices"
	"strings"
)

var allwedMethods = []string{"POST", "GET", "PATCH", "PUT", "DELETE"}

type wrapper struct {
	endpointMap map[string]map[string]http.HandlerFunc
	logger      Logger
}

type WrapperParam struct {
	Endpoint string
	Method   string
	Handler  http.HandlerFunc
}

func NewGRPCwrapper(logger Logger, opts ...WrapperParam) *wrapper {
	w := wrapper{
		endpointMap: make(map[string]map[string]http.HandlerFunc, len(opts)),
		logger:      logger,
	}

	for _, opt := range opts {
		if !slices.Contains(allwedMethods, strings.ToUpper(opt.Method)) {
			logger.Warning("unexpected method passed, ignoring opt: %v", opt)
		}
		if _, ok := w.endpointMap[opt.Endpoint]; !ok {
			w.endpointMap[opt.Endpoint] = make(map[string]http.HandlerFunc, 1)
		}
		w.endpointMap[opt.Endpoint][opt.Method] = opt.Handler
	}
	return &w
}

func (w *wrapper) GetHandler(urlPath, method string) http.HandlerFunc {
	handlers, ok := w.endpointMap[strings.TrimSuffix(urlPath, "/")]
	if !ok {
		w.logger.Error("not found handler: ", urlPath)
		return http.NotFound
	}
	w.logger.Info("found handler: ", urlPath)
	methodHandler, ok := handlers[strings.ToUpper(method)]
	if !ok {
		w.logger.Error("not found handler for method: ", urlPath, method)
		return http.NotFound
	}
	return methodHandler
}
