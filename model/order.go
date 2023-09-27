package model

import (
	"Patronus/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Order struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID    string             `json:"userId,omitempty" bson:"userId,omitempty"`
	Type      util.OrderType     `json:"type,omitempty" bson:"type,omitempty"`
	Bid       bool               `json:"bid" bson:"bid"`
	Size      float64            `json:"size" bson:"size"`
	Price     float64            `json:"price" bson:"price"`
	Coin      string             `json:"coin" bson:"coin"`
	ToFill    float64            `json:"toFill" bson:"toFill"`
	Status    util.OrderStatus   `json:"status,omitempty" bson:"status,omitempty"`
	CreatedAt time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	FilledAt  time.Time          `json:"filledAt,omitempty" bson:"filledAt,omitempty"`
}

func FilteredOrderResponse(order *Order) Order {
	return Order{
		ID:        order.ID,
		Type:      order.Type,
		Bid:       order.Bid,
		Size:      order.Size,
		Price:     order.Price,
		Coin:      order.Coin,
		ToFill:    order.ToFill,
		Status:    order.Status,
		CreatedAt: order.CreatedAt,
		FilledAt:  order.FilledAt,
	}
}

type Orders []*Order

func (o Orders) Len() int           { return len(o) }
func (o Orders) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o Orders) Less(i, j int) bool { return o[i].CreatedAt.UnixNano() < o[j].CreatedAt.UnixNano() }

func (o *Order) IsFilled() bool {
	return o.ToFill == 0.0
}
