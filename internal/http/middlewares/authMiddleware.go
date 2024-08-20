package middlewares

import (
	"api/internal/errors"
	jwtlib "api/internal/lib/jwt"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errors.ErrAuthHeaderMissing})
			return
		}

		parts := strings.Split(header, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errors.ErrAuthHeaderInvalid})
			return
		}

		fmt.Println(parts[1])

		token, err := jwtlib.ParseToken(parts[1], secretKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errors.ErrTokenSignatureInvalid})
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errors.ErrTokenInvalid})
			return
		}

		c.Next()
	}
}
