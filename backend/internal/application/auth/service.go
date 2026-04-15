package auth

import "dangde-world/backend/internal/domain/user"

type Service struct {
	users user.Repository
}

func NewService(users user.Repository) Service {
	return Service{users: users}
}

func (s Service) Login(userID uint) (user.User, error) {
	return s.users.FindByID(userID)
}
