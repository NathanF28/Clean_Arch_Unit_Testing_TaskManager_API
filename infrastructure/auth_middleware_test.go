package infrastructure_test

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"task7/infrastructure"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testSecret = []byte("supersecretkeyforunittests123")

func createTestContext(w *httptest.ResponseRecorder) *gin.Context {
	c, _ := gin.CreateTestContext(w)
	return c
}

func generateTestToken(t *testing.T, username, role string, expiration time.Time) string {
	claims := jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      expiration.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(testSecret)
	require.NoError(t, err, "Failed to sign test token")
	return tokenString
}

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	originalSecret := infrastructure.GetJWTSecret()
	defer func() {
		infrastructure.SetJWTSecret(originalSecret)
	}()

	tests := []struct {
		name         string
		authHeader   string
		expectedCode int
		expectedBody string
		assertNext   bool
		expectedUser string
		expectedRole string
	}{
		{
			name:         "Missing Authorization Header",
			authHeader:   "",
			expectedCode: http.StatusUnauthorized,
			expectedBody: `{"error":"Authorization header missing"}`,
			assertNext:   false,
		},
		{
			name:         "Invalid Authorization Header Format - No Bearer",
			authHeader:   "InvalidToken",
			expectedCode: http.StatusUnauthorized,
			expectedBody: `{"error":"Invalid Authorization header format"}`,
			assertNext:   false,
		},
		{
			name:         "Invalid Authorization Header Format - Missing Token",
			authHeader:   "Bearer",
			expectedCode: http.StatusUnauthorized,
			expectedBody: `{"error":"Invalid Authorization header format"}`,
			assertNext:   false,
		},
		{
			name:         "Malformed JWT",
			authHeader:   "Bearer abc.def.ghi",
			expectedCode: http.StatusUnauthorized,
			expectedBody: `{"error":"Invalid or expired token"}`,
			assertNext:   false,
		},
		{
			name: "Expired JWT",
			authHeader: func() string {
				token := generateTestToken(t, "testuser", "regular", time.Now().Add(-1*time.Hour))
				return "Bearer " + token
			}(),
			expectedCode: http.StatusUnauthorized,
			expectedBody: `{"error":"Invalid or expired token"}`,
			assertNext:   false,
		},
		{
			name: "Valid JWT - User Role",
			authHeader: func() string {
				token := generateTestToken(t, "user1", "regular", time.Now().Add(1*time.Hour))
				return "Bearer " + token
			}(),
			expectedCode: http.StatusOK,
			assertNext:   true,
			expectedUser: "user1",
			expectedRole: "regular",
		},
		{
			name: "Valid JWT - Admin Role",
			authHeader: func() string {
				token := generateTestToken(t, "admin1", "admin", time.Now().Add(1*time.Hour))
				return "Bearer " + token
			}(),
			expectedCode: http.StatusOK,
			assertNext:   true,
			expectedUser: "admin1",
			expectedRole: "admin",
		},
		{
			name: "JWT with Invalid Signature (wrong secret)",
			authHeader: func() string {
				otherSecret := []byte("another_secret_key_different_from_test_secret")
				claims := jwt.MapClaims{
					"username": "attacker",
					"role":     "admin",
					"exp":      time.Now().Add(1 * time.Hour).Unix(),
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString(otherSecret)
				return "Bearer " + tokenString
			}(),
			expectedCode: http.StatusUnauthorized,
			expectedBody: `{"error":"Invalid or expired token"}`,
			assertNext:   false,
		},
		{
			name: "JWT with Missing Username Claim",
			authHeader: func() string {
				claims := jwt.MapClaims{
					"role": "regular",
					"exp":  time.Now().Add(1 * time.Hour).Unix(),
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString(testSecret)
				return "Bearer " + tokenString
			}(),
			expectedCode: http.StatusOK,
			assertNext:   true,
			expectedUser: "",
			expectedRole: "regular",
		},
		{
			name: "JWT with Missing Role Claim",
			authHeader: func() string {
				claims := jwt.MapClaims{
					"username": "userWithoutRole",
					"exp":      time.Now().Add(1 * time.Hour).Unix(),
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString(testSecret)
				return "Bearer " + tokenString
			}(),
			expectedCode: http.StatusOK,
			assertNext:   true,
			expectedUser: "userWithoutRole",
			expectedRole: "",
		},
		{
			name: "JWT with Unexpected Signing Method (e.g., RS256)",
			authHeader: func() string {
				header := `{"alg":"RS256","typ":"JWT"}`
				payload := `{"username":"testuser","role":"regular","exp":` + `1234567890` + `}` // dummy exp
				return "Bearer " + base64Encode(header) + "." + base64Encode(payload) + ".fakesignature"
			}(),
			expectedCode: http.StatusUnauthorized,
			expectedBody: `{"error":"Invalid or expired token"}`,
			assertNext:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			infrastructure.SetJWTSecret(testSecret)

			w := httptest.NewRecorder()
			r := gin.New()

			nextCalled := false
			r.Use(infrastructure.AuthMiddleware())
			r.GET("/", func(c *gin.Context) {
				nextCalled = true
				if tt.assertNext {
					assert.Equal(t, tt.expectedUser, c.GetString("username"))
					assert.Equal(t, tt.expectedRole, c.GetString("role"))
					c.Status(http.StatusOK)
				}
			})

			req, _ := http.NewRequest("GET", "/", nil)
			req.Header.Set("Authorization", tt.authHeader)
			r.ServeHTTP(w, req)

			if !tt.assertNext {
				assert.Equal(t, tt.expectedCode, w.Code)
				var actualBody map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &actualBody)
				require.NoError(t, err, "Failed to unmarshal response body")

				var expectedBody map[string]string
				err = json.Unmarshal([]byte(tt.expectedBody), &expectedBody)
				require.NoError(t, err, "Failed to unmarshal expected body string")

				assert.Equal(t, expectedBody, actualBody)
				assert.False(t, nextCalled, "next handler should not be called")
			} else {
				assert.True(t, nextCalled, "next handler should be called")
				assert.Equal(t, http.StatusOK, w.Code, "Expected 200 OK status if Next() is called")
				assert.Equal(t, 0, w.Body.Len(), "Response body should be empty if Next() is called by middleware")
			}
		})
	}
}

func base64Encode(s string) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString([]byte(s)), "=")
}

func TestAdminAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		setRole      interface{}
		expectedCode int
		expectedBody string
		assertNext   bool
	}{
		{
			name:         "No role in context",
			setRole:      nil,
			expectedCode: http.StatusForbidden,
			expectedBody: `{"error":"Admin access only"}`,
			assertNext:   false,
		},
		{
			name:         "Role is not admin (regular user)",
			setRole:      "regular",
			expectedCode: http.StatusForbidden,
			expectedBody: `{"error":"Admin access only"}`,
			assertNext:   false,
		},
		{
			name:         "Role is admin",
			setRole:      "admin",
			expectedCode: http.StatusOK,
			expectedBody: "",
			assertNext:   true,
		},
		{
			name:         "Role is wrong type (e.g., int)",
			setRole:      123,
			expectedCode: http.StatusForbidden,
			expectedBody: `{"error":"Admin access only"}`,
			assertNext:   false,
		},
		{
			name:         "Role is empty string",
			setRole:      "",
			expectedCode: http.StatusForbidden,
			expectedBody: `{"error":"Admin access only"}`,
			assertNext:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := gin.New()

			nextCalled := false
			r.Use(func(c *gin.Context) {
				if tt.setRole != nil {
					c.Set("role", tt.setRole)
				}
				c.Next()
			})
			r.Use(infrastructure.AdminAuth())
			r.GET("/", func(c *gin.Context) {
				nextCalled = true
				c.Status(http.StatusOK)
			})

			req, _ := http.NewRequest("GET", "/", nil)
			r.ServeHTTP(w, req)

			if !tt.assertNext {
				assert.Equal(t, tt.expectedCode, w.Code)
				var actualBody map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &actualBody)
				require.NoError(t, err, "Failed to unmarshal response body")

				var expectedBody map[string]string
				err = json.Unmarshal([]byte(tt.expectedBody), &expectedBody)
				require.NoError(t, err, "Failed to unmarshal expected body string")

				assert.Equal(t, expectedBody, actualBody)
				assert.False(t, nextCalled, "next handler should not be called")
			} else {
				assert.True(t, nextCalled, "next handler should be called")
				assert.Equal(t, http.StatusOK, w.Code, "Expected 200 OK status if Next() is called")
				assert.Equal(t, 0, w.Body.Len(), "Response body should be empty if Next() is called")
			}
		})
	}
}
