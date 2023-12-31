package service

import "Patronus/model"

type WalletService interface {
	Save(wallet *model.Wallet) (*model.Wallet, error)
	FindUserWalletForNetwork(userId string, network string) (*model.Wallet, error)
}
