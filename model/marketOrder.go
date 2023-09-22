package model

type MarketOrder struct {
	//Matches     []Match
	Orders      Order
	Bid         bool
	TotalVolume float64
}
