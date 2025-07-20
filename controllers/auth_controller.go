package controllers

import (
	"github.com/gin-gonic/gin"
	"task6/models"
	"task6/data"
	"task6/middleware"
)

func RegisterUser(c *gin.Context) {
	var newUser models.User

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(400, gin.H{"error": "Bad Request"})
		return
	}

	if len(newUser.PasswordHash) < 8 {
		c.JSON(400, gin.H{"error": "Password must be at least 8 characters"})
		return
	}

	err := data.RegisterUser(&newUser)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "User registered successfully"})
}


func LoginUser(c *gin.Context){
	var existingUser models.User
	if err := c.ShouldBindJSON(&existingUser); err != nil {
		c.JSON(400, gin.H{"error": "Bad Request"})
		return
	}
	// authenticate user
	user,err := data.AuthenticateUser(&existingUser)
	if err != nil {
		c.JSON(401, gin.H{"message":"Invalid username or password"})
		return
	}
	// generate jwt token for user

	token, err := middleware.GenerateJWT(user)
	if err != nil {
        c.JSON(500, gin.H{"error": "Could not generate token"})
        return
    }
	c.JSON(200,gin.H{"token": token})
}