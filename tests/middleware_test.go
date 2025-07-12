package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"go-users-api/middleware"
)

func TestCORS(t *testing.T) {
	router := setupTestRouter()
	router.Use(middleware.CORS())

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	tests := []struct {
		name           string
		origin         string
		expectedOrigin string
		expectedStatus int
	}{
		{
			name:           "Allowed origin - localhost",
			origin:         "http://localhost:4200",
			expectedOrigin: "http://localhost:4200",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Allowed origin - cloudfront",
			origin:         "https://d31rarudcmsl1r.cloudfront.net",
			expectedOrigin: "https://d31rarudcmsl1r.cloudfront.net",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Disallowed origin",
			origin:         "http://malicious-site.com",
			expectedOrigin: "",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "No origin header",
			origin:         "",
			expectedOrigin: "",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/test", nil)
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedOrigin != "" {
				assert.Equal(t, tt.expectedOrigin, w.Header().Get("Access-Control-Allow-Origin"))
			} else {
				assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
			}

			// Verificar otros headers de CORS
			assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
			assert.Contains(t, w.Header().Get("Access-Control-Allow-Headers"), "Content-Type")
			assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "GET")
		})
	}
}

func TestCORS_Options(t *testing.T) {
	router := setupTestRouter()
	router.Use(middleware.CORS())

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	req, _ := http.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "http://localhost:4200")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "http://localhost:4200", w.Header().Get("Access-Control-Allow-Origin"))
}

func TestLogger(t *testing.T) {
	router := setupTestRouter()
	router.Use(middleware.Logger())

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRecovery(t *testing.T) {
	router := setupTestRouter()
	router.Use(middleware.Recovery())

	router.GET("/panic", func(c *gin.Context) {
		panic("test panic")
	})

	router.GET("/normal", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "normal"})
	})

	// Test panic recovery
	t.Run("Panic recovery", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/panic", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Internal Server Error", response["error"])
		assert.Equal(t, "test panic", response["message"])
		assert.Equal(t, float64(500), response["code"])
	})

	// Test normal request
	t.Run("Normal request", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/normal", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "normal", response["message"])
	})
}

func TestRecovery_NonStringPanic(t *testing.T) {
	router := setupTestRouter()
	router.Use(middleware.Recovery())

	router.GET("/panic", func(c *gin.Context) {
		panic(123) // Non-string panic
	})

	req, _ := http.NewRequest("GET", "/panic", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Internal Server Error", response["error"])
	assert.Equal(t, "An unexpected error occurred", response["message"])
	assert.Equal(t, float64(500), response["code"])
}
