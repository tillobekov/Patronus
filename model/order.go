package model

import (
	"Patronus/util"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Order struct {
	ID        string    `json:"id" bson:"id"`
	Size      float64   `json:"size" bson:"size"`
	Bid       bool      `json:"bid" bson:"bid"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}

type OrderRequestModel struct {
	ID        string              `json:"id,omitempty" bson:"id,omitempty"`
	User      UserDBResponseModel `json:"user" bson:"user"`
	Type      util.OrderType      `json:"type" bson:"type"`
	Bid       bool                `json:"bid" bson:"bid"`
	Size      float64             `json:"size" bson:"size"`
	Price     float64             `json:"price" bson:"price"`
	Coin      CoinSymbol          `json:"coin" bson:"coin"`
	Tofill    float64             `json:"tofill" bson:"tofill"`
	Status    util.OrderStatus    `json:"status" bson:"status"`
	CreatedAt time.Time           `json:"createdAt" bson:"createdAt"`
	FilledAt  time.Time           `json:"filledAt" bson:"filledAt"`
}

type OrderDBResponseModel struct {
	ID        primitive.ObjectID  `json:"id" bson:"_id"`
	User      UserDBResponseModel `json:"user" bson:"user"`
	Type      util.OrderType      `json:"type" bson:"type"`
	Bid       bool                `json:"bid" bson:"bid"`
	Size      float64             `json:"size" bson:"size"`
	Price     float64             `json:"price" bson:"price"`
	Coin      CoinSymbol          `json:"coin" bson:"coin"`
	Tofill    float64             `json:"tofill" bson:"tofill"`
	Status    util.OrderStatus    `json:"status" bson:"status"`
	CreatedAt time.Time           `json:"createdAt" bson:"createdAt"`
	FilledAt  time.Time           `json:"filledAt" bson:"filledAt"`
}

type OrderResponseModel struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
	//User      UserDBResponseModel `json:"user" bson:"user"`
	Type      util.OrderType   `json:"type" bson:"type"`
	Bid       bool             `json:"bid" bson:"bid"`
	Size      float64          `json:"size" bson:"size"`
	Price     float64          `json:"price" bson:"price"`
	Coin      CoinSymbol       `json:"coin" bson:"coin"`
	Tofill    float64          `json:"tofill" bson:"tofill"`
	Status    util.OrderStatus `json:"status" bson:"status"`
	CreatedAt time.Time        `json:"createdAt" bson:"createdAt"`
	FilledAt  time.Time        `json:"filledAt" bson:"filledAt"`
}

func FilteredOrderResponse(order *OrderDBResponseModel) OrderResponseModel {
	return OrderResponseModel{
		ID:        order.ID,
		Type:      order.Type,
		Bid:       order.Bid,
		Size:      order.Size,
		Price:     order.Price,
		Coin:      order.Coin,
		Tofill:    order.Tofill,
		Status:    order.Status,
		CreatedAt: order.CreatedAt,
		FilledAt:  order.FilledAt,
	}
}

type Orders []*Order

func (o Orders) Len() int           { return len(o) }
func (o Orders) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o Orders) Less(i, j int) bool { return o[i].CreatedAt.UnixNano() < o[j].CreatedAt.UnixNano() }

func NewOrder(order *OrderDBResponseModel) *Order {
	return &Order{
		ID:        order.ID.Hex(),
		Size:      order.Size,
		Bid:       order.Bid,
		CreatedAt: order.CreatedAt,
	}
}

func (o *Order) String() string {
	return fmt.Sprintf("[size: %.2f]", o.Size)
}

func (o *Order) IsFilled() bool {
	return o.Size == 0.0
}

//type OrderCreateModel struct {
//	owner           primitive.ObjectID `json:"owner" bson:"owner"`
//	base            primitive.ObjectID `json:"base" bson:"to"`
//	Name            string             `json:"name" bson:"name" binding:"required"`
//	Email           string             `json:"email" bson:"email" binding:"required"`
//	Password        string             `json:"password" bson:"password" binding:"required,min=8"`
//	PasswordConfirm string             `json:"passwordConfirm" bson:"passwordConfirm,omitempty" binding:"required"`
//	Role            string             `json:"role" bson:"role"`
//	Verified        bool               `json:"verified" bson:"verified"`
//	CreatedAt       time.Time          `json:"created_at" bson:"created_at"`
//	UpdatedAt       time.Time          `json:"updated_at" bson:"updated_at"`
//}
