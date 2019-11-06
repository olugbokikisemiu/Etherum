package main

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"github.com/olugbokikisemiu/EthereumDemo/inbox"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("could not load environment: %v\n", err)
	}

	ctx := context.Background()

	client, err := ethclient.Dial(os.Getenv("LOCAL_GATEWAY"))
	if err != nil {
		log.Fatalf("could not connect to Ethereum gateway: %v\n", err)
	}
	defer client.Close()

	privateKey, err := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		log.Fatalf("Unable to convert private key: %v\n", err)
	}
	publicKey := privateKey.Public()

	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatalf("Error occurred: %v\n", err)
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		log.Fatalf("Unable to get nonce: %v\n", err)
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatalf("Gas Price error occurred: %v\n", err)
	}

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = 0
	auth.GasPrice = gasPrice

	address, tx, instance, err := inbox.DeployInbox(auth, client, "Hello World!!!")
	if err != nil {
		log.Fatalf("Deploy error occurred: %v\n", err)
	}

	log.Println("Address: ", address.Hex())

	time.Sleep(30 * time.Second)
	callOpts := &bind.CallOpts{Pending: true, Context: ctx, From: address}

	msg, err := instance.TestMessage(callOpts)
	if err != nil {
		log.Fatalf("Text error occurred: %v\n", err)
	}

	log.Println("tx hash: ", tx.Hash().Hex())
	log.Println("Message: ", msg)

}
