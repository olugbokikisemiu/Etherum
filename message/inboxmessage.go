package message

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"github.com/olugbokikisemiu/EthereumDemo/inbox"
)

var myenv map[string]string

const envLoc = ".env"

// InboxSession struct type of inbox values
type InboxSession struct {
	Ctx      context.Context
	Session  inbox.InboxSession
	Client   *ethclient.Client
	Local    bool
	Message  string
	Keystore *keystore.KeyStore
}

func loadEnv() {
	var err error
	if myenv, err = godotenv.Read(envLoc); err != nil {
		log.Printf("could not load env from %s: %v", envLoc, err)
	}
}

// NewSession creates an inbox Session
func (i *InboxSession) NewSession() inbox.InboxSession {
	loadEnv()
	var auth *bind.TransactOpts
	switch i.Local {
	case true:
		keystore, err := os.Open(myenv["KEYSTORE"])
		if err != nil {
			log.Fatalf("Cannot load keystore from location %s: %v\n", os.Getenv("KEYSTORE"), err)
		}

		defer keystore.Close()

		keypass := myenv["KEYSTOREPASS"]
		auth, err = bind.NewTransactor(keystore, keypass)
		if err != nil {
			log.Fatalf("Error occurreed %v\n", err)
		}
	default:
		privatekey, _ := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
		auth = bind.NewKeyedTransactor(privatekey)
	}

	auth.Nonce = big.NewInt(1)
	auth.GasLimit = 3000000
	auth.GasPrice = big.NewInt(1200000000000)

	return inbox.InboxSession{
		TransactOpts: *auth,
		CallOpts: bind.CallOpts{
			From:    auth.From,
			Context: i.Ctx,
		},
	}
}

// DeployInboxContract deploys inbox contract if none exist
func (i *InboxSession) DeployInboxContract() inbox.InboxSession {
	loadEnv()

	contractAddress, tx, instance, err := inbox.DeployInbox(&i.Session.TransactOpts, i.Client, i.Message)
	if err != nil {
		log.Fatalf("Deployment error %+v", err)
	}

	fmt.Printf("Contract deployed! Wait for tx %s to be confirmed.\n", tx.Hash().Hex())

	i.Session.Contract = instance
	switch i.Local {
	case true:
		updateEnvFile("LOCAL_ADDRESS", contractAddress.Hex())
	default:
		updateEnvFile("ADDRESS", contractAddress.Hex())
	}

	return i.Session
}

// LoadInboxContract load existing contracts
func (i *InboxSession) LoadInboxContract() inbox.InboxSession {
	loadEnv()
	var addr common.Address
	switch i.Local {
	case true:
		addr = common.HexToAddress(myenv["LOCAL_ADDRESS"])
	default:
		addr = common.HexToAddress(myenv["ADDRESS"])
	}
	instance, err := inbox.NewInbox(addr, i.Client)
	if err != nil {
		log.Fatalf("Error loading contract: %+v", err)
	}
	i.Session.Contract = instance
	return i.Session
}

// ReadMessage reads message passed to contract constructor while deploying
func (i *InboxSession) ReadMessage() string {
	msg, err := i.Session.TestMessage()
	if err != nil {
		return err.Error()
	}
	return msg
}

func (i *InboxSession) SetMessage() string {
	tx, err := i.Session.SetMessage(i.Message)
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
