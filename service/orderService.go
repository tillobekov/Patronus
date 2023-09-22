package service

import (
	"Patronus/model"
)

type OrderService interface {
	Save(order *model.OrderRequestModel) (*model.OrderDBResponseModel, error)
	Update(order *model.OrderRequestModel) (*model.OrderDBResponseModel, error)
	FindAll(user model.UserDBResponseModel) ([]model.OrderDBResponseModel, error)
	Cancel(id string) (*model.OrderDBResponseModel, error)
	FindOneById(string) (*model.OrderDBResponseModel, error)
}
