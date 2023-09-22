package service

import "Patronus/model"

type OrderBookService interface {
	FindOrderBookForCoin(symbol *model.CoinSymbol) (*model.OrderBook, error)
	//FindOrderBookById(string) (*model.OrderDBResponseModel, error)
}
