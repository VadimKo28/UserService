package handler

import (
	// "bytes"
	// "encoding/json"
	// "fmt"
	// "net/http"
	// "net/http/httptest"
	"testing"

	// mock_logger "user_advt/internal/lib/logger/handlers"
	// mock_service "user_advt/internal/service/mocks"
	// "user_advt/internal/storage"

	// "github.com/gin-gonic/gin"
	// "github.com/go-openapi/testify/v2/require"
	// "github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	
	test_cases := []struct {
		title   string
		name    string
		email   string
		password string
    respError string 
		mockError error}{
		{
			title:     "Success",
			name:      "TestUser",
			email:     "test@mail.ru",
			password:  "test123",
		},
		// {
		// 	title:     "EmptyName",
		// 	name:      "",
		// 	email:     "test@mail.ru",
		// 	password:  "test123",
		// },
		// {
		// 	title:     "EmptyEmail",
		// 	name:      "TestUser",
		// 	email:     "",
		// 	password:  "test123",
		// },
		// {
		// 	title:     "EmptyPassword",
		// 	name:      "TestUser",
		// 	email:     "test@mail.ru",
		// 	password:  "",
		// },
		// {
		// 	title:     "InvalidEmail",
		// 	name:      "TestUser",
		// 	email:     "invalid-email",
		// 	password:  "test123",
		// },
		// {
		// 	title:     "ShortPassword",
		// 	name:      "TestUser",
		// 	email:     "test@mail.ru",
		// 	password:  "123",
		// },
		// {
		// 	title:     "UserAlreadyExists",
		// 	name:      "TestUser",
		// 	email:     "test@mail.ru",
		// 	password:  "test123",
		// 	respError: "user already exists",
		// 	mockError: storage.ErrUserAlreadyExists,
		// },
	}

	for _, tc := range test_cases {
		t.Run(tc.title, func(t *testing.T) {
			t.Parallel()
			// mockUserService := mock_service.NewUserService(t)

      // handler := NewHandler(nil, mockUserService)

      // handler.service.SignUp(nil, tc.name, tc.email, tc.password)

	
			// require.NoError(t, err)

			// rr := httptest.NewRecorder()

			// router := gin.Default()
			// router.POST("/users", handler.CreateUser)
			// router.ServeHTTP(rr, req)

			// assert.Equal(t, http.StatusOK, rr.Code)

			// body := rr.Body.String()

			// var resp Response

			// require.NoError(t, json.Unmarshal([]byte(body), &resp))

			// require.Equal(t, tc.respError, resp.Error)

			// require.Equal(t, tc.respError == "", resp.Status)

			// if tc.respError == "" {
			// 	require.Equal(t, tc.name, resp.Name)
			// }
			
		})
	}
}
