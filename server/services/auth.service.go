package services

import (
	"context"
	"crud_app/models"
	"crud_app/repositories"

	"crud_app/utils"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo repositories.IUserRepository
}

func NewAuthService(userRepo repositories.IUserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) RegisterUser(ctx context.Context, user *models.User) error {
	existingUser, _ := s.userRepo.FindUserByEmail(ctx, user.Email)

	if existingUser != nil {
		return errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return s.userRepo.CreateUser(ctx, user)
}

func (s *AuthService) LoginUser(ctx context.Context, email string, password string) (string, error) {
	user, err := s.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if user == nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := utils.GenerateJWT(*user)
	if err != nil {
		return "", err
	}

	return token, nil
}
