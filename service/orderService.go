package service

import (
	"Patronus/model"
)

type OrderService interface {
	Save(order *model.Order) (*model.Order, error)
	Update(order *model.Order) (*model.Order, error)
	FindAll(userId string) ([]model.Order, error)
	Cancel(id string) (*model.Order, error)
	FindOneById(string) (*model.Order, error)
}
