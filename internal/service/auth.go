package service

import (
	"MusicService/internal/model"
	"MusicService/internal/repository"
	"MusicService/pkg/jwt"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(user *model.RegisterRequest) (*model.UserResponse, error)
	Login(credentials *model.LoginRequest) (string, error)
}

type authService struct {
	userRepo   repository.UserRepository
	jwtService jwt.JWTService
}

func NewAuthService(userRepo repository.UserRepository, jwtService jwt.JWTService) AuthService {
	return &authService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (s *authService) Register(req *model.RegisterRequest) (*model.UserResponse, error) {
	_, err := s.userRepo.FindByEmail(req.Email)
	if err == nil {
		return nil, errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return &model.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (s *authService) Login(req *model.LoginRequest) (string, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := s.jwtService.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
