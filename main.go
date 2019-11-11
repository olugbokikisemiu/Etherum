package main

import (
	"context"
	"log"
	"os"

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

	inboxSession := &message.InboxSession{
		Ctx:    c,
		Client: client,
		Local:  false,
	}

	inboxSession.Session = inboxSession.NewSession()

	inboxSession.LoadInboxContract()

	log.Println(inboxSession.ReadMessage())
}
