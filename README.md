# Go KORD

Golang implementation of the KORD protocol.

## Build

Building the kord binary requires a Go compiler. Once installed, run:

```
go build -o bin/kord ./cmd/kord
```

The binary will now be available in `bin/kord`.

## KORD Node

### Development

To run a local development node, run:

```
kord node --dev
```

### Testnet Node

To run a testnet node, create a testnet data directory with a single key with
an empty passphrase (just hit enter at the passphrase prompt):

```
$ mkdir -p tmp/testnet
$ geth account new --keystore tmp/testnet/keystore
```

Run the node pointing at the testnet data directory:

```
$ kord node --testnet --datadir tmp/testnet
```

## KORD Graphs

Create a KORD ID, entering a passphrase to encrypt the private key:

```
$ kord id new
```

It will output a hex encoded address like `0xba9CA0f65Fb0D4B77ae8c44cCaC7D92EC0D55e88`
which can be used to refer to the KORD ID.

Create a graph:

```
$ kord graph create 0xba9CA0f65Fb0D4B77ae8c44cCaC7D92EC0D55e88
```
