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

func (os OrderServiceImpl) Save(order *model.Order) (*model.Order, error) {
	order.ToFill = order.Size
	order.Status = util.OrderStatusACTIVE
	order.CreatedAt = time.Now()

	res, err := os.collection.InsertOne(os.ctx, &order)

	if err != nil {
		return nil, err
	}

	var newOrder *model.Order
	query := bson.M{"_id": res.InsertedID}

	err = os.collection.FindOne(os.ctx, query).Decode(&newOrder)

	return newOrder, err

}

func (os OrderServiceImpl) Update(order *model.Order) (*model.Order, error) {
	if order.ToFill == 0.0 {
		order.Status = util.OrderStatusFILLED
	} else {
		order.Status = util.OrderStatusACTIVE
	}
	//primitiveID, _ := primitive.ObjectIDFromHex(order.ID.Hex())
	filter := bson.M{"_id": bson.M{"$eq": order.ID}}
	update := bson.M{"$set": bson.M{"status": order.Status, "toFill": order.ToFill, "filledAt": time.Now()}}

	_, err := os.collection.UpdateOne(os.ctx, filter, update)
	if err != nil {
		return &model.Order{}, err
	}

	var updatedOrder *model.Order
	query := bson.M{"_id": bson.M{"$eq": order.ID}}
	err = os.collection.FindOne(os.ctx, query).Decode(&updatedOrder)

	return updatedOrder, err

}

func (os *OrderServiceImpl) FindAll(userId string) ([]model.Order, error) {
	var orders []model.Order

	query := bson.M{"userId": userId}
	cursor, err := os.collection.Find(os.ctx, query)

	if err != nil {
		return nil, err
	}

	err = cursor.All(os.ctx, &orders)
	return orders, err
}

func (os OrderServiceImpl) Cancel(id string) (*model.Order, error) {
	primitiveID, _ := primitive.ObjectIDFromHex(id)

	filter := bson.M{"_id": bson.M{"$eq": primitiveID}}
	update := bson.M{"$set": bson.M{"status": util.OrderStatusCANCELED}}

	_, err := os.collection.UpdateOne(os.ctx, filter, update)
	if err != nil {
		return &model.Order{}, err
	}

	var updatedOrder *model.Order
	query := bson.M{"_id": bson.M{"$eq": primitiveID}}
	err = os.collection.FindOne(os.ctx, query).Decode(&updatedOrder)

	return updatedOrder, err
}

func (os *OrderServiceImpl) FindOneById(id string) (*model.Order, error) {
	oid, _ := primitive.ObjectIDFromHex(id)

	var order *model.Order

	query := bson.M{"_id": oid}
	err := os.collection.FindOne(os.ctx, query).Decode(&order)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &model.Order{}, err
		}
		return nil, err
	}

	return order, nil

}
