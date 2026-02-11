package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	mock_service "app/internal/handler/mocks"
	mock_logger "app/internal/lib/logger/handlers"
	"app/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/testify/v2/require"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSignUp(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type response struct {
		Success bool   `json:"success"`
		UserID  int    `json:"user_id,omitempty"`
		Message string `json:"message,omitempty"`
		Status  int    `json:"status,omitempty"`
	}

	testCases := []struct {
		title             string
		name              string
		email             string
		password          string
		mockError         error
		mockUserID        int
		expectedStatus    int
		expectedSuccess   bool
		expectedErrMsg    string
		expectServiceCall bool
	}{
		{
			title:             "Success",
			name:              "TestUser",
			email:             "test@mail.ru",
			password:          "test123",
			mockUserID:        10,
			expectedStatus:    http.StatusOK,
			expectedSuccess:   true,
			expectServiceCall: true,
		},
		{
			title:           "EmptyName",
			name:            "",
			email:           "test@mail.ru",
			password:        "test123",
			expectedStatus:  http.StatusUnprocessableEntity,
			expectedSuccess: false,
		},
		{
			title:           "EmptyEmail",
			name:            "TestUser",
			email:           "",
			password:        "test123",
			expectedStatus:  http.StatusUnprocessableEntity,
			expectedSuccess: false,
		},
		{
			title:           "EmptyPassword",
			name:            "TestUser",
			email:           "test@mail.ru",
			password:        "",
			expectedStatus:  http.StatusUnprocessableEntity,
			expectedSuccess: false,
		},
		{
			title:           "InvalidEmail",
			name:            "TestUser",
			email:           "invalid-email",
			password:        "test123",
			expectedStatus:  http.StatusUnprocessableEntity,
			expectedSuccess: false,
		},
		{
			title:           "ShortPassword",
			name:            "TestUser",
			email:           "test@mail.ru",
			password:        "123",
			expectedStatus:  http.StatusUnprocessableEntity,
			expectedSuccess: false,
		},
		{
			title:             "UserAlreadyExists",
			name:              "TestUser",
			email:             "test@mail.ru",
			password:          "test123",
			mockError:         storage.ErrUserAlreadyExists,
			expectedStatus:    http.StatusBadRequest,
			expectedSuccess:   false,
			expectedErrMsg:    storage.ErrUserAlreadyExists.Error(),
			expectServiceCall: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			mockUserService := mock_service.NewUserService(t)

			if tc.expectServiceCall {
				mockUserService.
					On("SignUpUser", mock.Anything, tc.name, tc.email, tc.password).
					Return(tc.mockUserID, tc.mockError).
					Once()
			}

			logger := mock_logger.NewDiscardLogger()
			router := gin.New()
			h := NewHandler(logger, mockUserService, nil, nil)
			h.Register(router)

			body, err := json.Marshal(map[string]string{
				"name":     tc.name,
				"email":    tc.email,
				"password": tc.password,
			})
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/users/sign-up", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			var resp response
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
			assert.Equal(t, tc.expectedSuccess, resp.Success)
			if tc.expectedSuccess {
				assert.Equal(t, tc.mockUserID, resp.UserID)
			} else if tc.expectedErrMsg != "" {
				assert.Equal(t, tc.expectedErrMsg, resp.Message)
			}

			if !tc.expectServiceCall {
				mockUserService.AssertNotCalled(t, "SignUpUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			}
		})
	}
}

func TestSignIn(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type response struct {
		Success      bool   `json:"success"`
		AccessToken  string `json:"access_token,omitempty"`
		RefreshToken string `json:"refresh_token,omitempty"`
		Message      string `json:"message,omitempty"`
		Status       int    `json:"status,omitempty"`
	}

	testCases := []struct {
		title             string
		email             string
		password          string
		mockError         error
		mockAccessToken   string
		mockRefreshToken  string
		expectedStatus    int
		expectedSuccess   bool
		expectedErrMsg    string
		expectServiceCall bool
	}{
		{
			title:             "Success",
			email:             "test@mail.ru",
			password:          "test123",
			mockAccessToken:   "access-token",
			mockRefreshToken:  "refresh-token",
			expectedStatus:    http.StatusOK,
			expectedSuccess:   true,
			expectServiceCall: true,
		},
		{
			title:           "EmptyEmail",
			email:           "",
			password:        "test123",
			expectedStatus:  http.StatusUnprocessableEntity,
			expectedSuccess: false,
		},
		{
			title:           "EmptyPassword",
			email:           "test@mail.ru",
			password:        "",
			expectedStatus:  http.StatusUnprocessableEntity,
			expectedSuccess: false,
		},
		{
			title:           "InvalidEmail",
			email:           "invalid-email",
			password:        "test123",
			expectedStatus:  http.StatusUnprocessableEntity,
			expectedSuccess: false,
		},
		{
			title:           "ShortPassword",
			email:           "test@mail.ru",
			password:        "123",
			expectedStatus:  http.StatusUnprocessableEntity,
			expectedSuccess: false,
		},
		{
			title:             "InvalidCredentials",
			email:             "test@mail.ru",
			password:          "test123",
			mockError:         storage.ErrUserInvalidCredentials,
			expectedStatus:    http.StatusUnauthorized,
			expectedSuccess:   false,
			expectedErrMsg:    storage.ErrUserInvalidCredentials.Error(),
			expectServiceCall: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			mockUserService := mock_service.NewUserService(t)

			if tc.expectServiceCall {
				mockUserService.
					On("SignInUser", mock.Anything, tc.email, tc.password).
					Return(tc.mockAccessToken, tc.mockRefreshToken, tc.mockError).
					Once()
			}

			logger := mock_logger.NewDiscardLogger()
			router := gin.New()
			h := NewHandler(logger, mockUserService, nil, nil)
			h.Register(router)

			body, err := json.Marshal(map[string]string{
				"email":    tc.email,
				"password": tc.password,
			})
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/users/sign-in", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			var resp response
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
			assert.Equal(t, tc.expectedSuccess, resp.Success)
			if tc.expectedSuccess {
				assert.Equal(t, tc.mockAccessToken, resp.AccessToken)
				assert.Equal(t, tc.mockRefreshToken, resp.RefreshToken)
			} else if tc.expectedErrMsg != "" {
				assert.Equal(t, tc.expectedErrMsg, resp.Message)
			}

			if !tc.expectServiceCall {
				mockUserService.AssertNotCalled(t, "SignInUser", mock.Anything, mock.Anything, mock.Anything)
			}
		})
	}
}

func TestLogOut(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type response struct {
		Success bool   `json:"success"`
		Message string `json:"message,omitempty"`
		Status  int    `json:"status,omitempty"`
	}

	testCases := []struct {
		title             string
		refreshToken      string
		mockError         error
		expectedStatus    int
		expectedSuccess   bool
		expectedErrMsg    string
		expectServiceCall bool
	}{
		{
			title:             "Success",
			refreshToken:      "refresh-token",
			expectedStatus:    http.StatusOK,
			expectedSuccess:   true,
			expectServiceCall: true,
		},
		{
			title:           "EmptyRefreshToken",
			refreshToken:    "",
			expectedStatus:  http.StatusUnprocessableEntity,
			expectedSuccess: false,
		},
		{
			title:             "RefreshTokenNotFound",
			refreshToken:      "refresh-token",
			mockError:         storage.ErrRefreshTokenNotFound,
			expectedStatus:    http.StatusUnauthorized,
			expectedSuccess:   false,
			expectedErrMsg:    storage.ErrRefreshTokenNotFound.Error(),
			expectServiceCall: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			mockUserService := mock_service.NewUserService(t)

			if tc.expectServiceCall {
				mockUserService.
					On("LogOut", mock.Anything, tc.refreshToken).
					Return(tc.mockError).
					Once()
			}

			logger := mock_logger.NewDiscardLogger()
			router := gin.New()
			h := NewHandler(logger, mockUserService, nil, nil)
			h.Register(router)

			body, err := json.Marshal(map[string]string{
				"refresh_token": tc.refreshToken,
			})
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/logout", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			var resp response
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
			assert.Equal(t, tc.expectedSuccess, resp.Success)
			if !tc.expectedSuccess && tc.expectedErrMsg != "" {
				assert.Equal(t, tc.expectedErrMsg, resp.Message)
			}

			if !tc.expectServiceCall {
				mockUserService.AssertNotCalled(t, "LogOut", mock.Anything, mock.Anything)
			}
		})
	}
}

func TestRefresh(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type response struct {
		Success      bool   `json:"success"`
		AccessToken  string `json:"access_token,omitempty"`
		RefreshToken string `json:"refresh_token,omitempty"`
		Message      string `json:"message,omitempty"`
		Status       int    `json:"status,omitempty"`
	}

	testCases := []struct {
		title             string
		refreshToken      string
		mockError         error
		mockAccessToken   string
		mockRefreshToken  string
		expectedStatus    int
		expectedSuccess   bool
		expectedErrMsg    string
		expectServiceCall bool
	}{
		{
			title:             "Success",
			refreshToken:      "refresh-token",
			mockAccessToken:   "access-token",
			mockRefreshToken:  "new-refresh-token",
			expectedStatus:    http.StatusOK,
			expectedSuccess:   true,
			expectServiceCall: true,
		},
		{
			title:           "EmptyRefreshToken",
			refreshToken:    "",
			expectedStatus:  http.StatusUnprocessableEntity,
			expectedSuccess: false,
		},
		{
			title:             "RefreshTokenExpired",
			refreshToken:      "refresh-token",
			mockError:         storage.ErrRefreshTokenExpired,
			expectedStatus:    http.StatusUnauthorized,
			expectedSuccess:   false,
			expectedErrMsg:    storage.ErrRefreshTokenExpired.Error(),
			expectServiceCall: true,
		},
		{
			title:             "RefreshTokenNotFound",
			refreshToken:      "refresh-token",
			mockError:         storage.ErrRefreshTokenNotFound,
			expectedStatus:    http.StatusUnauthorized,
			expectedSuccess:   false,
			expectedErrMsg:    storage.ErrRefreshTokenNotFound.Error(),
			expectServiceCall: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			mockUserService := mock_service.NewUserService(t)

			if tc.expectServiceCall {
				mockUserService.
					On("RefreshTokens", mock.Anything, tc.refreshToken).
					Return(tc.mockAccessToken, tc.mockRefreshToken, tc.mockError).
					Once()
			}

			logger := mock_logger.NewDiscardLogger()
			router := gin.New()
			h := NewHandler(logger, mockUserService, nil, nil)
			h.Register(router)

			body, err := json.Marshal(map[string]string{
				"refresh_token": tc.refreshToken,
			})
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			var resp response
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
			assert.Equal(t, tc.expectedSuccess, resp.Success)
			if tc.expectedSuccess {
				assert.Equal(t, tc.mockAccessToken, resp.AccessToken)
				assert.Equal(t, tc.mockRefreshToken, resp.RefreshToken)
			} else if tc.expectedErrMsg != "" {
				assert.Equal(t, tc.expectedErrMsg, resp.Message)
			}

			if !tc.expectServiceCall {
				mockUserService.AssertNotCalled(t, "RefreshTokens", mock.Anything, mock.Anything)
			}
		})
	}
}
