package impl

import (
	"Patronus/model"
	"Patronus/service"
	"Patronus/util"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type OrderServiceImpl struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewOrderServiceImpl(collection *mongo.Collection, ctx context.Context) service.OrderService {
	return &OrderServiceImpl{collection, ctx}
}

func (os OrderServiceImpl) Save(order *model.OrderRequestModel) (*model.OrderDBResponseModel, error) {
	order.Tofill = order.Size
	order.Status = util.OrderStatusACTIVE
	order.CreatedAt = time.Now()

	res, err := os.collection.InsertOne(os.ctx, &order)

	if err != nil {
		return nil, err
	}

	var newOrder *model.OrderDBResponseModel
	query := bson.M{"_id": res.InsertedID}

	err = os.collection.FindOne(os.ctx, query).Decode(&newOrder)

	return newOrder, err

}

func (os OrderServiceImpl) Update(order *model.OrderRequestModel) (*model.OrderDBResponseModel, error) {
	if order.Tofill == 0.0 {
		order.Status = util.OrderStatusFILLED
	} else {
		order.Status = util.OrderStatusACTIVE
	}
	primitiveID, _ := primitive.ObjectIDFromHex(order.ID)
	filter := bson.M{"_id": bson.M{"$eq": primitiveID}}
	update := bson.M{"$set": bson.M{"status": order.Status, "tofill": order.Tofill, "filledAt": time.Now()}}

	_, err := os.collection.UpdateOne(os.ctx, filter, update)
	if err != nil {
		return &model.OrderDBResponseModel{}, err
	}

	var updatedOrder *model.OrderDBResponseModel
	query := bson.M{"_id": bson.M{"$eq": primitiveID}}
	err = os.collection.FindOne(os.ctx, query).Decode(&updatedOrder)

	return updatedOrder, err

}

func (os *OrderServiceImpl) FindAll(user model.UserDBResponseModel) ([]model.OrderDBResponseModel, error) {
	var orders []model.OrderDBResponseModel

	query := bson.M{"user": user}
	cursor, err := os.collection.Find(os.ctx, query)

	if err != nil {
		return nil, err
	}

	err = cursor.All(os.ctx, &orders)
	return orders, err
}

func (os OrderServiceImpl) Cancel(id string) (*model.OrderDBResponseModel, error) {
	primitiveID, _ := primitive.ObjectIDFromHex(id)

	filter := bson.M{"_id": bson.M{"$eq": primitiveID}}
	update := bson.M{"$set": bson.M{"status": util.OrderStatusCANCELED}}

	_, err := os.collection.UpdateOne(os.ctx, filter, update)
	if err != nil {
		return &model.OrderDBResponseModel{}, err
	}

	var updatedOrder *model.OrderDBResponseModel
	query := bson.M{"_id": bson.M{"$eq": primitiveID}}
	err = os.collection.FindOne(os.ctx, query).Decode(&updatedOrder)

	return updatedOrder, err
}

func (os *OrderServiceImpl) FindOneById(id string) (*model.OrderDBResponseModel, error) {
	oid, _ := primitive.ObjectIDFromHex(id)

	var order *model.OrderDBResponseModel

	query := bson.M{"_id": oid}
	err := os.collection.FindOne(os.ctx, query).Decode(&order)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &model.OrderDBResponseModel{}, err
		}
		return nil, err
	}

	return order, nil

}
