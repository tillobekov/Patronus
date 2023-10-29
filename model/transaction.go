package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Transaction struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	SenderID   string             `json:"senderID" bson:"senderID"`
	ReceiverID string             `json:"receiverID" bson:"receiverID"`
	//OrderID    string             `json:"order" bson:"order"`
	Value     float64   `json:"value" bson:"value"`
	Coin      string    `json:"coin" bson:"coin"`
	Confirmed bool      `json:"confirmed" bson:"confirmed"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}
