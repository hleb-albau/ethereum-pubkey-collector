package main

import (
	"encoding/binary"
	"github.com/ethereum/go-ethereum/common"
	"github.com/syndtr/goleveldb/leveldb"
)

type Db struct {
	base  *leveldb.DB
	index map[string]struct{}
}

var lastBlockKey = []byte("last_processed_block")

func OpenDb(path string) (Db, error) {

	db, err := leveldb.OpenFile(path, nil)
	if err == nil {
		index := make(map[string]struct{})
		iter := db.NewIterator(nil, nil)
		for iter.Next() {
			address := common.BytesToAddress(iter.Key())
			index[address.Hex()] = struct{}{}
		}
		iter.Release()
		return Db{base: db, index: index}, iter.Error()
	}
	return Db{}, err
}

func (db *Db) GetLastProcessedBlock() uint64 {

	value, err := db.base.Get(lastBlockKey, nil)
	if err != nil {
		return 0
	}
	return binary.LittleEndian.Uint64(value)
}

func (db *Db) SaveLastProcessedBlock(num uint64) {

	numbAsBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(numbAsBytes, num)
	_ = db.base.Put(lastBlockKey, numbAsBytes, nil)
}

func (db *Db) SaveAddressPublicKey(address common.Address, pubkey []byte) {
	_ = db.base.Put(address[:], pubkey, nil)
	db.index[address.Hex()] = struct{}{}
}

func (db *Db) GetKnownAddressedCount() int {
	return len(db.index)
}
