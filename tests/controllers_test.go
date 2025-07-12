package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"go-users-api/controllers"
	"go-users-api/models"
)

func TestCreateUser(t *testing.T) {
	mockService := NewMockUserService()
	controller := controllers.NewUserController(mockService)
	router := setupTestRouter()

	router.POST("/users", controller.CreateUser)

	tests := []struct {
		name           string
		requestBody    models.CreateUserRequest
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "Valid user creation",
			requestBody: models.CreateUserRequest{
				Name:    "John Doe",
				Email:   "john.doe@example.com",
				Age:     30,
				Phone:   "+1234567890",
				Address: "123 Main St, City, Country",
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name: "Invalid user data - missing name",
			requestBody: models.CreateUserRequest{
				Name:    "",
				Email:   "john.doe@example.com",
				Age:     30,
				Phone:   "+1234567890",
				Address: "123 Main St, City, Country",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if !tt.expectedError {
				var response models.SuccessResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotEmpty(t, response.Message)
				assert.NotNil(t, response.Data)
			}
		})
	}
}

func TestGetUsers(t *testing.T) {
	mockService := NewMockUserService()
	controller := controllers.NewUserController(mockService)
	router := setupTestRouter()

	router.GET("/users", controller.GetUsers)

	// Crear algunos usuarios en el mock
	req1 := models.CreateUserRequest{Name: "John Doe", Email: "john@example.com", Age: 30}
	req2 := models.CreateUserRequest{Name: "Jane Smith", Email: "jane@example.com", Age: 25}
	mockService.CreateUser(context.Background(), req1)
	mockService.CreateUser(context.Background(), req2)

	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.UsersResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), response.Total)
	assert.Len(t, response.Users, 2)
}

func TestGetUserByID(t *testing.T) {
	mockService := NewMockUserService()
	controller := controllers.NewUserController(mockService)
	router := setupTestRouter()

	router.GET("/users/:id", controller.GetUserByID)

	// Crear un usuario en el mock
	req := models.CreateUserRequest{Name: "John Doe", Email: "john@example.com", Age: 30}
	user, _ := mockService.CreateUser(context.Background(), req)

	// Test obtener usuario existente
	req2, _ := http.NewRequest("GET", "/users/"+user.UUID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req2)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Data)

	// Test obtener usuario inexistente
	req3, _ := http.NewRequest("GET", "/users/non-existent-id", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req3)

	assert.Equal(t, http.StatusInternalServerError, w2.Code)
}

func TestUpdateUser(t *testing.T) {
	mockService := NewMockUserService()
	controller := controllers.NewUserController(mockService)
	router := setupTestRouter()

	router.PUT("/users/:id", controller.UpdateUser)

	// Crear un usuario en el mock
	req := models.CreateUserRequest{Name: "John Doe", Email: "john@example.com", Age: 30}
	user, _ := mockService.CreateUser(context.Background(), req)

	// Test actualizar usuario existente
	updateReq := models.UpdateUserRequest{
		Name:    "John Updated",
		Email:   "john.updated@example.com",
		Age:     31,
		Phone:   "+0987654321",
		Address: "456 Oak Ave, Town, State",
	}

	jsonBody, _ := json.Marshal(updateReq)
	req2, _ := http.NewRequest("PUT", "/users/"+user.UUID, bytes.NewBuffer(jsonBody))
	req2.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req2)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Data)

	// Test actualizar usuario inexistente
	req3, _ := http.NewRequest("PUT", "/users/non-existent-id", bytes.NewBuffer(jsonBody))
	req3.Header.Set("Content-Type", "application/json")

	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req3)

	assert.Equal(t, http.StatusInternalServerError, w2.Code)
}

func TestDeleteUser(t *testing.T) {
	mockService := NewMockUserService()
	controller := controllers.NewUserController(mockService)
	router := setupTestRouter()

	router.DELETE("/users/:id", controller.DeleteUser)

	// Crear un usuario en el mock
	req := models.CreateUserRequest{Name: "John Doe", Email: "john@example.com", Age: 30}
	user, _ := mockService.CreateUser(context.Background(), req)

	// Test eliminar usuario existente
	req2, _ := http.NewRequest("DELETE", "/users/"+user.UUID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req2)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "User deleted successfully", response.Message)

	// Test eliminar usuario inexistente
	req3, _ := http.NewRequest("DELETE", "/users/non-existent-id", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req3)

	assert.Equal(t, http.StatusInternalServerError, w2.Code)
}
