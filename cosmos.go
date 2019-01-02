package main

import (
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/tendermint/btcd/btcec"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/bech32"

	"github.com/spf13/cobra"
	"log"
)

const (
	addressFlag   = "address"
	accPrefixFlag = "acc-prefix"
)

func CosmosAddressCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "cosmos-address",
		Short: "Calculates for given eth address cosmos-based chain address",
		Run: func(cmd *cobra.Command, args []string) {

			db, err := OpenDb("/home/hlb/projects/eth-pub-keys/eth-pubkeys")
			if err != nil {
				log.Fatal(err)
			}

			ethAddrHex := viper.GetString(addressFlag)

			if !common.IsHexAddress(ethAddrHex) {
				panic(errors.New("ETH address provided in wrong format"))
			}

			ethAddr := common.HexToAddress(ethAddrHex)
			ethRawPubkey := db.GetAddressPublicKey(ethAddr)

			if ethRawPubkey == nil {
				panic(errors.New("No public key found for provided address"))
			}

			cosmosAddr := CosmosAddressFromEthKey(ethRawPubkey)
			log.Printf("[Eth: %s] [Cosmos: %s]", ethAddrHex, EncodeToHex(cosmosAddr, viper.GetString(accPrefixFlag)))
		},
	}
	cmd.Flags().String(addressFlag, "", "hex encoded eth address")
	cmd.Flags().String(accPrefixFlag, "cosmos", "cosmos-based chain acc prefix")

	_ = viper.BindPFlag(addressFlag, cmd.Flags().Lookup(addressFlag))
	_ = viper.BindPFlag(accPrefixFlag, cmd.Flags().Lookup(accPrefixFlag))
	return cmd
}

func CosmosAddressFromEthKey(ethRawPubkey []byte) types.AccAddress {

	var ethCompressedPubkey [33]byte
	ethPubkey, _ := btcec.ParsePubKey(ethRawPubkey[:], btcec.S256())
	copy(ethCompressedPubkey[:], ethPubkey.SerializeCompressed()[:])

	cbdPubKey := secp256k1.PubKeySecp256k1(ethCompressedPubkey)
	cbdAddr := types.AccAddress(cbdPubKey.Address())
	return cbdAddr
}

func EncodeToHex(address types.AccAddress, accPrefix string) string {
	bech32Addr, err := bech32.ConvertAndEncode(accPrefix, address.Bytes())
	if err != nil {
		panic(err)
	}
	return bech32Addr
}
