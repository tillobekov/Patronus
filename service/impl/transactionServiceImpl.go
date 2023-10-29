package impl

import (
	"Patronus/model"
	"Patronus/service"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type TransactionServiceImpl struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewTransactionServiceImpl(collection *mongo.Collection, ctx context.Context) service.TransactionService {
	return &TransactionServiceImpl{collection, ctx}
}

func (ts *TransactionServiceImpl) Save(transaction *model.Transaction) (*model.Transaction, error) {
	transaction.Confirmed = true
	transaction.CreatedAt = time.Now()
	transaction.UpdatedAt = transaction.CreatedAt

	res, err := ts.collection.InsertOne(ts.ctx, &transaction)

	if err != nil {
		return nil, err
	}

	var newTransaction *model.Transaction
	query := bson.M{"_id": res.InsertedID}

	err = ts.collection.FindOne(ts.ctx, query).Decode(&newTransaction)

	return newTransaction, err
}

func (ts *TransactionServiceImpl) FindAll(userID string) ([]*model.Transaction, error) {
	var transactions []*model.Transaction
	var senders []*model.Transaction
	var receivers []*model.Transaction

	query := bson.M{"senderID": userID}
	cursor, err := ts.collection.Find(ts.ctx, query)

	if err != nil {
		return nil, err
	}
	err = cursor.All(ts.ctx, &senders)

	query = bson.M{"receiverID": userID}
	cursor, err = ts.collection.Find(ts.ctx, query)

	if err != nil {
		return nil, err
	}
	err = cursor.All(ts.ctx, &receivers)

	transactions = append(transactions, senders...)
	transactions = append(transactions, receivers...)

	return transactions, err
}
