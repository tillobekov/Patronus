package model

import (
	"sort"
)

type Limit struct {
	//Bid         bool    `json:"bid" bson:"bid"`
	Price       float64 `json:"price" bson:"price"`
	Orders      Orders  `json:"orders" bson:"orders"`
	TotalVolume float64 `json:"totalVolume" bson:"totalVolume"`
}

//type LimitOrderRequestModel struct {
//	Bid         bool                   `json:"bid" bson:"bid"`
//	Price       float64                `json:"price" bson:"price"`
//	Orders      []OrderDBResponseModel `json:"orders" bson:"orders"`
//	TotalVolume float64                `json:"totalVolume" bson:"totalVolume"`
//	CreatedAt   time.Time              `json:"createdAt" bson:"createdAt"`
//	UpdatedAt   time.Time              `json:"updatedAt" bson:"updatedAt"`
//}
//
//type LimitOrderDBResponseModel struct {
//	ID          primitive.ObjectID     `json:"id" bson:"_id"`
//	Bid         bool                   `json:"bid" bson:"bid"`
//	Price       float64                `json:"price" bson:"price"`
//	Orders      []OrderDBResponseModel `json:"orders" bson:"orders"`
//	TotalVolume float64                `json:"totalVolume" bson:"totalVolume"`
//	CreatedAt   time.Time              `json:"createdAt" bson:"createdAt"`
//	UpdatedAt   time.Time              `json:"updatedAt" bson:"updatedAt"`
//}
//
//type LimitRequestModel struct {
//}

func NewLimit(price float64) *Limit {
	return &Limit{
		Price:  price,
		Orders: []*Order{},
	}
}

type Limits []*Limit
type ByBestAsk struct{ Limits }

func (a ByBestAsk) Len() int           { return len(a.Limits) }
func (a ByBestAsk) Swap(i, j int)      { a.Limits[i], a.Limits[j] = a.Limits[j], a.Limits[i] }
func (a ByBestAsk) Less(i, j int) bool { return a.Limits[i].Price < a.Limits[j].Price }

type ByBestBid struct{ Limits }

func (b ByBestBid) Len() int           { return len(b.Limits) }
func (b ByBestBid) Swap(i, j int)      { b.Limits[i], b.Limits[j] = b.Limits[j], b.Limits[i] }
func (b ByBestBid) Less(i, j int) bool { return b.Limits[i].Price > b.Limits[j].Price }

func (l *Limit) AddOrder(o *Order) {
	l.Orders = append(l.Orders, o)
	l.TotalVolume += o.Size
}

func (l *Limit) DeleteOrder(o *Order) {
	for i := 0; i < len(l.Orders); i++ {
		if l.Orders[i] == o {
			l.Orders[i] = l.Orders[len(l.Orders)-1]
			l.Orders = l.Orders[:len(l.Orders)-1]
		}
	}
	l.TotalVolume -= o.Size

	sort.Sort(l.Orders)
}
