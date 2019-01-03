package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/syndtr/goleveldb/leveldb"
)

// Usage: eth-pub-keys merge-dbs from to
// Usage: eth-pub-keys collector-6350 eth-pubkeys
func MergeDbsCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "merge-dbs",
		Short: "Copy all public keys from one db into other",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			fromDb, err := OpenDb(args[0])
			if err != nil {
				return err
			}

			toDb, err := OpenDb(args[1])
			if err != nil {
				return err
			}

			batch := &leveldb.Batch{}
			fromDb.IteratePubkeys(func(address []byte, key []byte) {
				batch.Put(address, key)
			})
			toDb.SavePublicKeysBatch(batch)
			fmt.Printf("%v keys added to %s ", batch.Len(), args[1])

			return nil
		},
	}
	return cmd
}
