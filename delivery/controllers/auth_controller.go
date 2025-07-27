package controllers

import (
	"task7/domain"
	"task7/infrastructure"
	services "task7/usecases"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	userService    services.UserService
	tokenGenerator infrastructure.TokenGenerator
}

func NewAuthController(us services.UserService, tg infrastructure.TokenGenerator) *AuthController {
	return &AuthController{
		userService:    us,
		tokenGenerator: tg,
	}
}

func (a *AuthController) RegisterUser(c *gin.Context) {
	var newUser domain.User

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(400, gin.H{"error": "Bad Request"})
		return
	}

	if len(newUser.PasswordHash) < 8 {
		c.JSON(400, gin.H{"error": "Password must be at least 8 characters"})
		return
	}

	err := a.userService.RegisterUser(&newUser)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "User registered successfully"})
}

func (a AuthController) LoginUser(c *gin.Context) {
	var existingUser domain.User
	if err := c.ShouldBindJSON(&existingUser); err != nil {
		c.JSON(400, gin.H{"error": "Bad Request"})
		return
	}
	// authenticate user
	user, err := a.userService.LoginUser(&existingUser)
	if err != nil {
		c.JSON(401, gin.H{"message": "Invalid username or password"})
		return
	}
	// generate jwt token for user

	token, err := a.tokenGenerator.GenerateToken(&user)
	if err != nil {
		c.JSON(500, gin.H{"error": "Could not generate token"})
		return
	}
	c.JSON(200, gin.H{"token": token})
}

func (a AuthController) PromoteUser(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	err := a.userService.PromoteUser(req.Username)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "User promoted to admin"})
}
