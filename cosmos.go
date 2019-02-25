package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hleb-albau/ethereum-pubkey-collector/crypto"
	"github.com/hleb-albau/ethereum-pubkey-collector/storage"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	addressFlag   = "address"
	accPrefixFlag = "acc-prefix"
)

// Usage: eth-pub-keys cosmos-address --address=0x7C4401aE98F12eF6de39aE24cf9fc51f80EBa16B --acc-prefix=cbd
func CosmosAddressCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "cosmos-address",
		Short: "Calculates for given eth address cosmos-based chain address",
		RunE: func(cmd *cobra.Command, args []string) error {

			db, err := storage.OpenDb("eth-pubkeys")
			if err != nil {
				return err
			}

			ethAddrHex := viper.GetString(addressFlag)

			if !common.IsHexAddress(ethAddrHex) {
				return errors.New("ETH address provided in wrong format")
			}

			ethAddr := common.HexToAddress(ethAddrHex)
			ethRawPubkey := db.GetAddressPublicKey(ethAddr)

			if ethRawPubkey == nil {
				return errors.New("No public key found for provided address")
			}

			cosmosAddr := crypto.CosmosAddressFromEthKey(ethRawPubkey)
			fmt.Printf("[Eth: %s] [Cosmos: %s]", ethAddrHex, crypto.EncodeToHex(cosmosAddr, viper.GetString(accPrefixFlag)))
			return nil
		},
	}
	cmd.Flags().String(addressFlag, "", "hex encoded eth address")
	cmd.Flags().String(accPrefixFlag, "cosmos", "cosmos-based chain acc prefix")

	_ = viper.BindPFlag(addressFlag, cmd.Flags().Lookup(addressFlag))
	_ = viper.BindPFlag(accPrefixFlag, cmd.Flags().Lookup(accPrefixFlag))
	return cmd
}
