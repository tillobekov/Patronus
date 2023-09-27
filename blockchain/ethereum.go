package blockchain

import (
	"Patronus/model"
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

type EthereumManager struct {
	ctx    context.Context
	client *ethclient.Client
}

func NewEthereumManager(url string) Manager {
	client, err := ethclient.Dial(url)
	if err != nil {
		fmt.Printf("EthereumManager: Error while connecting to a client: %+v\n", err)
		panic(err)
	}
	return &EthereumManager{
		ctx:    context.Background(),
		client: client}
}

func (eth *EthereumManager) GetBalance(address string) string {
	account := common.HexToAddress(address)
	balance, err := eth.client.BalanceAt(eth.ctx, account, nil)
	if err != nil {
		fmt.Printf("EthereumManager: Error while reading the balance: %+v\n", err)
		panic(err)
	}
	//fbalance := new(big.Float)
	//fbalance.SetString(balance.String())
	//ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	return balance.String()
}

func (eth *EthereumManager) CreateNewWallet() *model.Wallet {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		fmt.Printf("EthereumManager: Error while generating a private key: %+v\n", err)
		panic(err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	//fmt.Println(hexutil.Encode(privateKeyBytes)[2:])

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		fmt.Printf("EthereumManager: Error while casting public key to ECDSA: %+v\n", err)
		panic(err)
	}

	//publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	return &model.Wallet{
		Address:    address,
		PrivateKey: hexutil.Encode(privateKeyBytes),
	}
}

//func (eth *EthereumManager) CreateNewAccount() {
//	ks := keystore.NewKeyStore("./wallets", keystore.StandardScryptN, keystore.StandardScryptP)
//	password := "secret"
//	account, err := ks.NewAccount(password)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Println(account.Address.Hex())
//	fmt.Println(eth.GetBalance(account.Address.Hex()))
//}

func (eth *EthereumManager) SendTransaction(from string, to string, valueStr string) bool {
	privateKey, err := crypto.HexToECDSA(from)
	if err != nil {
		fmt.Printf("EthereumManager: Error while sending a transaction (1): %+v\n", err)
		panic(err)
		return false
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		fmt.Printf("EthereumManager: Error while sending a transaction (casting public key to ECDSA): %+v\n", err)
		panic(err)
		return false
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := eth.client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		fmt.Printf("EthereumManager: Error while sending a transaction (3): %+v\n", err)
		panic(err)
		return false
	}

	gasLimit := uint64(21000) // in units
	gasPrice, err := eth.client.SuggestGasPrice(context.Background())
	if err != nil {
		fmt.Printf("EthereumManager: Error while sending a transaction (4): %+v\n", err)
		panic(err)
		return false
	}

	value := new(big.Int)
	value, ok = value.SetString(valueStr, 10)
	if !ok {
		fmt.Printf("EthereumManager: Error while sending a transaction (5): %+v\n", err)
		panic(err)
		return false
	}

	toAddress := common.HexToAddress(to)
	var data []byte
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	chainID, err := eth.client.NetworkID(context.Background())
	if err != nil {
		fmt.Printf("EthereumManager: Error while sending a transaction (6): %+v\n", err)
		panic(err)
		return false
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		fmt.Printf("EthereumManager: Error while sending a transaction (7): %+v\n", err)
		panic(err)
		return false
	}

	err = eth.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		fmt.Printf("EthereumManager: Error while sending a transaction (8): %+v\n", err)
		panic(err)
		return false
	}

	//fmt.Printf("tx sent: %s\n", signedTx.Hash().Hex())
	return true

}
