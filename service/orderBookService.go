package service

import "Patronus/model"

type OrderBookService interface {
	FindOrderBookForCoin(symbol *string) (*model.OrderBook, error)
	//FindOrderBookById(string) (*model.OrderDBResponseModel, error)
}
