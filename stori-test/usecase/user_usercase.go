package usecase

import (
	"github.com/ortaman/stori-test/entity"
)

type UserUsecaseRepo struct {
	userRepository entity.UserRepository
}

func NewUserUsecase(r entity.UserRepository) entity.UserUsecase {
	return &UserUsecaseRepo{
		userRepository: r,
	}
}

func (userUsecaseRepo *UserUsecaseRepo) GetOrCreateUser(email string) int {
	userId := userUsecaseRepo.userRepository.GetUserByEmail(email)

	if userId >= 0 {
		return userId
	}

	return userUsecaseRepo.userRepository.CreateUser(email)
}
