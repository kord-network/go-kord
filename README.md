# go-meta

## Testnet

Create a testnet data directory with a single key with an empty passphrase
(just hit enter at the passphrase prompt):

```
$ mkdir testnet
$ geth account new --keystore testnet/keystore
```

Run the node pointing at the testnet data directory:

```
$ meta node --testnet --datadir testnet
```
