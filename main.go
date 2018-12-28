package main

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var ctx = context.Background()

func main() {

	var rootCmd = &cobra.Command{
		Use:   "ethereum-pubkey-collector",
		Short: "Collects Ethereum public keys from signed transactions on the chain.",
		Long: `This script connects to a web3 client and pulls transaction data from the blockchain. Collected keys
			stored as local LevelDb database.`,
	}

	rootCmd.AddCommand(CollectCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
