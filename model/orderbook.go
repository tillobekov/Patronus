package model

import (
	"sort"
)

//type Match struct {
//	Ask        *Order
//	Bid        *Order
//	SizeFilled float64
//	Price      float64
//}

type OrderBook struct {
	asks []*Limit `json:"asks" bson:"asks"`
	bids []*Limit `json:"bids" bson:"bids"`

	MarketAsks []*Order `json:"marketAsks" bson:"marketAsks"`
	MarketBids []*Order `json:"marketBids" bson:"marketBids"`

	orders map[string]*Order `json:"orders" bson:"orders"`
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		asks:       []*Limit{},
		bids:       []*Limit{},
		MarketAsks: []*Order{},
		MarketBids: []*Order{},
		orders:     make(map[string]*Order),
	}
}

func (ob *OrderBook) FindLimitByOrderID(id string, bid bool) *Limit {
	if bid {
		for _, limit := range ob.bids {
			for _, order := range limit.Orders {
				if order.ID.Hex() == id {
					return limit
				}
			}
		}
	} else {
		for _, limit := range ob.asks {
			for _, order := range limit.Orders {
				if order.ID.Hex() == id {
					return limit
				}
			}
		}
	}
	return nil
}

func (ob *OrderBook) FindLimitByPrice(price float64, bid bool) *Limit {
	if bid {
		for _, limit := range ob.bids {
			if limit.Price == price {
				return limit
			}
		}
	} else {
		for _, limit := range ob.asks {
			if limit.Price == price {
				return limit
			}
		}
	}
	return nil
}

func (ob *OrderBook) PlaceMarketOrder(o *Order) ([]*Order, float64) {
	var filledOrders []*Order
	sizeFilled := o.Size

	if o.Bid {
		if ob.AskTotalVolume() == 0.0 {
			ob.MarketBids = append(ob.MarketBids, o)
			return nil, 0.0
		} else {
			for _, limit := range ob.Asks() {
				filled, size := ob.fillOrders(limit.Orders, o)
				filledOrders = append(filledOrders, filled...)
				limit.TotalVolume -= size

				if o.IsFilled() {
					return filledOrders, sizeFilled
				}

			}
			ob.MarketBids = append(ob.MarketBids, o)
			return filledOrders, sizeFilled - o.ToFill
		}

	} else {
		if ob.BidTotalVolume() == 0.0 {
			ob.MarketAsks = append(ob.MarketAsks, o)
			return nil, 0.0
		} else {
			for _, limit := range ob.Bids() {
				filled, size := ob.fillOrders(limit.Orders, o)
				filledOrders = append(filledOrders, filled...)
				limit.TotalVolume -= size

				if o.IsFilled() {
					return filledOrders, sizeFilled
				}
			}
			ob.MarketAsks = append(ob.MarketAsks, o)
			return filledOrders, sizeFilled - o.ToFill
		}
	}
}

func (ob *OrderBook) PlaceLimitOrder(price float64, o *Order) (*Limit, []*Order) {
	var filledOrders []*Order

	if o.Bid {
		filled, _ := ob.fillOrders(ob.MarketAsks, o)
		filledOrders = append(filledOrders, filled...)
	} else {
		filled, _ := ob.fillOrders(ob.MarketBids, o)
		filledOrders = append(filledOrders, filled...)
	}

	var limit *Limit
	limit = ob.FindLimitByPrice(price, o.Bid)

	if limit == nil {
		limit = NewLimit(price)
		if o.Bid {
			ob.bids = append(ob.bids, limit)
		} else {
			ob.asks = append(ob.asks, limit)
		}
	}
	ob.orders[o.ID.Hex()] = o
	limit.AddOrder(o)

	return limit, filledOrders
}

func (ob *OrderBook) fillOrders(orders []*Order, o *Order) ([]*Order, float64) {
	var filledOrders []*Order
	sizeFilled := 0.0

	for _, order := range orders {
		if order.Size <= o.Size {
			o.ToFill -= order.Size
			sizeFilled += order.Size
			order.ToFill = 0.0
		} else {
			order.ToFill -= o.Size
			sizeFilled += o.Size
			o.ToFill = 0.0
		}
		filledOrders = append(filledOrders, order)
		if o.IsFilled() {
			break
		}
	}
	return filledOrders, sizeFilled
}

func (ob *OrderBook) clearLimit(bid bool, l *Limit) {
	if bid {
		for i := 0; i < len(ob.bids); i++ {
			if ob.bids[i] == l {
				ob.bids[i] = ob.bids[len(ob.bids)-1]
				ob.bids = ob.bids[:len(ob.bids)-1]
			}
		}
	} else {
		for i := 0; i < len(ob.asks); i++ {
			if ob.asks[i] == l {
				ob.asks[i] = ob.asks[len(ob.asks)-1]
				ob.asks = ob.asks[:len(ob.asks)-1]
			}
		}
	}
}

func (ob *OrderBook) ClearFilled() {
	for _, l := range ob.bids {
		if l.TotalVolume == 0.0 {
			ob.clearLimit(true, l)
		} else {
			for _, o := range l.Orders {
				if o.Size == 0.0 {
					l.DeleteOrder(o)
				}
			}
		}
	}
	for _, l := range ob.asks {
		if l.TotalVolume == 0.0 {
			ob.clearLimit(false, l)
		} else {
			for _, o := range l.Orders {
				if o.Size == 0.0 {
					l.DeleteOrder(o)
				}
			}
		}
	}
	for i := 0; i < len(ob.MarketBids); i++ {
		if ob.MarketBids[i].Size == 0.0 {
			ob.MarketBids[i] = ob.MarketBids[len(ob.MarketBids)-1]
			ob.MarketBids = ob.MarketBids[:len(ob.MarketBids)-1]
		}
	}
	for i := 0; i < len(ob.MarketAsks); i++ {
		if ob.MarketAsks[i].Size == 0.0 {
			ob.MarketAsks[i] = ob.MarketAsks[len(ob.MarketAsks)-1]
			ob.MarketAsks = ob.MarketAsks[:len(ob.MarketAsks)-1]
		}
	}
	for id, order := range ob.orders {
		if order.Size == 0.0 {
			delete(ob.orders, id)
		}
	}
}

func (ob *OrderBook) CancelOrder(id string, bid bool) {
	limit := ob.FindLimitByOrderID(id, bid)
	limit.DeleteOrder(ob.orders[id])
	if limit.Price == 0.0 {
		ob.clearLimit(bid, limit)
	}

}

func (ob *OrderBook) BidTotalVolume() float64 {
	totalVolume := 0.0

	for i := 0; i < len(ob.bids); i++ {
		totalVolume += ob.bids[i].TotalVolume
	}
	return totalVolume
}

func (ob *OrderBook) AskTotalVolume() float64 {
	totalVolume := 0.0

	for i := 0; i < len(ob.asks); i++ {
		totalVolume += ob.asks[i].TotalVolume
	}
	return totalVolume
}

func (ob *OrderBook) Asks() []*Limit {
	sort.Sort(ByBestAsk{ob.asks})
	return ob.asks
}

func (ob *OrderBook) Bids() []*Limit {
	sort.Sort(ByBestBid{ob.bids})
	return ob.bids
}
