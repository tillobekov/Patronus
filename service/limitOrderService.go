package service

import "Patronus/model"

type LimitOrderService interface {
	SaveOrUpdate(order *model.LimitOrderRequestModel) (*model.LimitOrderDBResponseModel, error)
	FindOrderByPriceAndBid(float64, bool) (*model.LimitOrderDBResponseModel, error)
}
