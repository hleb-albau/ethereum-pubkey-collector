# Ethereum Public Key Collector


[![version](https://img.shields.io/github/release/hleb-albau/ethereum-pubkey-collector.svg?style=flat-square)](https://github.com/hleb-albau/ethereum-pubkey-collector/releases/latest)
[![CircleCI](https://img.shields.io/circleci/project/github/hleb-albau/ethereum-pubkey-collector.svg?style=flat-square)](https://circleci.com/gh/hleb-albau/ethereum-pubkey-collector/tree/master)
[![license](https://img.shields.io/github/license/hleb-albau/ethereum-pubkey-collector.svg?style=flat-square)](https://github.com/hleb-albau/ethereum-pubkey-collector/blob/master/LICENSE)
[![LoC](https://tokei.rs/b1/github/hleb-albau/ethereum-pubkey-collector)](https://github.com/hleb-albau/ethereum-pubkey-collector)
[![contributors](https://img.shields.io/github/contributors/hleb-albau/ethereum-pubkey-collector.svg?style=flat-square)](https://github.com/hleb-albau/ethereum-pubkey-collector/graphs/contributors)
[![contribute](https://img.shields.io/badge/contributions-welcome-orange.svg?style=flat-square)](https://github.com/hleb-albau/ethereum-pubkey-collector/graphs/contributors)

Collects Ethereum public keys from signed transactions on the chain. In particular, service connects to a web3 client and pulls transaction data from the blockchain, extracts r,v,s signature components of each transaction and calculates the secp256k1 public key associated with the Ethereum account that created the transaction. Collected keys stored withing local LevelDB database, and thus can be accessed later.

## Installing

You can download already compiled binaries from [release section](https://github.com/hleb-albau/ethereum-pubkey-collector/releases). But, unfortunfortunatelyently, not for all platforms binaries are available (see [issue #4](https://github.com/hleb-albau/ethereum-pubkey-collector/issues/4)). 

You can build executable by your own (Go 1.11+ **required**):
```
git clone git@github.com:hleb-albau/ethereum-pubkey-collector.git
cd  ethereum-pubkey-collector
go build -o eth-pubkeys ./
# than copy eth-pubkeys to your path, or use it as local executable `./eth-pubkeys` instead of `eth-pubkeys`
```

## Usage

To simple start collecting keys use: 
```
eth-pubkeys collect --node-url=ws://127.0.0.1:8546  --threads=10
```

This project, beside collecting pubkeys, supports other commands such as **create csv file**. To see all available commands use:
```
eth-pubkeys help
```

## Issues

If you have any problems with or questions about this program, please contact us
through a [GitHub issue](https://github.com/hleb-albau/ethereum-pubkey-collector/issues).

## Contributing

You are invited to contribute new features, fixes, or updates, large or small;
I am always thrilled to receive pull requests, and do my best to process them
as fast as I can.
