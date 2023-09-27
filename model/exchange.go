package model

type Exchange struct {
	OrderBooks map[string]*OrderBook `json:"orderBooks" bson:"orderBooks"`
}

func NewExchange() *Exchange {

	orderBooks := make(map[string]*OrderBook)
	orderBooks[CoinETH] = NewOrderBook()

	return &Exchange{
		OrderBooks: orderBooks,
	}
}
