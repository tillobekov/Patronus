package model

type Exchange struct {
	OrderBooks map[CoinSymbol]*OrderBook `json:"orderBooks" bson:"orderBooks"`
}

func NewExchange() *Exchange {

	orderBooks := make(map[CoinSymbol]*OrderBook)
	orderBooks[CoinETH] = NewOrderBook()

	return &Exchange{
		OrderBooks: orderBooks,
	}
}
