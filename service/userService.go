package service

import "Patronus/model"

type UserService interface {
	FindUserById(string) (*model.UserDBResponseModel, error)
	FindUserByEmail(string) (*model.UserDBResponseModel, error)
}
