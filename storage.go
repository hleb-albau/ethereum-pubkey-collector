package main

import (
	"encoding/binary"
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
			addressAsBytes := iter.Key()
			index[string(addressAsBytes)] = struct{}{}
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

	var numbAsBytes []byte
	binary.LittleEndian.PutUint64(numbAsBytes, num)
	_ = db.base.Put(lastBlockKey, numbAsBytes, nil)
}

func (db *Db) SaveAddressPublicKey(address string, key string) {
	if _, knownAddress := db.index[address]; !knownAddress {
		_ = db.base.Put([]byte(address), []byte(key), nil)
	}
}
