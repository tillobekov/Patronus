package impl

import (
	"Patronus/model"
	"Patronus/service"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type LimitOrderServiceImpl struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewLimitOrderServiceImpl(collection *mongo.Collection, ctx context.Context) service.LimitOrderService {
	return &LimitOrderServiceImpl{collection, ctx}
}

func (los LimitOrderServiceImpl) SaveOrUpdate(order *model.LimitOrderRequestModel) (*model.LimitOrderDBResponseModel, error) {

	fmt.Println("Starting the function...")

	oldLimit, error := los.FindOrderByPriceAndBid(order.Price, order.Bid)
	if oldLimit == nil {
		fmt.Println(error)
		// save new
		order.CreatedAt = time.Now()
		order.UpdatedAt = order.CreatedAt

		res, err := los.collection.InsertOne(los.ctx, &order)

		if err != nil {
			//if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			//	return nil, errors.New("the limit with the given price already exists")
			//}
			fmt.Println("Return 1")
			return nil, err
		}

		// Create a unique index for the email field
		//opt := options.Index()
		//opt.SetUnique(true)
		//index := mongo.IndexModel{Keys: bson.M{"price": 1}, Options: opt}
		//
		//if _, err := los.collection.Indexes().CreateOne(los.ctx, index); err != nil {
		//	return nil, errors.New("could not create index for price")
		//}

		var newLimit *model.LimitOrderDBResponseModel
		query := bson.M{"_id": res.InsertedID}

		err = los.collection.FindOne(los.ctx, query).Decode(&newLimit)

		fmt.Println("Return 2")
		return newLimit, err
	} else {
		fmt.Println(error)
		// update
		oldLimit.UpdatedAt = time.Now()
		oldLimit.Orders = append(oldLimit.Orders, order.Orders[0])
		oldLimit.TotalVolume = oldLimit.TotalVolume + order.Orders[0].Size

		id := oldLimit.ID
		filter := bson.M{"_id": bson.M{"$eq": id}}
		update := bson.M{"$set": bson.M{"updatedAt": oldLimit.UpdatedAt, "orders": oldLimit.Orders, "totalVolume": oldLimit.TotalVolume}}

		_, err := los.collection.UpdateOne(los.ctx, filter, update)
		if err != nil {
			fmt.Println("Return 3")
			return nil, err
		}

		var updatedLimit *model.LimitOrderDBResponseModel
		query := bson.M{"_id": bson.M{"$eq": id}}

		err = los.collection.FindOne(los.ctx, query).Decode(&updatedLimit)

		fmt.Println("Return 4")
		return updatedLimit, err

	}

}
func (los LimitOrderServiceImpl) FindOrderByPriceAndBid(price float64, bid bool) (*model.LimitOrderDBResponseModel, error) {
	var limitModel *model.LimitOrderDBResponseModel

	query := bson.M{"price": price, "bid": bid}
	err := los.collection.FindOne(los.ctx, query).Decode(&limitModel)

	if err != nil {
		//if err == mongo.ErrNoDocuments {
		//	return &model.LimitOrderDBResponseModel{}, err
		//}
		return nil, err
	}

	return limitModel, nil
}
