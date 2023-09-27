package blockchain

import "Patronus/model"

type Manager interface {
	GetBalance(address string) string
	CreateNewWallet() *model.Wallet
	SendTransaction(from string, to string, valueStr string) bool
}
