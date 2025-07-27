package controllers_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"task7/delivery/controllers"
	"task7/domain"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) RegisterUser(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

// *** FIX HERE: Correctly dereference the mock return value ***
func (m *MockUserService) LoginUser(user *domain.User) (domain.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return domain.User{}, args.Error(1) // Return zero-value struct on nil/error
	}
	// Assert that the returned mock value is *domain.User, then dereference it
	return *(args.Get(0).(*domain.User)), args.Error(1) // <-- CORRECTED LINE
}

func (m *MockUserService) PromoteUser(username string) error {
	args := m.Called(username)
	return args.Error(0)
}

type MockTokenGenerator struct {
	mock.Mock
}

// *** FIX HERE: This signature must match infrastructure.TokenGenerator interface (which is *domain.User) ***
func (m *MockTokenGenerator) GenerateToken(user *domain.User) (string, error) { // <-- CORRECTED SIGNATURE
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

type AuthControllerTestSuite struct {
	suite.Suite

	mockUserService    *MockUserService
	mockTokenGenerator *MockTokenGenerator

	authController *controllers.AuthController

	ginContext *gin.Context
	recorder   *httptest.ResponseRecorder
}

func (s *AuthControllerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)

	s.recorder = httptest.NewRecorder()

	s.ginContext, _ = gin.CreateTestContext(s.recorder)

	s.mockUserService = new(MockUserService)
	s.mockTokenGenerator = new(MockTokenGenerator)

	// controllers.NewAuthController expects services.UserService and infrastructure.TokenGenerator
	// Our mocks now correctly implement these interfaces.
	s.authController = controllers.NewAuthController(s.mockUserService, s.mockTokenGenerator)
}

func (s *AuthControllerTestSuite) TearDownTest() {
	s.mockUserService.AssertExpectations(s.T())
	s.mockTokenGenerator.AssertExpectations(s.T())
}

func (s *AuthControllerTestSuite) TestRegisterUser_Success() {
	userJSON := `{"username": "testuser", "passwordHash": "securepassword123"}`
	reqBody := bytes.NewBufferString(userJSON)

	req, _ := http.NewRequest(http.MethodPost, "/register", reqBody)
	req.Header.Set("Content-Type", "application/json")

	s.ginContext.Request = req

	s.mockUserService.On("RegisterUser", mock.AnythingOfType("*domain.User")).Return(nil).Once()

	s.authController.RegisterUser(s.ginContext)

	s.Equal(http.StatusCreated, s.recorder.Code)
	s.Contains(s.recorder.Body.String(), `{"message":"User registered successfully"}`)
}

func (s *AuthControllerTestSuite) TestRegisterUser_InvalidJSON() {
	userJSON := `{"username": "testuser", "passwordHash": }`
	reqBody := bytes.NewBufferString(userJSON)

	req, _ := http.NewRequest(http.MethodPost, "/register", reqBody)
	req.Header.Set("Content-Type", "application/json")
	s.ginContext.Request = req

	s.authController.RegisterUser(s.ginContext)

	s.Equal(http.StatusBadRequest, s.recorder.Code)
	s.Contains(s.recorder.Body.String(), `{"error":"Bad Request"}`)
}

func (s *AuthControllerTestSuite) TestRegisterUser_ShortPassword() {
	userJSON := `{"username": "testuser", "passwordHash": "short"}`
	reqBody := bytes.NewBufferString(userJSON)

	req, _ := http.NewRequest(http.MethodPost, "/register", reqBody)
	req.Header.Set("Content-Type", "application/json")
	s.ginContext.Request = req

	s.authController.RegisterUser(s.ginContext)

	s.Equal(http.StatusBadRequest, s.recorder.Code)
	s.Contains(s.recorder.Body.String(), `{"error":"Password must be at least 8 characters"}`)
}

func (s *AuthControllerTestSuite) TestRegisterUser_ServiceError() {
	userJSON := `{"username": "existinguser", "passwordHash": "securepassword123"}`
	reqBody := bytes.NewBufferString(userJSON)

	req, _ := http.NewRequest(http.MethodPost, "/register", reqBody)
	req.Header.Set("Content-Type", "application/json")
	s.ginContext.Request = req

	expectedServiceError := errors.New("user 'existinguser' already exists")
	s.mockUserService.On("RegisterUser", mock.AnythingOfType("*domain.User")).Return(expectedServiceError).Once()

	s.authController.RegisterUser(s.ginContext)

	s.Equal(http.StatusInternalServerError, s.recorder.Code)
	s.Contains(s.recorder.Body.String(), `{"error":"user 'existinguser' already exists"}`)
}

func (s *AuthControllerTestSuite) TestLoginUser_Success() {
	userJSON := `{"username": "testuser", "passwordHash": "correctpassword"}`
	reqBody := bytes.NewBufferString(userJSON)

	req, _ := http.NewRequest(http.MethodPost, "/login", reqBody)
	req.Header.Set("Content-Type", "application/json")
	s.ginContext.Request = req

	// authenticatedUser is a POINTER to domain.User, because that's what mockUserService.LoginUser
	// is set up to return (and it will be dereferenced by the MockUserService.LoginUser method)
	// and what TokenGenerator.GenerateToken expects.
	authenticatedUserPtr := &domain.User{ // <-- This is a pointer
		ID:           primitive.NewObjectID(),
		Username:     "testuser",
		PasswordHash: "hashed_password_from_db",
		Role:         "user",
	}

	// Mock UserService: Expect LoginUser to be called, returning the authenticated user (POINTER) and no error.
	s.mockUserService.On("LoginUser", mock.AnythingOfType("*domain.User")).Return(authenticatedUserPtr, nil).Once() // <-- Returns POINTER

	expectedToken := "mock_jwt_token_for_user123"
	// *** FIX HERE: TokenGenerator.GenerateToken expects *domain.User, so use *domain.User in mock.AnythingOfType ***
	s.mockTokenGenerator.On("GenerateToken", mock.AnythingOfType("*domain.User")).Return(expectedToken, nil).Once() // <-- Expects POINTER

	s.authController.LoginUser(s.ginContext)

	s.Equal(http.StatusOK, s.recorder.Code)
	s.Contains(s.recorder.Body.String(), `{"token":"`+expectedToken+`"}`)
}

func (s *AuthControllerTestSuite) TestLoginUser_InvalidJSON() {
	userJSON := `{"username": "testuser", "passwordHash": }`
	reqBody := bytes.NewBufferString(userJSON)

	req, _ := http.NewRequest(http.MethodPost, "/login", reqBody)
	req.Header.Set("Content-Type", "application/json")
	s.ginContext.Request = req

	s.authController.LoginUser(s.ginContext)

	s.Equal(http.StatusBadRequest, s.recorder.Code)
	s.Contains(s.recorder.Body.String(), `{"error":"Bad Request"}`)
}

func (s *AuthControllerTestSuite) TestLoginUser_AuthenticationFailed() {
	userJSON := `{"username": "wronguser", "passwordHash": "wrongpassword"}`
	reqBody := bytes.NewBufferString(userJSON)

	req, _ := http.NewRequest(http.MethodPost, "/login", reqBody)
	req.Header.Set("Content-Type", "application/json")
	s.ginContext.Request = req

	authError := errors.New("username or password mismatch")
	// Return nil for the user, as login failed.
	s.mockUserService.On("LoginUser", mock.AnythingOfType("*domain.User")).Return(nil, authError).Once()

	s.authController.LoginUser(s.ginContext)

	s.Equal(http.StatusUnauthorized, s.recorder.Code)
	s.Contains(s.recorder.Body.String(), `{"message":"Invalid username or password"}`)
}

func (s *AuthControllerTestSuite) TestLoginUser_TokenGenerationFailed() {
	userJSON := `{"username": "testuser", "passwordHash": "correctpassword"}`
	reqBody := bytes.NewBufferString(userJSON)

	req, _ := http.NewRequest(http.MethodPost, "/login", reqBody)
	req.Header.Set("Content-Type", "application/json")
	s.ginContext.Request = req

	// authenticatedUser is a POINTER to domain.User
	authenticatedUserPtr := &domain.User{ // <-- This is a pointer
		ID:           primitive.NewObjectID(),
		Username:     "testuser",
		PasswordHash: "hashed_password",
		Role:         "user",
	}

	// Mock UserService: Simulate successful login, returning a domain.User pointer.
	s.mockUserService.On("LoginUser", mock.AnythingOfType("*domain.User")).Return(authenticatedUserPtr, nil).Once() // <-- Returns POINTER

	tokenGenError := errors.New("internal server error during token signing")
	// *** FIX HERE: TokenGenerator.GenerateToken expects *domain.User, so use *domain.User in mock.AnythingOfType ***
	s.mockTokenGenerator.On("GenerateToken", mock.AnythingOfType("*domain.User")).Return("", tokenGenError).Once() // <-- Expects POINTER

	s.authController.LoginUser(s.ginContext)

	s.Equal(http.StatusInternalServerError, s.recorder.Code)
	s.Contains(s.recorder.Body.String(), `{"error":"Could not generate token"}`)
}

func (s *AuthControllerTestSuite) TestPromoteUser_Success() {
	promoteJSON := `{"username": "user_to_promote"}`
	reqBody := bytes.NewBufferString(promoteJSON)

	req, _ := http.NewRequest(http.MethodPost, "/promote", reqBody)
	req.Header.Set("Content-Type", "application/json")
	s.ginContext.Request = req

	s.mockUserService.On("PromoteUser", "user_to_promote").Return(nil).Once()

	s.authController.PromoteUser(s.ginContext)

	s.Equal(http.StatusOK, s.recorder.Code)
	s.Contains(s.recorder.Body.String(), `{"message":"User promoted to admin"}`)
}

func (s *AuthControllerTestSuite) TestPromoteUser_InvalidJSON() {
	promoteJSON := `{"username": }`
	reqBody := bytes.NewBufferString(promoteJSON)

	req, _ := http.NewRequest(http.MethodPost, "/promote", reqBody)
	req.Header.Set("Content-Type", "application/json")
	s.ginContext.Request = req

	s.authController.PromoteUser(s.ginContext)

	s.Equal(http.StatusBadRequest, s.recorder.Code)
	s.Contains(s.recorder.Body.String(), `{"error":"Invalid request"}`)
}

func (s *AuthControllerTestSuite) TestPromoteUser_ServiceError() {
	promoteJSON := `{"username": "nonexistent_user"}`
	reqBody := bytes.NewBufferString(promoteJSON)

	req, _ := http.NewRequest(http.MethodPost, "/promote", reqBody)
	req.Header.Set("Content-Type", "application/json")
	s.ginContext.Request = req

	expectedServiceError := errors.New("user 'nonexistent_user' not found")
	s.mockUserService.On("PromoteUser", "nonexistent_user").Return(expectedServiceError).Once()

	s.authController.PromoteUser(s.ginContext)

	s.Equal(http.StatusNotFound, s.recorder.Code)
	s.Contains(s.recorder.Body.String(), `{"error":"user 'nonexistent_user' not found"}`)
}

func TestAuthController(t *testing.T) {
	suite.Run(t, new(AuthControllerTestSuite))
}
