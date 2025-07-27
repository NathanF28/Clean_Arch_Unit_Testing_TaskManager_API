package mongo

import (
	"context"
	"fmt"
	"task7/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type MongoUserRepository struct { // mongo implementer
	UserCollection *mongo.Collection
}

func NewMongoUserRepository(userCol *mongo.Collection) *MongoUserRepository { // instance of mongo implementer
	return &MongoUserRepository{
		UserCollection: userCol,
	}
}

func (m *MongoUserRepository) RegisterUser(newUser *domain.User) error {
	if newUser.PasswordHash == "" {
		return fmt.Errorf("password cannot be empty")
	}

	filter := bson.M{"username": newUser.Username}
	existingUser := &domain.User{}
	err := m.UserCollection.FindOne(context.TODO(), filter).Decode(existingUser)

	if err == nil {
		return fmt.Errorf("user with username '%s' already exists", newUser.Username)
	}
	if err != mongo.ErrNoDocuments {
		return fmt.Errorf("database error checking for existing user: %w", err)
	}

	hashed_pw, err := bcrypt.GenerateFromPassword([]byte(newUser.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	newUser.PasswordHash = string(hashed_pw)

	count, err := m.UserCollection.CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		return fmt.Errorf("database error counting users: %w", err)
	}

	if count == 0 {
		newUser.Role = "admin"
	} else {
		newUser.Role = "regular"
	}

	newUser.ID = primitive.NewObjectID()
	_, err = m.UserCollection.InsertOne(context.TODO(), newUser)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("user with username '%s' already exists (duplicate key error)", newUser.Username)
		}
		return fmt.Errorf("failed to insert user into database: %w", err)
	}

	return nil
}

func (m *MongoUserRepository) LoginUser(existingUser *domain.User) (domain.User, error) {
	filter := bson.M{"username": existingUser.Username}
	var user domain.User
	err := m.UserCollection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return domain.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(existingUser.PasswordHash))
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (m *MongoUserRepository) PromoteUser(username string) error {
	filter := bson.M{"username": username}
	update := bson.M{"$set": bson.M{"role": "admin"}}
	result, err := m.UserCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}
