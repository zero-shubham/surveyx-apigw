package wrapper

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/zero-shubham/surveyx-apigw/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Logger interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Warning(v ...interface{})
	Error(v ...interface{})
	Critical(v ...interface{})
	Fatal(v ...interface{})
}

//go:generate mockgen -package mocks -source=wrapper.go -destination=../mocks/wrapper.go *
type AuthServiceClient interface {
	client.AuthServiceClient
}

type wrapperClient struct {
	grpcClient client.AuthServiceClient
	logger     Logger
}

func NewWrapperClient(grpcClient client.AuthServiceClient, logger Logger) *wrapperClient {
	w := wrapperClient{
		grpcClient: grpcClient,
		logger:     logger,
	}

	return &w
}

func (wc *wrapperClient) HandleUserToken(respWtr http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	email := req.FormValue("email")
	password := req.FormValue("password")

	for k, vals := range req.Header {
		for _, v := range vals {
			ctx = metadata.AppendToOutgoingContext(ctx, k, v)
		}
	}
	var respHeader metadata.MD
	resp, err := wc.grpcClient.UserToken(ctx, &client.UserTokenRequest{
		Email:    email,
		Password: password,
	}, grpc.Header(&respHeader))
	if err != nil {
		wc.logger.Error("error while making grpc call: ", err)
		respWtr.WriteHeader(http.StatusInternalServerError)
		return
	}

	wc.logger.Info("call to UserToken successful")

	// Copy headers, status codes, and body from the backend to the response writer
	for k, hs := range respHeader {
		for _, h := range hs {
			respWtr.Header().Add(k, h)
		}
	}
	respWtr.Header().Add("Content-Type", "application/json")

	respWtr.WriteHeader(http.StatusOK)
	if resp == nil {
		wc.logger.Warning("grpc response for UserToken is nil")
		return
	}

	respBody, err := json.Marshal(resp)
	if err != nil {
		wc.logger.Error("error while marshaling resp: %v", err)
		respWtr.WriteHeader(http.StatusInternalServerError)
		return
	}
	wc.logger.Info("writing response body from UserToken")

	_, err = respWtr.Write(respBody)
	if err != nil {
		wc.logger.Error("error while writing resp: %v", err)
		respWtr.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (wc *wrapperClient) HandleCreateUser(respWtr http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	// Read and parse JSON request body
	var requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		OrgID    string `json:"org_id"`
		AppGrpID string `json:"app_grp_id"`
	}

	if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
		wc.logger.Error("error while decoding request body: ", err)
		respWtr.WriteHeader(http.StatusBadRequest)
		return
	}

	// Forward all headers to gRPC context
	for k, vals := range req.Header {
		for _, v := range vals {
			ctx = metadata.AppendToOutgoingContext(ctx, k, v)
		}
	}

	var respHeader metadata.MD
	resp, err := wc.grpcClient.CreateUser(ctx, &client.UserRequest{
		Email:      requestBody.Email,
		Password:   requestBody.Password,
		OrgId:      requestBody.OrgID,
		AppGroupId: requestBody.AppGrpID,
	}, grpc.Header(&respHeader))
	if err != nil {
		wc.logger.Error("error while making grpc call: ", err)
		respWtr.WriteHeader(http.StatusInternalServerError)
		return
	}

	wc.logger.Info("call to CreateUser successful")

	// Copy headers from the backend to the response writer
	for k, hs := range respHeader {
		for _, h := range hs {
			respWtr.Header().Add(k, h)
		}
	}
	respWtr.Header().Add("Content-Type", "application/json")

	respWtr.WriteHeader(http.StatusOK)
	if resp == nil {
		wc.logger.Warning("grpc response for CreateUser is nil")
		return
	}

	respBody, err := json.Marshal(resp)
	if err != nil {
		wc.logger.Error("error while marshaling resp: %v", err)
		respWtr.WriteHeader(http.StatusInternalServerError)
		return
	}
	wc.logger.Info("writing response body from CreateUser")

	_, err = respWtr.Write(respBody)
	if err != nil {
		wc.logger.Error("error while writing resp: %v", err)
		respWtr.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (wc *wrapperClient) HandleCreateApp(respWtr http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	// Read and parse JSON request body
	var requestBody struct {
		OrgID      string `json:"org_id"`
		AppGroupID string `json:"app_group_id"`
	}

	if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
		wc.logger.Error("error while decoding request body: ", err)
		respWtr.WriteHeader(http.StatusBadRequest)
		return
	}

	// Forward all headers to gRPC context
	for k, vals := range req.Header {
		for _, v := range vals {
			ctx = metadata.AppendToOutgoingContext(ctx, k, v)
		}
	}

	var respHeader metadata.MD
	resp, err := wc.grpcClient.CreateApp(ctx, &client.AppRequest{
		AppGroupId: requestBody.AppGroupID,
		OrgId:      requestBody.OrgID,
	}, grpc.Header(&respHeader))
	if err != nil {
		wc.logger.Error("error while making grpc call: ", err)
		respWtr.WriteHeader(http.StatusInternalServerError)
		return
	}

	wc.logger.Info("call to CreateApp successful")

	// Copy headers from the backend to the response writer
	for k, hs := range respHeader {
		for _, h := range hs {
			respWtr.Header().Add(k, h)
		}
	}
	respWtr.Header().Add("Content-Type", "application/json")

	respWtr.WriteHeader(http.StatusOK)
	if resp == nil {
		wc.logger.Warning("grpc response for CreateApp is nil")
		return
	}

	respBody, err := json.Marshal(resp)
	if err != nil {
		wc.logger.Error("error while marshaling resp: %v", err)
		respWtr.WriteHeader(http.StatusInternalServerError)
		return
	}
	wc.logger.Info("writing response body from CreateApp")

	_, err = respWtr.Write(respBody)
	if err != nil {
		wc.logger.Error("error while writing resp: %v", err)
		respWtr.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (wc *wrapperClient) HandleCreateAppGroup(respWtr http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	// Read and parse JSON request body
	var requestBody struct {
		Name   string   `json:"name"`
		Scopes []string `json:"scopes"`
		OrgID  string   `json:"org_id"`
	}

	if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
		wc.logger.Error("error while decoding request body: ", err)
		respWtr.WriteHeader(http.StatusBadRequest)
		return
	}

	// Forward all headers to gRPC context
	for k, vals := range req.Header {
		for _, v := range vals {
			ctx = metadata.AppendToOutgoingContext(ctx, k, v)
		}
	}

	var respHeader metadata.MD
	resp, err := wc.grpcClient.CreateAppGroup(ctx, &client.AppGroupRequest{
		Name:   requestBody.Name,
		Scopes: requestBody.Scopes,
		OrgId:  requestBody.OrgID,
	}, grpc.Header(&respHeader))
	if err != nil {
		wc.logger.Error("error while making grpc call: ", err)
		respWtr.WriteHeader(http.StatusInternalServerError)
		return
	}

	wc.logger.Info("call to CreateAppGroup successful")

	// Copy headers from the backend to the response writer
	for k, hs := range respHeader {
		for _, h := range hs {
			respWtr.Header().Add(k, h)
		}
	}
	respWtr.Header().Add("Content-Type", "application/json")

	respWtr.WriteHeader(http.StatusOK)
	if resp == nil {
		wc.logger.Warning("grpc response for CreateAppGroup is nil")
		return
	}

	respBody, err := json.Marshal(resp)
	if err != nil {
		wc.logger.Error("error while marshaling resp: %v", err)
		respWtr.WriteHeader(http.StatusInternalServerError)
		return
	}
	wc.logger.Info("writing response body from CreateAppGroup")

	_, err = respWtr.Write(respBody)
	if err != nil {
		wc.logger.Error("error while writing resp: %v", err)
		respWtr.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (wc *wrapperClient) HandleGetAppGroup(respWtr http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	// Get app group ID from path parameter
	appGroupID := strings.TrimPrefix(req.URL.Path, "/app-groups/")
	if appGroupID == "" {
		wc.logger.Error("app_group_id is required in path")
		respWtr.WriteHeader(http.StatusBadRequest)
		return
	}
	appGroupID = strings.TrimSuffix(appGroupID, "/")

	// Forward all headers to gRPC context
	for k, vals := range req.Header {
		for _, v := range vals {
			ctx = metadata.AppendToOutgoingContext(ctx, k, v)
		}
	}

	var respHeader metadata.MD
	resp, err := wc.grpcClient.GetAppGroup(ctx, &client.GetAppGroupRequest{
		Id: appGroupID,
	}, grpc.Header(&respHeader))
	if err != nil {
		wc.logger.Error("error while making grpc call: ", err)
		respWtr.WriteHeader(http.StatusInternalServerError)
		return
	}

	wc.logger.Info("call to GetAppGroup successful")

	// Copy headers from the backend to the response writer
	for k, hs := range respHeader {
		for _, h := range hs {
			respWtr.Header().Add(k, h)
		}
	}
	respWtr.Header().Add("Content-Type", "application/json")

	respWtr.WriteHeader(http.StatusOK)
	if resp == nil {
		wc.logger.Warning("grpc response for GetAppGroup is nil")
		return
	}

	respBody, err := json.Marshal(resp)
	if err != nil {
		wc.logger.Error("error while marshaling resp: %v", err)
		respWtr.WriteHeader(http.StatusInternalServerError)
		return
	}
	wc.logger.Info("writing response body from GetAppGroup")

	_, err = respWtr.Write(respBody)
	if err != nil {
		wc.logger.Error("error while writing resp: %v", err)
		respWtr.WriteHeader(http.StatusInternalServerError)
		return
	}
}
