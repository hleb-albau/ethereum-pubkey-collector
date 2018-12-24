package main

import (
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"time"
)

var ctx = context.Background()

func main() {

	db, err := OpenDb("./eth-pubkeys")
	if err != nil {
		log.Fatal(err)
	}

	lastProcessedBlock := int64(db.GetLastProcessedBlock())
	log.Println("Last processed block", lastProcessedBlock)

	client, err := ethclient.Dial("https://mainnet.infura.io")
	if err != nil {
		log.Fatal(err)
	}

	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	lastNetworkBlock := header.Number.Int64()
	log.Println("Last network block", lastNetworkBlock)

	for blockNum := lastProcessedBlock; blockNum <= lastNetworkBlock; blockNum++ {

		if blockNum%100 == 0 {
			log.Printf("Processing block %v. Collected %v addresses.", blockNum, db.GetKnownAddressedCount())
			db.SaveLastProcessedBlock(uint64(blockNum))
		}

		block, err := client.BlockByNumber(ctx, big.NewInt(blockNum))
		if err != nil {
			// retry after 5 secs
			time.Sleep(time.Second * 5)
			blockNum--
			continue
		}

		for _, tx := range block.Transactions() {
			// process only first txes for each address
			if tx.Nonce() == 0 {
				address, pubkey := GetPubKey(tx, blockNum)
				db.SaveAddressPublicKey(address, pubkey)
			}
		}
	}
}
