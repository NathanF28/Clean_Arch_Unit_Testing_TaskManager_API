package services

import (
    "task7/domain"
    "task7/repository/interfaces"
)

type UserService interface {
    RegisterUser(user *domain.User) error
    LoginUser(user *domain.User) (domain.User, error)
	PromoteUser(username string) error
} 

type userService struct {   // one type of userService to implement the interface
    userRepo interfaces.UserRepository // can be any db as long as it implements UserRepository interface
}

func NewUserService(repo interfaces.UserRepository) UserService { // object creation , new type implementer 
    return &userService{
        userRepo: repo,
    }
}

func (s *userService) RegisterUser(user *domain.User) error {
    return s.userRepo.RegisterUser(user)
}

func (s *userService) LoginUser(user *domain.User) (domain.User, error) {
    return s.userRepo.LoginUser(user)
}

func (s *userService) PromoteUser (username string) error {
	return s.userRepo.PromoteUser(username)
}