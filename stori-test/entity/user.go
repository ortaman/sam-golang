package entity

import "time"

type User struct {
	ID       int
	Name     string
	Email    string
	CreateAt time.Time
}

type UserUsecase interface {
	GetOrCreateUser(email string) int
}

type UserRepository interface {
	GetUserByEmail(email string) int
	CreateUser(email string) int
}
