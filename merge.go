package main

import (
	"github.com/hleb-albau/ethereum-pubkey-collector/storage"
	"github.com/spf13/cobra"
)

// Usage: eth-pub-keys merge-dbs from to
// Usage: eth-pub-keys collector-6350 eth-pubkeys
func MergeDbsCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "merge-dbs",
		Short: "Copy all public keys from one db into other",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			fromDb, err := storage.OpenDb(args[0])
			if err != nil {
				return err
			}

			toDb, err := storage.OpenDb(args[1])
			if err != nil {
				return err
			}

			toDb.MergeDbs(&fromDb)
			return nil
		},
	}
	return cmd
}
