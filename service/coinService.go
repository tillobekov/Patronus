package service

import "Patronus/model"

type CoinService interface {
	Save(coin *model.CoinRequestModel) (*model.CoinDBResponseModel, error)
	FindBySymbol(string) (*model.CoinDBResponseModel, error)
}
