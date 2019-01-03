package main

import (
	"encoding/binary"
	"github.com/ethereum/go-ethereum/common"
	"github.com/syndtr/goleveldb/leveldb"
)

type Db struct {
	base *leveldb.DB
}

var lastBlockKey = []byte("last_processed_block")

func OpenDb(path string) (Db, error) {
	db, err := leveldb.OpenFile(path, nil)
	return Db{base: db}, err
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

func (db *Db) SavePublicKeysBatch(batch *leveldb.Batch) {
	_ = db.base.Write(batch, nil)
}

func (db *Db) SaveAddressPublicKey(address common.Address, pubkey []byte) {
	_ = db.base.Put(address[:], pubkey, nil)
}

func (db *Db) GetAddressPublicKey(address common.Address) (key []byte) {
	key, _ = db.base.Get(address[:], nil)
	return
}

func (db *Db) IteratePubkeys(execute func(address []byte, key []byte)) {
	iter := db.base.NewIterator(nil, nil)
	for iter.Next() {
		execute(iter.Key(), iter.Value())
	}
	iter.Release()
}
