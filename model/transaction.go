package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Transaction struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	From      string             `json:"from" bson:"from"`
	To        string             `json:"to" bson:"to"`
	Value     string             `json:"value" bson:"value"`
	Coin      string             `json:"coin" bson:"coin"`
	Confirmed bool               `json:"confirmed" bson:"confirmed"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}
