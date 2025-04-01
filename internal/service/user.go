package service

import (
	"MusicService/internal/model"
	"MusicService/internal/repository"
)

type UserService interface {
	GetProfile(userID uint) (*model.UserResponse, error)
	UpdateProfile(userID uint, update *model.UserResponse) (*model.UserResponse, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetProfile(userID uint) (*model.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	return &model.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (s *userService) UpdateProfile(userID uint, update *model.UserResponse) (*model.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	user.Username = update.Username
	user.Email = update.Email

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return &model.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
