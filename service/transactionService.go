package service

import "Patronus/model"

type TransactionService interface {
	Save(transaction model.Transaction) (model.Transaction, error)
}
