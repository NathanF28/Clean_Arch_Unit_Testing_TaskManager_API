package infrastructure

import "task7/domain"

type TokenGenerator interface {
	GenerateToken(user *domain.User) (string, error)
}
