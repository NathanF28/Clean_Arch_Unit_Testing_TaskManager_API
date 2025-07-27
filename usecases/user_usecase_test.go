package services_test

import (
	"errors"
	"task7/domain"
	services "task7/usecases"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) RegisterUser(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) LoginUser(user *domain.User) (domain.User, error) {
	args := m.Called(user)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserRepository) PromoteUser(username string) error {
	args := m.Called(username)
	return args.Error(0)
}

type UserServiceSuite struct {
	suite.Suite
	mockRepo    *MockUserRepository
	userService services.UserService
}

func (s *UserServiceSuite) SetupTest() {
	s.mockRepo = new(MockUserRepository)
	s.userService = services.NewUserService(s.mockRepo)
}

func TestUserServiceSuite(t *testing.T) {
	suite.Run(t, new(UserServiceSuite))
}

func (s *UserServiceSuite) TestRegisterUser_Success() {
	user := &domain.User{
		Username:     "newuser",
		PasswordHash: "password123",
	}

	s.mockRepo.On("RegisterUser", user).Return(nil).Once()

	err := s.userService.RegisterUser(user)
	s.NoError(err, "RegisterUser should not return an error on success")
	s.mockRepo.AssertExpectations(s.T())
}

func (s *UserServiceSuite) TestRegisterUser_RepositoryError() {
	user := &domain.User{
		Username:     "erroruser",
		PasswordHash: "password123",
	}
	repoError := errors.New("database registration failed")

	s.mockRepo.On("RegisterUser", user).Return(repoError).Once()

	err := s.userService.RegisterUser(user)
	s.Error(err, "RegisterUser should return an error when repository fails")
	s.Equal(repoError, err, "Error returned should be the repository error")
	s.mockRepo.AssertExpectations(s.T())
}

func (s *UserServiceSuite) TestLoginUser_Success() {
	loginCreds := &domain.User{
		Username:     "existinguser",
		PasswordHash: "correctpassword",
	}
	expectedUser := domain.User{
		ID:           primitive.NewObjectID(), // FIX: Changed to generate a new valid ObjectID
		Username:     "existinguser",
		Role:         "regular",
		PasswordHash: "hashedpassword",
	}

	s.mockRepo.On("LoginUser", loginCreds).Return(expectedUser, nil).Once()

	loggedInUser, err := s.userService.LoginUser(loginCreds)
	s.NoError(err, "LoginUser should not return an error on successful login")
	assert.Equal(s.T(), expectedUser.Username, loggedInUser.Username, "Logged in user username should match")
	assert.Equal(s.T(), expectedUser.Role, loggedInUser.Role, "Logged in user role should match")
	s.mockRepo.AssertExpectations(s.T())
}

func (s *UserServiceSuite) TestLoginUser_RepositoryError() {
	loginCreds := &domain.User{
		Username:     "nonexistent",
		PasswordHash: "anypass",
	}
	repoError := errors.New("user not found or invalid credentials")

	s.mockRepo.On("LoginUser", loginCreds).Return(domain.User{}, repoError).Once()

	loggedInUser, err := s.userService.LoginUser(loginCreds)
	s.Error(err, "LoginUser should return an error on repository failure")
	s.Equal(repoError, err, "Error returned should be the repository error")
	s.Equal(domain.User{}, loggedInUser, "Logged in user should be empty on error")
	s.mockRepo.AssertExpectations(s.T())
}

func (s *UserServiceSuite) TestPromoteUser_Success() {
	username := "user_to_promote"

	s.mockRepo.On("PromoteUser", username).Return(nil).Once()

	err := s.userService.PromoteUser(username)
	s.NoError(err, "PromoteUser should not return an error on success")
	s.mockRepo.AssertExpectations(s.T())
}

func (s *UserServiceSuite) TestPromoteUser_RepositoryError() {
	username := "non_existent_user"
	repoError := errors.New("user not found for promotion")

	s.mockRepo.On("PromoteUser", username).Return(repoError).Once()

	err := s.userService.PromoteUser(username)
	s.Error(err, "PromoteUser should return an error when repository fails")
	s.Equal(repoError, err, "Error returned should be the repository error")
	s.mockRepo.AssertExpectations(s.T())
}
