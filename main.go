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

	ctx := context.Background()

	client, err := ethclient.Dial(os.Getenv("LOCAL_GATEWAY"))
	if err != nil {
		log.Fatalf("could not connect to Ethereum gateway: %v\n", err)
	}
	defer client.Close()

	session := message.NewSession(ctx, "KEYSTORE", "KEYSTOREPASS")

	sess := message.LoadInboxContract(session, client, "ADDRESS")

	log.Println(message.ReadMessage(sess))

}
