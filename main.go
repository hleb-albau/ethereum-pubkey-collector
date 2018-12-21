package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

func main() {

	db, err := OpenDb("./eth-pubkeys")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Last processed block %v", db.GetLastProcessedBlock())

	client, err := ethclient.Dial("https://mainnet.infura.io")
	if err != nil {
		log.Fatal(err)
	}

	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(header.Number.String()) // 5671744
}
