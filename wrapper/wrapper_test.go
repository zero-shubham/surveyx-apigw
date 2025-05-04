package wrapper_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zero-shubham/surveyx-apigw/client"
	"github.com/zero-shubham/surveyx-apigw/mocks"
	"github.com/zero-shubham/surveyx-apigw/wrapper"
	"go.uber.org/mock/gomock"
)

func TestWrapper(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedClient := mocks.NewMockAuthServiceClient(ctrl)
	mockedLogger := mocks.NewMockLogger(ctrl)

	mw := wrapper.NewWrapperClient(mockedClient, mockedLogger)

	// Common test data
	const (
		testEmail      = "test@example.com"
		testPassword   = "testpass"
		testOrgID      = "org123"
		testAppGrpID   = "appgrp123"
		testUserID     = "user123"
		testAppID      = "app123"
		testToken      = "test-token"
		testAppGrpName = "Test App Group"
	)

	var testScopes = []string{"read", "write", "admin"}

	t.Run("should be able to retrieve user token", func(t *testing.T) {
		// Create test request
		req := httptest.NewRequest(http.MethodPost, "/user/token", nil)
		req.Form = map[string][]string{
			"email":    {testEmail},
			"password": {testPassword},
		}

		// Create response recorder
		w := httptest.NewRecorder()

		// Setup mock expectations
		expectedResp := &client.TokenResponse{
			AccessToken: testToken,
		}

		mockedClient.EXPECT().
			UserToken(gomock.Any(), &client.UserTokenRequest{
				Email:    testEmail,
				Password: testPassword,
			}, gomock.Any()).
			Return(expectedResp, nil)

		mockedLogger.EXPECT().Info("call to UserToken successful")
		mockedLogger.EXPECT().Info("writing response body from UserToken")

		// Call the handler
		mw.HandleUserToken(w, req)

		// Assert response
		assert.Equal(t, http.StatusOK, w.Code)

		var response client.TokenResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, testToken, response.AccessToken)
	})

	t.Run("error while trying to retrieve user token", func(t *testing.T) {
		// Create test request
		req := httptest.NewRequest(http.MethodPost, "/user/token", nil)
		req.Form = map[string][]string{
			"email":    {testEmail},
			"password": {testPassword},
		}

		// Create response recorder
		w := httptest.NewRecorder()

		// Setup mock expectations
		mockedClient.EXPECT().
			UserToken(gomock.Any(), &client.UserTokenRequest{
				Email:    testEmail,
				Password: testPassword,
			}, gomock.Any()).
			Return(nil, fmt.Errorf("internal error"))

		mockedLogger.EXPECT().Error("error while making grpc call: ", gomock.Any())

		// Call the handler
		mw.HandleUserToken(w, req)

		// Assert response
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should be able to create user", func(t *testing.T) {
		// Create test request body
		requestBody := struct {
			Email    string `json:"email"`
			Password string `json:"password"`
			OrgID    string `json:"org_id"`
			AppGrpID string `json:"app_grp_id"`
		}{
			Email:    testEmail,
			Password: testPassword,
			OrgID:    testOrgID,
			AppGrpID: testAppGrpID,
		}
		body, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		// Create test request
		req := httptest.NewRequest(http.MethodPost, "/user", nil)
		req.Body = io.NopCloser(bytes.NewReader(body))

		// Create response recorder
		w := httptest.NewRecorder()

		// Setup mock expectations
		expectedResp := &client.UserResponse{
			Id:         testUserID,
			Email:      testEmail,
			OrgId:      testOrgID,
			AppGroupId: testAppGrpID,
		}

		mockedClient.EXPECT().
			CreateUser(gomock.Any(), &client.UserRequest{
				Email:      testEmail,
				Password:   testPassword,
				OrgId:      testOrgID,
				AppGroupId: testAppGrpID,
			}, gomock.Any()).
			Return(expectedResp, nil)

		mockedLogger.EXPECT().Info("call to CreateUser successful")
		mockedLogger.EXPECT().Info("writing response body from CreateUser")

		// Call the handler
		mw.HandleCreateUser(w, req)

		// Assert response
		assert.Equal(t, http.StatusOK, w.Code)

		var response client.UserResponse
		err = json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, testUserID, response.Id)
		assert.Equal(t, testEmail, response.Email)
		assert.Equal(t, testOrgID, response.OrgId)
		assert.Equal(t, testAppGrpID, response.AppGroupId)
	})

	t.Run("should be able to create app", func(t *testing.T) {
		// Create test request body
		requestBody := struct {
			OrgID      string `json:"org_id"`
			AppGroupID string `json:"app_group_id"`
		}{
			OrgID:      testOrgID,
			AppGroupID: testAppGrpID,
		}
		body, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		// Create test request
		req := httptest.NewRequest(http.MethodPost, "/app", nil)
		req.Body = io.NopCloser(bytes.NewReader(body))

		// Create response recorder
		w := httptest.NewRecorder()

		// Setup mock expectations
		expectedResp := &client.AppResponse{
			Id:         testAppID,
			OrgId:      testOrgID,
			AppGroupId: testAppGrpID,
		}

		mockedClient.EXPECT().
			CreateApp(gomock.Any(), &client.AppRequest{
				OrgId:      testOrgID,
				AppGroupId: testAppGrpID,
			}, gomock.Any()).
			Return(expectedResp, nil)

		mockedLogger.EXPECT().Info("call to CreateApp successful")
		mockedLogger.EXPECT().Info("writing response body from CreateApp")

		// Call the handler
		mw.HandleCreateApp(w, req)

		// Assert response
		assert.Equal(t, http.StatusOK, w.Code)

		var response client.AppResponse
		err = json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, testAppID, response.Id)
		assert.Equal(t, testOrgID, response.OrgId)
		assert.Equal(t, testAppGrpID, response.AppGroupId)
	})

	t.Run("should be able to create app group", func(t *testing.T) {
		// Create test request body
		requestBody := struct {
			Name   string   `json:"name"`
			Scopes []string `json:"scopes"`
			OrgID  string   `json:"org_id"`
		}{
			Name:   testAppGrpName,
			Scopes: testScopes,
			OrgID:  testOrgID,
		}
		body, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		// Create test request
		req := httptest.NewRequest(http.MethodPost, "/app-group", nil)
		req.Body = io.NopCloser(bytes.NewReader(body))

		// Create response recorder
		w := httptest.NewRecorder()

		// Setup mock expectations
		expectedResp := &client.AppGroupResponse{
			Id:     testAppGrpID,
			Name:   testAppGrpName,
			Scopes: testScopes,
			OrgId:  testOrgID,
		}

		mockedClient.EXPECT().
			CreateAppGroup(gomock.Any(), &client.AppGroupRequest{
				Name:   testAppGrpName,
				Scopes: testScopes,
				OrgId:  testOrgID,
			}, gomock.Any()).
			Return(expectedResp, nil)

		mockedLogger.EXPECT().Info("call to CreateAppGroup successful")
		mockedLogger.EXPECT().Info("writing response body from CreateAppGroup")

		// Call the handler
		mw.HandleCreateAppGroup(w, req)

		// Assert response
		assert.Equal(t, http.StatusOK, w.Code)

		var response client.AppGroupResponse
		err = json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, testAppGrpID, response.Id)
		assert.Equal(t, testAppGrpName, response.Name)
		assert.Equal(t, testScopes, response.Scopes)
		assert.Equal(t, testOrgID, response.OrgId)
	})

	t.Run("should be able to get app group", func(t *testing.T) {
		// Create test request
		req := httptest.NewRequest(http.MethodGet, "/app-groups/"+testAppGrpID, nil)

		// Create response recorder
		w := httptest.NewRecorder()

		// Setup mock expectations
		expectedResp := &client.AppGroupResponse{
			Id:     testAppGrpID,
			Name:   testAppGrpName,
			Scopes: testScopes,
			OrgId:  testOrgID,
		}

		mockedClient.EXPECT().
			GetAppGroup(gomock.Any(), &client.GetAppGroupRequest{
				Id: testAppGrpID,
			}, gomock.Any()).
			Return(expectedResp, nil)

		mockedLogger.EXPECT().Info("call to GetAppGroup successful")
		mockedLogger.EXPECT().Info("writing response body from GetAppGroup")

		// Call the handler
		mw.HandleGetAppGroup(w, req)

		// Assert response
		assert.Equal(t, http.StatusOK, w.Code)

		var response client.AppGroupResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, testAppGrpID, response.Id)
		assert.Equal(t, testAppGrpName, response.Name)
		assert.Equal(t, testScopes, response.Scopes)
		assert.Equal(t, testOrgID, response.OrgId)
	})
}
