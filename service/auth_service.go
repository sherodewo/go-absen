package service

import (
	"github.com/kreditplus/scorepro/dto"
	"github.com/kreditplus/scorepro/models"
	"github.com/kreditplus/scorepro/repository"
	"github.com/kreditplus/scorepro/utils"
)

type AuthService struct {
	userRepository repository.UserRepository
}

func NewAuthService(repository repository.UserRepository) *AuthService {
	return &AuthService{
		userRepository: repository,
	}
}

func (s *AuthService) FindUserByEmail(email string) (*models.User, error) {
	data, err := s.userRepository.FindWhere(email)
	if err != nil {
		return nil, err
	} else {
	}
	return &data, nil
}

func (s *AuthService) FindUserBranch(id string) (*models.UserDetail, error) {
	data, err := s.userRepository.FindBranch(id)
	if err != nil {
		return nil, err
	} else {
	}
	return &data, nil
}

func (s *AuthService) RegisterUser(dto dto.UserDto) (*models.User, error) {
	hashPassword, _ := utils.HashPassword(dto.Password)
	entity := models.User{
		Nik:      dto.Nik,
		Email:    dto.Email,
		Password: hashPassword,
		IsActive: dto.IsActive,
		TypeUser: dto.TypeUser,
	}
	data, err := s.userRepository.Save(entity)
	if err != nil {
		return nil, err
	} else {
		return &data, nil
	}
}
