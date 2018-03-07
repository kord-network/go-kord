# Go KORD

Golang implementation of the KORD protocol.

## Build

Building the kord binary requires a Go compiler. Once installed, run:

```
go build -o bin/kord ./cmd/kord
```

The binary will now be available in `bin/kord`.

## Development Node

To run a local development node, run:

```
kord node --dev
```

## Testnet

Create a testnet data directory with a single key with an empty passphrase
(just hit enter at the passphrase prompt):

```
$ mkdir testnet
$ geth account new --keystore testnet/keystore
```

Run the node pointing at the testnet data directory:

```
$ kord node --testnet --datadir testnet
```
