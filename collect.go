package main

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"math/big"
	"time"
)

const (
	nodeUrlFlag = "node-url"
)

func CollectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collect",
		Short: "A simple task which connects to web3 provider and pull txes, extracting and collecting pub keys",
		Long: `This script connects to a web3 client and pulls transaction data from the blockchain. 
           In particular, it extracts r,v,s signature components of each transaction and calculates the secp256k1 
           public key associated with the Ethereum account that created the transaction. 
           Collected data are stored in LevelDb as current sub-folder "eth-pubkeys".`,
		Run: func(cmd *cobra.Command, args []string) {

			db, err := OpenDb("./eth-pubkeys")
			if err != nil {
				log.Fatal(err)
			}

			lastProcessedBlock := int64(db.GetLastProcessedBlock())
			log.Println("Last processed block", lastProcessedBlock)

			client, err := ethclient.Dial(viper.GetString(nodeUrlFlag))
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
			log.Printf("Collecting pubkeys till %v block finished", lastNetworkBlock)
		},
	}
	cmd.Flags().String(nodeUrlFlag, "https://mainnet.infura.io", "web3 endpoint")
	_ = viper.BindPFlag(nodeUrlFlag, cmd.Flags().Lookup(nodeUrlFlag))
	return cmd
}
