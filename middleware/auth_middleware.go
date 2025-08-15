package middleware

import (
	"deni-be-crm/config"
	"deni-be-crm/internal/common"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, common.Error("Token required"))
			c.Abort()
			return
		}
		secret := config.GetEnv("JWT_SECRET", "default_secret")
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(secret), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userIDFloat, ok := claims["user_id"].(float64)
			role, ok := claims["role"].(string)
			isLeader, ok := claims["isLeader"].(bool)
			if !ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
				return
			}
			userID := int(userIDFloat)
			c.Set("user_id", userID)
			c.Set("role", role)
			c.Set("isLeader", isLeader)
			c.Next()
		}

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token not valid"})
			c.Abort()
			return
		}

		c.Next()
	}
}
