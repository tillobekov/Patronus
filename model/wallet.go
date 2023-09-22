package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Wallet struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	User       string             `json:"user" bson:"user"`
	Network    string             `json:"network" bson:"network"`
	Address    string             `json:"address" bson:"address"`
	PrivateKey string             `json:"privateKey" bson:"privateKey"`
	Balance    float64            `json:"balance" bson:"balance"`
	CreatedAt  time.Time          `json:"createdAt" bson:"createdAt"`
}

//type WalletRequestModel struct {
//	User       string  `json:"user" bson:"user"`
//	Network    string  `json:"network" bson:"network"`
//	Address    string  `json:"address" bson:"address"`
//	PrivateKey string  `json:"privateKey" bson:"privateKey"`
//	Balance    float64 `json:"balance" bson:"balance"`
//}

//type WalletDBResponseModel struct {
//	ID         primitive.ObjectID `json:"id" bson:"_id"`
//	User       string             `json:"user" bson:"user"`
//	Network    string             `json:"network" bson:"network"`
//	Address    string             `json:"address" bson:"address"`
//	PrivateKey string             `json:"privateKey" bson:"privateKey"`
//	Balance    float64            `json:"balance" bson:"balance"`
//	CreatedAt  time.Time          `json:"createdAt" bson:"createdAt"`
//}

type WalletResponseModel struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	User       string             `json:"user" bson:"user"`
	Network    string             `json:"network" bson:"network"`
	Address    string             `json:"address" bson:"address"`
	PrivateKey string             `json:"privateKey" bson:"privateKey"`
	Balance    float64            `json:"balance" bson:"balance"`
	CreatedAt  time.Time          `json:"createdAt" bson:"createdAt"`
}

func FilteredWalletResponse(wallet *Wallet) *WalletResponseModel {
	return &WalletResponseModel{
		ID:         wallet.ID,
		User:       wallet.User,
		Network:    wallet.Network,
		Address:    wallet.Address,
		PrivateKey: wallet.PrivateKey,
		Balance:    wallet.Balance,
		CreatedAt:  wallet.CreatedAt}
}
