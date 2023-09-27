package impl

import (
	"Patronus/model"
	"Patronus/service"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type WalletServiceImpl struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewWalletServiceImpl(collection *mongo.Collection, ctx context.Context) service.WalletService {
	return &WalletServiceImpl{collection: collection, ctx: ctx}
}

func (ws WalletServiceImpl) Save(wallet *model.Wallet) (*model.Wallet, error) {
	wallet.CreatedAt = time.Now()
	res, err := ws.collection.InsertOne(ws.ctx, &wallet)

	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("the wallet with the given address already exists")
		}
		return nil, err
	}

	// Create a unique index for the address field
	opt := options.Index()
	opt.SetUnique(true)
	index := mongo.IndexModel{Keys: bson.M{"address": 2}, Options: opt}

	if _, err := ws.collection.Indexes().CreateOne(ws.ctx, index); err != nil {
		return nil, errors.New("could not create index for wallet")
	}

	var newWallet *model.Wallet
	query := bson.M{"_id": res.InsertedID}

	err = ws.collection.FindOne(ws.ctx, query).Decode(&newWallet)
	if err != nil {
		return nil, err
	}

	return newWallet, nil
}

func (ws WalletServiceImpl) FindUserWalletForNetwork(userId string, network string) (*model.Wallet, error) {
	var walletModel *model.Wallet

	query := bson.M{"user": userId, "network": network}
	err := ws.collection.FindOne(ws.ctx, query).Decode(&walletModel)

	if err != nil {
		return nil, err
	}
	return walletModel, nil
}
