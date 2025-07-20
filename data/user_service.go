package data

import (
    "fmt"
    "context"
    "task6/models"
    "golang.org/x/crypto/bcrypt"
    "go.mongodb.org/mongo-driver/bson"
)

func RegisterUser(newUser *models.User) error {
    hashed_pw, err := bcrypt.GenerateFromPassword([]byte(newUser.PasswordHash), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    filter := bson.D{}
    count, err := UserCollection.CountDocuments(context.TODO(), filter)
    if err != nil {
        return err
    }

    if count == 0 {
        newUser.Role = "admin"
    } else {
        newUser.Role = "regular"
    }

    newUser.PasswordHash = string(hashed_pw)

    _, err = UserCollection.InsertOne(context.TODO(), newUser)
    if err != nil {
        return err
    }

    return nil
}



func AuthenticateUser (existingUser *models.User) (models.User, error){
	filter := bson.M{"username" : existingUser.Username}
	var user models.User
	err := UserCollection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return models.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(existingUser.PasswordHash))
	if err != nil {
		return models.User{}, err
	}	
	return user,nil
}


func PromoteUser(username string) error {
	filter := bson.M{"username": username}
	update := bson.M{"$set" : bson.M {"role": "admin"}}

	result,err := UserCollection.UpdateOne(context.TODO(),filter,update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}