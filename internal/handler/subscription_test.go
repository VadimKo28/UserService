package handler

import (
	"app/internal/domain/subscription"
	"app/internal/domain/users"
	mock_service "app/internal/handler/mocks"
	mock_logger "app/internal/lib/logger/handlers"
	"app/internal/storage"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/testify/v2/require"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockSubscriptionService struct {
	mock.Mock
}

func (m *mockSubscriptionService) CreateSubscription(ctx context.Context, dto subscription.CreateSubscriptionDTO) (int, error) {
	args := m.Called(ctx, dto)
	return args.Int(0), args.Error(1)
}

func (m *mockSubscriptionService) GetSubscriptionsByUserID(ctx context.Context, userID, limit, offset int) ([]subscription.Subscription, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]subscription.Subscription), args.Error(1)
}

func (m *mockSubscriptionService) UpdateSubscription(ctx context.Context, dto subscription.UpdateSubscriptionDTO) error {
	args := m.Called(ctx, dto)
	return args.Error(0)
}

func setupSubscriptionRouter(t *testing.T, userSvc *mock_service.UserService, tokenSvc *mock_service.TokenService, subSvc *mockSubscriptionService) *gin.Engine {
	t.Helper()
	logger := mock_logger.NewDiscardLogger()
	router := gin.New()
	h := NewHandler(logger, userSvc, tokenSvc, subSvc)
	h.Register(router)
	return router
}

func expectAuth(t *testing.T, userSvc *mock_service.UserService, tokenSvc *mock_service.TokenService, userID string) {
	t.Helper()
	tokenSvc.
		On("ParseToken", "test-token").
		Return(userID, nil).
		Once()

	userSvc.
		On("GetUser", mock.Anything, userID).
		Return(users.User{ID: 1, Name: "AuthUser", Email: "auth@mail.ru"}, nil).
		Once()
}

func TestCreateUserSubscription(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type response struct {
		Success        bool   `json:"success"`
		SubscriptionID int    `json:"subscription_id,omitempty"`
		Message        string `json:"message,omitempty"`
		Status         int    `json:"status,omitempty"`
	}

	testCases := []struct {
		title             string
		path              string
		body              map[string]any
		mockID            int
		mockErr           error
		expectedStatus    int
		expectedSuccess   bool
		expectedErrMsg    string
		expectServiceCall bool
	}{
		{
			title: "Success",
			path:  "/api/users/1/subscriptions",
			body: map[string]any{
				"service_name": "Netflix",
				"price":        499,
				"start_date":   "2024-05",
				"end_date":     "2024-08",
			},
			mockID:            10,
			expectedStatus:    http.StatusCreated,
			expectedSuccess:   true,
			expectServiceCall: true,
		},
		{
			title: "InvalidUserIDParam",
			path:  "/api/users/abc/subscriptions",
			body: map[string]any{
				"service_name": "Netflix",
				"price":        499,
				"start_date":   "2024-05",
				"end_date":     "2024-08",
			},
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
			expectedErrMsg:  storage.ErrUserIDInvalidParams.Error(),
		},
		{
			title:           "InvalidBody",
			path:            "/api/users/1/subscriptions",
			body:            map[string]any{},
			expectedStatus:  http.StatusUnprocessableEntity,
			expectedSuccess: false,
		},
		{
			title: "ServiceError",
			path:  "/api/users/1/subscriptions",
			body: map[string]any{
				"service_name": "Netflix",
				"price":        499,
				"start_date":   "2024-05",
				"end_date":     "2024-08",
			},
			mockErr:           storage.ErrInternalServerError,
			expectedStatus:    http.StatusInternalServerError,
			expectedSuccess:   false,
			expectedErrMsg:    storage.ErrInternalServerError.Error(),
			expectServiceCall: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			mockUserService := mock_service.NewUserService(t)
			mockTokenService := mock_service.NewTokenService(t)
			mockSubscriptionService := new(mockSubscriptionService)

			expectAuth(t, mockUserService, mockTokenService, "1")

			if tc.expectServiceCall {
				mockSubscriptionService.
					On("CreateSubscription", mock.Anything, mock.MatchedBy(func(dto subscription.CreateSubscriptionDTO) bool {
						startDate, _ := time.Parse("2006-01", "2024-05")
						endDate, _ := time.Parse("2006-01", "2024-08")
						return dto.UserID == 1 &&
							dto.ServiceName == "Netflix" &&
							dto.Price == 499 &&
							dto.StartDate.Equal(startDate) &&
							dto.EndDate.Equal(endDate)
					})).
					Return(tc.mockID, tc.mockErr).
					Once()
			}

			router := setupSubscriptionRouter(t, mockUserService, mockTokenService, mockSubscriptionService)

			body, err := json.Marshal(tc.body)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, tc.path, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer test-token")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			var resp response
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
			assert.Equal(t, tc.expectedSuccess, resp.Success)
			if tc.expectedSuccess {
				assert.Equal(t, tc.mockID, resp.SubscriptionID)
			} else if tc.expectedErrMsg != "" {
				assert.Equal(t, tc.expectedErrMsg, resp.Message)
			}
		})
	}
}

func TestGetUserSubscriptions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type response struct {
		Success       bool                        `json:"success"`
		Limit         int                         `json:"limit,omitempty"`
		Offset        int                         `json:"offset,omitempty"`
		Subscriptions []subscription.Subscription `json:"subscriptions,omitempty"`
		Message       string                      `json:"message,omitempty"`
		Status        int                         `json:"status,omitempty"`
	}

	testCases := []struct {
		title             string
		path              string
		mockItems         []subscription.Subscription
		mockErr           error
		expectedStatus    int
		expectedSuccess   bool
		expectedErrMsg    string
		expectServiceCall bool
	}{
		{
			title: "Success",
			path:  "/api/users/1/subscriptions?limit=2&offset=1",
			mockItems: []subscription.Subscription{
				{ID: 1, UserID: 1, ServiceName: "Netflix", Price: 499, StartDate: time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)},
				{ID: 2, UserID: 1, ServiceName: "Spotify", Price: 299, StartDate: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)},
			},
			expectedStatus:    http.StatusOK,
			expectedSuccess:   true,
			expectServiceCall: true,
		},
		{
			title:           "InvalidLimit",
			path:            "/api/users/1/subscriptions?limit=0",
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
			expectedErrMsg:  storage.ErrInvalidPaginationParams.Error(),
		},
		{
			title:           "InvalidOffset",
			path:            "/api/users/1/subscriptions?offset=-1",
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
			expectedErrMsg:  storage.ErrInvalidPaginationParams.Error(),
		},
		{
			title:             "ServiceError",
			path:              "/api/users/1/subscriptions",
			mockErr:           storage.ErrInternalServerError,
			expectedStatus:    http.StatusInternalServerError,
			expectedSuccess:   false,
			expectedErrMsg:    storage.ErrInternalServerError.Error(),
			expectServiceCall: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			mockUserService := mock_service.NewUserService(t)
			mockTokenService := mock_service.NewTokenService(t)
			mockSubscriptionService := new(mockSubscriptionService)

			expectAuth(t, mockUserService, mockTokenService, "1")

			if tc.expectServiceCall {
				mockSubscriptionService.
					On("GetSubscriptionsByUserID", mock.Anything, 1, mock.AnythingOfType("int"), mock.AnythingOfType("int")).
					Return(tc.mockItems, tc.mockErr).
					Once()
			}

			router := setupSubscriptionRouter(t, mockUserService, mockTokenService, mockSubscriptionService)

			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			req.Header.Set("Authorization", "Bearer test-token")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			var resp response
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
			assert.Equal(t, tc.expectedSuccess, resp.Success)
			if tc.expectedSuccess {
				assert.Len(t, resp.Subscriptions, len(tc.mockItems))
			} else if tc.expectedErrMsg != "" {
				assert.Equal(t, tc.expectedErrMsg, resp.Message)
			}
		})
	}
}

func TestUpdateUserSubscription(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type response struct {
		Success bool   `json:"success"`
		Message string `json:"message,omitempty"`
		Status  int    `json:"status,omitempty"`
	}

	testCases := []struct {
		title             string
		path              string
		body              map[string]any
		mockErr           error
		expectedStatus    int
		expectedSuccess   bool
		expectedErrMsg    string
		expectServiceCall bool
	}{
		{
			title: "Success",
			path:  "/api/users/1/subscriptions/5",
			body: map[string]any{
				"service_name": "Netflix",
				"price":        499,
				"start_date":   "2024-05",
				"end_date":     "2024-08",
			},
			expectedStatus:    http.StatusOK,
			expectedSuccess:   true,
			expectServiceCall: true,
		},
		{
			title: "InvalidSubscriptionID",
			path:  "/api/users/1/subscriptions/abc",
			body: map[string]any{
				"service_name": "Netflix",
				"price":        499,
				"start_date":   "2024-05",
				"end_date":     "2024-08",
			},
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
			expectedErrMsg:  storage.ErrSubscriptionIDInvalidParams.Error(),
		},
		{
			title: "NotFound",
			path:  "/api/users/1/subscriptions/5",
			body: map[string]any{
				"service_name": "Netflix",
				"price":        499,
				"start_date":   "2024-05",
				"end_date":     "2024-08",
			},
			mockErr:           storage.ErrSubscriptionNotFound,
			expectedStatus:    http.StatusNotFound,
			expectedSuccess:   false,
			expectedErrMsg:    storage.ErrSubscriptionNotFound.Error(),
			expectServiceCall: true,
		},
		{
			title:           "InvalidBody",
			path:            "/api/users/1/subscriptions/5",
			body:            map[string]any{},
			expectedStatus:  http.StatusUnprocessableEntity,
			expectedSuccess: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			mockUserService := mock_service.NewUserService(t)
			mockTokenService := mock_service.NewTokenService(t)
			mockSubscriptionService := new(mockSubscriptionService)

			expectAuth(t, mockUserService, mockTokenService, "1")

			if tc.expectServiceCall {
				mockSubscriptionService.
					On("UpdateSubscription", mock.Anything, mock.MatchedBy(func(dto subscription.UpdateSubscriptionDTO) bool {
						startDate, _ := time.Parse("2006-01", "2024-05")
						endDate, _ := time.Parse("2006-01", "2024-08")
						return dto.ID == 5 &&
							dto.UserID == 1 &&
							dto.ServiceName == "Netflix" &&
							dto.Price == 499 &&
							dto.StartDate.Equal(startDate) &&
							dto.EndDate.Equal(endDate)
					})).
					Return(tc.mockErr).
					Once()
			}

			router := setupSubscriptionRouter(t, mockUserService, mockTokenService, mockSubscriptionService)

			body, err := json.Marshal(tc.body)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPut, tc.path, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer test-token")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			var resp response
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
			assert.Equal(t, tc.expectedSuccess, resp.Success)
			if !tc.expectedSuccess && tc.expectedErrMsg != "" {
				assert.Equal(t, tc.expectedErrMsg, resp.Message)
			}
		})
	}
}
