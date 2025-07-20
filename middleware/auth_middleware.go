package middleware

import (
    "fmt"
    "os"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "github.com/joho/godotenv"
    "task6/models"
)

func init() {
    godotenv.Load()
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateJWT(user models.User) (string, error) {
    claims := jwt.MapClaims{
        "username": user.Username,
        "role":  user.Role,
        "exp":   time.Now().Add(2 * time.Hour).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}


func AdminAuth() gin.HandlerFunc {
	return func (c *gin.Context){
		role,exists := c.Get("role")
		if !exists || role != "admin" {
			c.JSON(403,gin.H{"error":"Admin access only"})
			c.Abort()
			return
		}
		c.Next()
	}
}



func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(401, gin.H{"error": "Authorization header missing"})
            c.Abort()
            return
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
            c.JSON(401, gin.H{"error": "Invalid Authorization header format"})
            c.Abort()
            return
        }

        tokenStr := parts[1]
        token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            return jwtSecret, nil
        })

        if err != nil || token == nil || !token.Valid {
            c.JSON(401, gin.H{"error": "Invalid or expired token"})
            c.Abort()
            return
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            c.JSON(401, gin.H{"error": "Failed to parse claims"})
            c.Abort()
            return
        }

        c.Set("username", claims["username"])
        c.Set("role", claims["role"])

        c.Next() 

    }
}
