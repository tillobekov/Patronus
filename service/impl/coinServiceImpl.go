package impl

import (
	"Patronus/model"
	"Patronus/service"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)

type CoinServiceImpl struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewCoinServiceImpl(collection *mongo.Collection, ctx context.Context) service.CoinService {
	return &CoinServiceImpl{collection: collection, ctx: ctx}
}

func (ms CoinServiceImpl) Save(coin *model.CoinRequestModel) (*model.CoinDBResponseModel, error) {
	coin.CreatedAt = time.Now()

	res, err := ms.collection.InsertOne(ms.ctx, &coin)

	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("the coin with the given symbol already exists")
		}
		return nil, err
	}

	// Create a unique index for the email field
	opt := options.Index()
	opt.SetUnique(true)
	index := mongo.IndexModel{Keys: bson.M{"symbol": 1}, Options: opt}

	if _, err := ms.collection.Indexes().CreateOne(ms.ctx, index); err != nil {
		return nil, errors.New("could not create index for coin")
	}

	var newCoin *model.CoinDBResponseModel
	query := bson.M{"_id": res.InsertedID}

	err = ms.collection.FindOne(ms.ctx, query).Decode(&newCoin)
	if err != nil {
		return nil, err
	}

	return newCoin, nil
}

func (ms CoinServiceImpl) FindBySymbol(symbol string) (*model.CoinDBResponseModel, error) {
	var coinModel *model.CoinDBResponseModel

	query := bson.M{"symbol": strings.ToLower(symbol)}
	err := ms.collection.FindOne(ms.ctx, query).Decode(&coinModel)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &model.CoinDBResponseModel{}, err
		}
		return nil, err
	}

	return coinModel, nil
}
