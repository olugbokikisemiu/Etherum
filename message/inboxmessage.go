package message

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"github.com/olugbokikisemiu/EthereumDemo/inbox"
)

var myenv map[string]string

const envLoc = ".env"

func loadEnv() {
	var err error
	if myenv, err = godotenv.Read(envLoc); err != nil {
		log.Printf("could not load env from %s: %v", envLoc, err)
	}
}

// NewSession creates an inbox session
func NewSession(ctx context.Context, keyStore string, keyPass string) inbox.InboxSession {
	loadEnv()
	keystore, err := os.Open(myenv[keyStore])
	if err != nil {
		log.Fatalf("Cannot load keystore from location %s: %v\n", os.Getenv(keyStore), err)
	}

	defer keystore.Close()

	keypass := myenv[keyPass]
	auth, err := bind.NewTransactor(keystore, keypass)
	if err != nil {
		log.Fatalf("Error occurreed %v\n", err)
	}

	auth.GasLimit = 1000000
	auth.GasPrice = big.NewInt(1)


	return inbox.InboxSession{
		TransactOpts: *auth,
		CallOpts: bind.CallOpts{
			From:    auth.From,
			Context: ctx,
		},
	}
}

// DeployInboxContract deploys inbox contract if none exist
func DeployInboxContract(session inbox.InboxSession, client *ethclient.Client, message string) inbox.InboxSession {
	loadEnv()

	contractAddress, tx, instance, err := inbox.DeployInbox(&session.TransactOpts, client, message)
	if err != nil {
		log.Fatalf("Deployment error %+v", err)
	}

	fmt.Printf("Contract deployed! Wait for tx %s to be confirmed.\n", tx.Hash().Hex())

	session.Contract = instance
	updateEnvFile("ADDRESS", contractAddress.Hex())
	return session
}

// LoadInboxContract load existing contracts
func LoadInboxContract(session inbox.InboxSession, client *ethclient.Client, address string) inbox.InboxSession {
	loadEnv()

	addr := common.HexToAddress(myenv[address])
	instance, err := inbox.NewInbox(addr, client)
	if err != nil {
		log.Fatalf("Error loading contract: %+v", err)
	}
	session.Contract = instance
	return session
}

// ReadMessage reads message passed to contract constructor while deploying
func ReadMessage(session inbox.InboxSession) string {
	msg, err := session.TestMessage()
	if err != nil {
		return err.Error()
	}
	return msg
}

func SetMessage(session inbox.InboxSession, message string) string {
	tx, err := session.SetMessage(message)
	if err != nil {
		return err.Error()
	}
	return tx.Hash().Hex()
}

func updateEnvFile(k string, val string) {
	myenv[k] = val
	err := godotenv.Write(myenv, envLoc)
	if err != nil {
		log.Printf("failed to update %s: %v\n", envLoc, err)
	}
}
