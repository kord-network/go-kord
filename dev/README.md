# META Dev

To run META in development, start the node with the `--dev` flag and use the
`dev` directory as the data directory:

```
$ meta node --dev --datadir ./dev
```

Run `deploy.go` from the root of the repository to deploy ENS:

```
$ go run dev/deploy.go
```

It should deploy three contracts:

```
ENS Registry:             0x241be96854Fc2f0172dAA660EE7A14410957C15d
ENS Resolver:             0xD277b08f085121d287878A991e0C496488AAaEc6
FIFS Registrar for .meta: 0xA686b12D350e03C32eb5694Fc552a2f273C43dF5
```

Create a graph:

```
$ meta create test.meta
```

You should see the ENS record was registered and assigned the hash of an empty
graph:

```
$ geth --preload dev/ensutils.js --exec 'getContent("test.meta")' attach ./dev/geth.ipc
"0x46d91774fb564b0787fcaa4feb984153ce5bcf67b7978edd7488d844ef2c0ace"
```

Load some data into the graph:

```
$ meta load graph/data/testdata.nq test.meta
```
