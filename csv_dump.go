package main

import (
	"encoding/csv"
	"github.com/hleb-albau/ethereum-pubkey-collector/storage"
	"github.com/spf13/cobra"
	"os"
)

// Usage: eth-pub-keys csv_dump <path>
// Usage: eth-pub-keys csv_dump eth-pubkeys.csv
func DumpToCSVCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "csv_dump <path>",
		Short: "Copy all loaded public keys from db into selected file in CSV format",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			db, err := storage.OpenDb("eth-pubkeys")
			if err != nil {
				return err
			}

			resultFile, err := os.Create(args[0])
			if err != nil {
				return err
			}
			writer := csv.NewWriter(resultFile)
			writer.Comma = ' '
			defer writer.Flush()
			defer resultFile.Close()

			for address, key := range db.GetAddressesPublicKeys() {
				_ = writer.Write([]string{address, key})
			}
			return nil
		},
	}
	return cmd
}
