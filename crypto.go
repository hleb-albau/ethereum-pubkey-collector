package main

import (
	"github.com/ethereum/go-ethereum/common"
	eth "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"math/big"
)

var big8 = big.NewInt(8)

func GetPubKey(tx *eth.Transaction, blockNum int64) (common.Address, []byte) {

	signer := eth.MakeSigner(params.MainnetChainConfig, big.NewInt(blockNum))

	sighash := signer.Hash(tx)
	Vb, R, S := tx.RawSignatureValues()

	// EIP155 support
	if Vb.Int64() > 28 {
		Vb.Sub(Vb, tx.ChainId()).Sub(Vb, tx.ChainId()).Sub(Vb, big8)
	}

	V := byte(Vb.Uint64() - 27)
	r, s := R.Bytes(), S.Bytes()
	sig := make([]byte, 65)
	copy(sig[32-len(r):32], r)
	copy(sig[64-len(s):64], s)
	sig[64] = V

	// recover the public key from the signature
	var addr common.Address
	pubkey, _ := crypto.Ecrecover(sighash[:], sig)
	copy(addr[:], crypto.Keccak256(pubkey[1:])[12:])
	return addr, pubkey
}
