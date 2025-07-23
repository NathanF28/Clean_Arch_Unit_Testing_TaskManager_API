package interfaces

import (
	"task7/domain"
)

type UserRepository interface {            // choose any db that implements register and login
	RegisterUser(user *domain.User) error
	LoginUser(user *domain.User) (domain.User,error)
	PromoteUser(username string) error
}





