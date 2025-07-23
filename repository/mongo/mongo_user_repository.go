package repository

import (
	"fmt"
	"context"
	"task7/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)


type MongoUserRepository struct {         // mongo implementer
	UserCollection *mongo.Collection
}

func NewMongoUserRepository(userCol *mongo.Collection) *MongoUserRepository {   // instance of mongo implementer
	return &MongoUserRepository{
		UserCollection : userCol,
	}
}


func (m *MongoUserRepository) RegisterUser(newUser *domain.User) error {
    hashed_pw, err := bcrypt.GenerateFromPassword([]byte(newUser.PasswordHash), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    filter := bson.D{}
    count, err := m.UserCollection.CountDocuments(context.TODO(), filter)
    if err != nil {
        return err
    }

    if count == 0 {
        newUser.Role = "admin"
    } else {
        newUser.Role = "regular"
    }

    newUser.PasswordHash = string(hashed_pw)

    _, err = m.UserCollection.InsertOne(context.TODO(), newUser)
    if err != nil {
        return err
    }

    return nil
}

func (m *MongoUserRepository) LoginUser(existingUser *domain.User) (domain.User, error){
	filter := bson.M{"username" : existingUser.Username}
	var user domain.User
	err := m.UserCollection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return domain.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(existingUser.PasswordHash))
	if err != nil {
		return domain.User{}, err
	}	
	return user,nil
}


func (m *MongoUserRepository) PromoteUser(username string) error {
	filter := bson.M{"username": username}
	update := bson.M{"$set" : bson.M {"role": "admin"}}
	result,err := m.UserCollection.UpdateOne(context.TODO(),filter,update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}