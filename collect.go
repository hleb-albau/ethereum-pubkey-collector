package main

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"math/big"
	"sync"
	"time"
)

const (
	nodeUrlFlag = "node-url"
	threads     = "threads"
)

// Usage: eth-pub-keys collect --node-url=https://mainnet.infura.io --threads=10
func CollectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collect",
		Short: "A simple task which connects to web3 provider and pull txes, extracting and collecting pub keys",
		Long: `This script connects to a web3 client and pulls transaction data from the blockchain. 
           In particular, it extracts r,v,s signature components of each transaction and calculates the secp256k1 
           public key associated with the Ethereum account that created the transaction. 
           Collected data are stored in LevelDb as current sub-folder "eth-pubkeys".`,
		RunE: func(cmd *cobra.Command, args []string) error {

			db, err := OpenDb("eth-pubkeys")
			if err != nil {
				return err
			}

			lastProcessedBlock := int64(db.GetLastProcessedBlock())
			log.Println("Last processed block", lastProcessedBlock)

			client, err := ethclient.Dial(viper.GetString(nodeUrlFlag))
			if err != nil {
				return err
			}

			header, err := client.HeaderByNumber(ctx, nil)
			if err != nil {
				return err
			}

			lastNetworkBlock := header.Number.Int64()
			threads := viper.GetInt64(threads)
			log.Println("Last network block", lastNetworkBlock)

			for blockNum := lastProcessedBlock; blockNum <= lastNetworkBlock; blockNum += threads {

				var wg sync.WaitGroup
				wg.Add(int(threads))
				for i := int64(0); i < threads; i++ {
					go func(thread int64) {
						downloadAndProcessBlock(blockNum+thread, client, db)
						wg.Done()
					}(i)
				}
				wg.Wait()

				if blockNum-lastProcessedBlock > 100 {
					log.Printf("Processed block %v.", blockNum)
					db.SaveLastProcessedBlock(uint64(blockNum))
					lastProcessedBlock = blockNum
				}
			}
			log.Printf("Collecting pubkeys till %v block finished", lastNetworkBlock)
			return nil
		},
	}
	cmd.Flags().String(nodeUrlFlag, "https://mainnet.infura.io", "web3 endpoint")
	cmd.Flags().Int64(threads, 4, "number of concurrent collectors")
	_ = viper.BindPFlag(nodeUrlFlag, cmd.Flags().Lookup(nodeUrlFlag))
	_ = viper.BindPFlag(threads, cmd.Flags().Lookup(threads))
	return cmd
}

func downloadAndProcessBlock(blockNum int64, client *ethclient.Client, db Db) {

	// loop for retry
	for true {

		block, err := client.BlockByNumber(ctx, big.NewInt(blockNum))

		if err != nil {
			// retry after 5 secs
			log.Printf("Could not download block %v:\n %e", blockNum, err)
			time.Sleep(time.Second * 10)
			continue
		}

		for _, tx := range block.Transactions() {
			// process only first txes for each address
			if tx.Nonce() == 0 {
				address, pubkey := GetPubKey(tx)
				db.SaveAddressPublicKey(address, pubkey)
			}
		}
		break
	}
}
