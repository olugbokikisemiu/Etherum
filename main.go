package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"github.com/olugbokikisemiu/EthereumDemo/message"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("could not load env: %v\n", err)
	}

	c := context.Background()

	client, err := ethclient.Dial(os.Getenv("GATEWAY"))
	if err != nil {
		log.Fatalf("could not connect to Ethereum gateway: %v\n", err)
	}
	defer client.Close()

	keyStore := keystore.NewKeyStore("node1/keystore", keystore.LightScryptN, keystore.LightScryptP)

	inboxSession := &message.InboxSession{
		Ctx:      c,
		Client:   client,
		Local:    false,
		Keystore: keyStore,
	}

	inboxSession.Session = inboxSession.NewSession()

	fmt.Println("Accounts ", inboxSession.Keystore.Accounts()[0].Address.Hex())

	inboxSession.LoadInboxContract()

	log.Println(inboxSession.ReadMessage())
}
