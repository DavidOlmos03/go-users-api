package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// CORS middleware para manejar Cross-Origin Resource Sharing
func CORS() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Orígenes permitidos
		allowedOrigins := []string{
			"http://localhost:4200",
			"https://d31rarudcmsl1r.cloudfront.net", // En prueba por el momento
		}

		// Obtener el origen de la petición
		origin := c.Request.Header.Get("Origin")

		// Verificar si el origen está permitido
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				c.Header("Access-Control-Allow-Origin", origin)
				break
			}
		}

		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
}

// Logger middleware personalizado para logging de requests
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

// Recovery middleware para manejar pánicos
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.JSON(500, gin.H{
				"error":   "Internal Server Error",
				"message": err,
				"code":    500,
			})
		} else {
			c.JSON(500, gin.H{
				"error":   "Internal Server Error",
				"message": "An unexpected error occurred",
				"code":    500,
			})
		}
		c.Abort()
	})
}
