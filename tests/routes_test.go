package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"go-users-api/controllers"
	"go-users-api/models"
	"go-users-api/routes"
)

func TestSetupRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := NewMockUserService()
	controller := controllers.NewUserController(mockService)

	// Configurar rutas
	routes.SetupRoutes(router, controller)

	tests := []struct {
		name           string
		method         string
		path           string
		body           interface{}
		expectedStatus int
	}{
		{
			name:           "Health check",
			method:         "GET",
			path:           "/api/v1/health",
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Create user",
			method: "POST",
			path:   "/api/v1/users",
			body: models.CreateUserRequest{
				Name:    "John Doe",
				Email:   "john.doe@example.com",
				Age:     30,
				Phone:   "+1234567890",
				Address: "123 Main St, City, Country",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Get users",
			method:         "GET",
			path:           "/api/v1/users",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get user by ID",
			method:         "GET",
			path:           "/api/v1/users/test-id",
			expectedStatus: http.StatusInternalServerError, // Mock service returns error
		},
		{
			name:   "Update user",
			method: "PUT",
			path:   "/api/v1/users/test-id",
			body: models.UpdateUserRequest{
				Name:  "John Updated",
				Email: "john.updated@example.com",
				Age:   31,
			},
			expectedStatus: http.StatusInternalServerError, // Mock service returns error
		},
		{
			name:           "Delete user",
			method:         "DELETE",
			path:           "/api/v1/users/test-id",
			expectedStatus: http.StatusInternalServerError, // Mock service returns error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			var err error

			if tt.body != nil {
				jsonBody, _ := json.Marshal(tt.body)
				req, err = http.NewRequest(tt.method, tt.path, bytes.NewBuffer(jsonBody))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req, err = http.NewRequest(tt.method, tt.path, nil)
			}

			assert.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			// Verificar headers de CORS
			if tt.method != "OPTIONS" {
				assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
				assert.Contains(t, w.Header().Get("Access-Control-Allow-Headers"), "Content-Type")
				assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), tt.method)
			}
		})
	}
}

func TestHealthCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Configurar rutas
	mockService := NewMockUserService()
	controller := controllers.NewUserController(mockService)
	routes.SetupRoutes(router, controller)

	req, _ := http.NewRequest("GET", "/api/v1/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response["status"])
	assert.Equal(t, "API Users BRM is running", response["message"])
	assert.NotEmpty(t, response["time"])
}

func TestSwaggerEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Configurar rutas
	mockService := NewMockUserService()
	controller := controllers.NewUserController(mockService)
	routes.SetupRoutes(router, controller)

	req, _ := http.NewRequest("GET", "/swagger/index.html", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Swagger endpoint debería existir (aunque puede devolver 404 si no hay docs generados)
	assert.Contains(t, []int{http.StatusOK, http.StatusNotFound}, w.Code)
}

func TestCORSHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Configurar rutas
	mockService := NewMockUserService()
	controller := controllers.NewUserController(mockService)
	routes.SetupRoutes(router, controller)

	tests := []struct {
		name           string
		origin         string
		expectedOrigin string
	}{
		{
			name:           "Localhost origin",
			origin:         "http://localhost:4200",
			expectedOrigin: "http://localhost:4200",
		},
		{
			name:           "Cloudfront origin",
			origin:         "https://d31rarudcmsl1r.cloudfront.net",
			expectedOrigin: "https://d31rarudcmsl1r.cloudfront.net",
		},
		{
			name:           "Disallowed origin",
			origin:         "http://malicious-site.com",
			expectedOrigin: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/api/v1/health", nil)
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if tt.expectedOrigin != "" {
				assert.Equal(t, tt.expectedOrigin, w.Header().Get("Access-Control-Allow-Origin"))
			} else {
				assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
			}

			assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
			assert.Contains(t, w.Header().Get("Access-Control-Allow-Headers"), "Content-Type")
			assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "GET")
		})
	}
}
