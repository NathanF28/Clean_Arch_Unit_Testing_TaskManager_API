package mongo_test

import (
	"context"
	"fmt"
	"task7/domain"
	"task7/repository/mongo"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type MongoUserRepositorySuite struct {
	suite.Suite
	mongoClient    *mongodriver.Client
	userCollection *mongodriver.Collection
	userRepo       *mongo.MongoUserRepository
	databaseName   string
}

func TestMongoUserRepositorySuite(t *testing.T) {
	suite.Run(t, new(MongoUserRepositorySuite))
}

func (s *MongoUserRepositorySuite) SetupSuite() {

	s.databaseName = "task7_test_users_db"
	mongoURI := "mongodb://localhost:27017"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongodriver.Connect(ctx, options.Client().ApplyURI(mongoURI))
	s.Require().NoError(err, "Failed to connect to local MongoDB at "+mongoURI)
	s.mongoClient = client

	err = client.Ping(ctx, nil)
	s.Require().NoError(err, "Failed to ping local MongoDB. Is it running?")

	s.userCollection = client.Database(s.databaseName).Collection("users")
	s.userRepo = mongo.NewMongoUserRepository(s.userCollection)

	_, err = s.userCollection.Indexes().CreateOne(context.Background(), mongodriver.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	s.Require().NoError(err, "Failed to create unique index on username")

	fmt.Printf("Connected to local MongoDB. Using database: %s\n", s.databaseName)
}

func (s *MongoUserRepositorySuite) TearDownSuite() {
	fmt.Println("----- Tearing down MongoUserRepositorySuite (Local MongoDB) -----")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if s.mongoClient != nil {
		err := s.mongoClient.Database(s.databaseName).Drop(ctx)
		s.NoError(err, "Failed to drop test database")
	}
	if s.mongoClient != nil {
		err := s.mongoClient.Disconnect(ctx)
		s.NoError(err, "Failed to disconnect MongoDB client")
	}
	fmt.Println("Local MongoDB connection closed and test database dropped.")
}

func (s *MongoUserRepositorySuite) SetupTest() {
	fmt.Println("--- Clearing users collection before test ---")
	_, err := s.userCollection.DeleteMany(context.Background(), bson.D{})
	s.Require().NoError(err, "Failed to clear users collection")
}

func (s *MongoUserRepositorySuite) TestRegisterUser_FirstUserAsAdmin() {
	user := &domain.User{
		Username:     "adminuser",
		PasswordHash: "password123",
	}

	err := s.userRepo.RegisterUser(user)
	s.Require().NoError(err, "Failed to register first user")
	s.Equal("admin", user.Role, "First user should be assigned 'admin' role")

	var foundUser domain.User
	err = s.userCollection.FindOne(context.Background(), bson.M{"username": user.Username}).Decode(&foundUser)
	s.Require().NoError(err, "Registered user not found in DB")
	s.Equal("admin", foundUser.Role, "Role in DB should be 'admin'")
	s.NotEmpty(foundUser.ID, "User ID should be set by MongoDB")
	s.True(bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte("password123")) == nil, "Password should be correctly hashed and verifiable")
}

func (s *MongoUserRepositorySuite) TestRegisterUser_SubsequentUsersAsRegular() {
	adminUser := &domain.User{Username: "initialadmin", PasswordHash: "adminpass"}
	s.Require().NoError(s.userRepo.RegisterUser(adminUser), "Failed to register initial admin user")
	s.Equal("admin", adminUser.Role, "Initial user should be admin")

	regularUser := &domain.User{
		Username:     "regularuser",
		PasswordHash: "regularpass",
	}
	err := s.userRepo.RegisterUser(regularUser)
	s.Require().NoError(err, "Failed to register second user")
	s.Equal("regular", regularUser.Role, "Second user should be assigned 'regular' role")

	var foundUser domain.User
	err = s.userCollection.FindOne(context.Background(), bson.M{"username": regularUser.Username}).Decode(&foundUser)
	s.Require().NoError(err, "Registered regular user not found in DB")
	s.Equal("regular", foundUser.Role, "Role in DB should be 'regular'")
}

func (s *MongoUserRepositorySuite) TestRegisterUser_DuplicateUsername() {
	user1 := &domain.User{Username: "duplicate", PasswordHash: "pass1"}
	s.Require().NoError(s.userRepo.RegisterUser(user1), "Failed to register first user with duplicate username")

	user2 := &domain.User{Username: "duplicate", PasswordHash: "pass2"}
	err := s.userRepo.RegisterUser(user2)
	s.Require().Error(err, "Expected error when registering user with duplicate username")
	s.Contains(err.Error(), "user with username 'duplicate' already exists", "Error message should indicate a duplicate username")
}

func (s *MongoUserRepositorySuite) TestRegisterUser_EmptyPassword() {
	user := &domain.User{
		Username:     "empty_pass_user",
		PasswordHash: "",
	}
	err := s.userRepo.RegisterUser(user)
	s.Error(err, "Expected error for empty password")
	s.Contains(err.Error(), "password cannot be empty", "Error message should indicate password cannot be empty")
}

func (s *MongoUserRepositorySuite) TestLoginUser_Success() {
	password := "securepass"
	user := &domain.User{Username: "loginuser", PasswordHash: password}
	s.Require().NoError(s.userRepo.RegisterUser(user), "Failed to register user for login test setup")

	loginCreds := &domain.User{Username: "loginuser", PasswordHash: password}
	loggedInUser, err := s.userRepo.LoginUser(loginCreds)
	s.Require().NoError(err, "Login should be successful")
	s.Equal(user.Username, loggedInUser.Username, "Logged in username should match")
	s.Equal(user.Role, loggedInUser.Role, "Logged in role should match")
	s.NotEmpty(loggedInUser.ID, "Logged in user should have an ID")
}

func (s *MongoUserRepositorySuite) TestLoginUser_IncorrectPassword() {
	password := "correctpass"
	user := &domain.User{Username: "badpassuser", PasswordHash: password}
	s.Require().NoError(s.userRepo.RegisterUser(user), "Failed to register user for bad password test setup")

	loginCreds := &domain.User{Username: "badpassuser", PasswordHash: "wrongpass"}
	_, err := s.userRepo.LoginUser(loginCreds)
	s.Error(err, "Login should fail with incorrect password")
	s.ErrorIs(err, bcrypt.ErrMismatchedHashAndPassword, "Error should be bcrypt.ErrMismatchedHashAndPassword")
}

func (s *MongoUserRepositorySuite) TestLoginUser_NotFound() {
	loginCreds := &domain.User{Username: "nonexistentuser", PasswordHash: "anypass"}
	_, err := s.userRepo.LoginUser(loginCreds)
	s.Error(err, "Login should fail for non-existent user")
	s.Contains(err.Error(), "mongo: no documents in result", "Error message should indicate no user found")
}

func (s *MongoUserRepositorySuite) TestPromoteUser_RegularToAdmin() {
	adminUser := &domain.User{Username: "initial_admin_for_promote", PasswordHash: "pass"}
	s.Require().NoError(s.userRepo.RegisterUser(adminUser), "Failed to register initial admin")

	regularUser := &domain.User{Username: "promoteme", PasswordHash: "pass"}
	s.Require().NoError(s.userRepo.RegisterUser(regularUser), "Failed to register user to promote")
	s.Equal("regular", regularUser.Role, "User should initially be regular")

	err := s.userRepo.PromoteUser("promoteme")
	s.Require().NoError(err, "Failed to promote user")

	var updatedUser domain.User
	err = s.userCollection.FindOne(context.Background(), bson.M{"username": "promoteme"}).Decode(&updatedUser)
	s.Require().NoError(err, "Promoted user not found in DB")
	s.Equal("admin", updatedUser.Role, "User role should be updated to 'admin'")
}

func (s *MongoUserRepositorySuite) TestPromoteUser_AdminRemainsAdmin() {
	adminUser := &domain.User{Username: "alreadyadmin", PasswordHash: "pass"}
	s.Require().NoError(s.userRepo.RegisterUser(adminUser), "Failed to register admin user")
	s.Equal("admin", adminUser.Role, "User should initially be admin")

	err := s.userRepo.PromoteUser("alreadyadmin")
	s.Require().NoError(err, "Promoting an already admin user should not return an error")

	var updatedUser domain.User
	err = s.userCollection.FindOne(context.Background(), bson.M{"username": "alreadyadmin"}).Decode(&updatedUser)
	s.Require().NoError(err, "Admin user not found in DB after attempted promotion")
	s.Equal("admin", updatedUser.Role, "Admin user role should remain 'admin'")
}

func (s *MongoUserRepositorySuite) TestPromoteUser_NotFound() {
	err := s.userRepo.PromoteUser("nonexistent_user")
	s.Error(err, "Expected error when promoting non-existent user")
	s.Contains(err.Error(), "user not found", "Error message should indicate user not found")
}
