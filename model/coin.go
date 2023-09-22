package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type CoinSymbol string

const (
	CoinETH CoinSymbol = "ETH"
)

type CoinRequestModel struct {
	Symbol    string    `json:"symbol" bson:"symbol"`
	Name      string    `json:"name" bson:"name"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}

type CoinDBResponseModel struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Symbol    string             `json:"symbol" bson:"symbol"`
	Name      string             `json:"name" bson:"name"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
}
