package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"user_advt/internal/domain/users"
	mock_service "user_advt/internal/handler/mocks"
	mock_logger "user_advt/internal/lib/logger/handlers"
	"user_advt/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/testify/v2/require"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetUserById(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type response struct {
		Success bool   `json:"success"`
		Name    string `json:"name,omitempty"`
		Email   string `json:"email,omitempty"`
		Message string `json:"message,omitempty"`
		Status  int    `json:"status,omitempty"`
	}

	testCases := []struct {
		title              string
		path               string
		authHeader         string
		authUserID         string
		authErr            error
		authUser           users.User
		authUserErr        error
		handlerUserID      string
		handlerUser        users.User
		handlerErr         error
		expectedStatus     int
		expectedSuccess    bool
		expectedErrMsg     string
		expectParseToken   bool
		expectAuthUserCall bool
		expectHandlerCall  bool
	}{
		{
			title:              "Success",
			path:               "/api/users/1",
			authHeader:         "Bearer test-token",
			authUserID:         "1",
			authUser:           users.User{ID: 1, Name: "AuthUser", Email: "auth@mail.ru"},
			handlerUserID:      "1",
			handlerUser:        users.User{ID: 1, Name: "TestUser", Email: "test@mail.ru"},
			expectedStatus:     http.StatusOK,
			expectedSuccess:    true,
			expectParseToken:   true,
			expectAuthUserCall: true,
			expectHandlerCall:  true,
		},
		{
			title:              "InvalidIDParam",
			path:               "/api/users/abc",
			authHeader:         "Bearer test-token",
			authUserID:         "1",
			authUser:           users.User{ID: 1, Name: "AuthUser", Email: "auth@mail.ru"},
			expectedStatus:     http.StatusBadRequest,
			expectedSuccess:    false,
			expectedErrMsg:     storage.ErrUserIDInvalidParams.Error(),
			expectParseToken:   true,
			expectAuthUserCall: true,
		},
		{
			title:              "UserNotFound",
			path:               "/api/users/99",
			authHeader:         "Bearer test-token",
			authUserID:         "1",
			authUser:           users.User{ID: 1, Name: "AuthUser", Email: "auth@mail.ru"},
			handlerUserID:      "99",
			handlerErr:         storage.ErrUserNotFount,
			expectedStatus:     http.StatusNotFound,
			expectedSuccess:    false,
			expectedErrMsg:     storage.ErrUserNotFount.Error(),
			expectParseToken:   true,
			expectAuthUserCall: true,
			expectHandlerCall:  true,
		},
		{
			title:           "Unauthorized",
			path:            "/api/users/1",
			authHeader:      "",
			expectedStatus:  http.StatusUnauthorized,
			expectedSuccess: false,
			expectedErrMsg:  storage.ErrUnauthorized.Error(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			mockUserService := mock_service.NewUserService(t)
			mockTokenService := mock_service.NewTokenService(t)

			if tc.expectParseToken {
				mockTokenService.
					On("ParseToken", "test-token").
					Return(tc.authUserID, tc.authErr).
					Once()
			}

			if tc.expectAuthUserCall {
				mockUserService.
					On("GetUser", mock.Anything, tc.authUserID).
					Return(tc.authUser, tc.authUserErr).
					Once()
			}

			if tc.expectHandlerCall {
				mockUserService.
					On("GetUser", mock.Anything, tc.handlerUserID).
					Return(tc.handlerUser, tc.handlerErr).
					Once()
			}

			logger := mock_logger.NewDiscardLogger()
			router := gin.New()
			h := NewHandler(logger, mockUserService, mockTokenService)
			h.Register(router)

			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			var resp response
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
			assert.Equal(t, tc.expectedSuccess, resp.Success)
			if tc.expectedSuccess {
				assert.Equal(t, tc.handlerUser.Name, resp.Name)
				assert.Equal(t, tc.handlerUser.Email, resp.Email)
			} else if tc.expectedErrMsg != "" {
				assert.Equal(t, tc.expectedErrMsg, resp.Message)
			}

			if !tc.expectParseToken {
				mockTokenService.AssertNotCalled(t, "ParseToken", mock.Anything)
			}
			if !tc.expectAuthUserCall && !tc.expectHandlerCall {
				mockUserService.AssertNotCalled(t, "GetUser", mock.Anything, mock.Anything)
			}
		})
	}
}
