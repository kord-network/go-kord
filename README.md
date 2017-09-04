# Go META Library

[META protocol](https://github.com/meta-network/docs) implementation in Go.

[![GoDoc](https://godoc.org/github.com/meta-network/go-meta?status.svg)](https://godoc.org/github.com/meta-network/go-meta)
[![CircleCI](https://circleci.com/gh/meta-network/go-meta.svg?style=svg)](https://circleci.com/gh/meta-network/go-meta)

## Usage

The `meta` package encodes sets of properties as META objects, stores them in a
key-value store, and provides a mechanism for traversing through graphs of
those objects.

The `Object` type is an immutable representation of a META object which has a
Content IDentifier (a.k.a. CID, see https://github.com/ipld/cid) and a raw
byte representation which is the IPLD Canonical CBOR format
(see https://github.com/ipld/specs/tree/master/ipld#canonical-format).

To encode a set of properties:

```
obj, err := Encode(map[string]string{
        "key1": "val1",
        "key2": "val2",
        "key3": "val3",
})
```

The object's CID can be retrieved by calling `obj.Cid()` and the raw CBOR
representation with `obj.RawData()`.

Objects can be linked by assigning an object's CID as the value of a property:

```
jane := MustEncode(map[string]string{"name": "Jane"})
john := MustEncode(map[string]string{"name": "John"})
jack := MustEncode(map[string]string{"name": "Jack"})

me := MustEncode(map[string]string{
        "sister": jane.Cid(),
        "children": []*cid.Cid{
                john.Cid(),
                jack.Cid(),
        },
})
```

Encoding the object as JSON is valid IPLD:

```
json.MarshalIndent(me, "", "  ")

{
  "children": [
    {
      "/": "zdpuAoqDTaSJuFifCN1EDYexQhdz3b4WchYNfvcDWHoEMcXii"
    },
    {
      "/": "zdpuAqWtnAKfAG7RPihvYE4c9iR2nK28hHcJYBFqGnMEpSaeS"
    }
  ],
  "sister": {
    "/": "zdpuArJB9DZzwceaB91z5RE9v6ALhn5tkQxhWw9zTEVBqnkd4"
  }
}
```

Objects can be stored and retrieved using a Store object:

```
store := NewStore(datastore.NewMapDatastore())

err := store.Put(obj)

obj, err := store.Get(cid)
```

An object graph can be traversed using a Graph object:

```
graph := NewGraph(store, root)

v, err := graph.Get("some", "path", "through", "the", "graph")
```

## XML

The `xml` directory contains a Go package which can be used to convert XML
documents and schemas into META object graphs.

## CLI

The `cmd/meta` directory contains a Go program which can be used as a command
line tool to perform the following:

* convert XML documents and schemas into META object graphs
* traverse META object graphs and print the result as a JSON encoded string
* start a HTTP server with an API to convert XML and retreive META objects

For simplicity, the CLI currently stores META objects in a `.meta` directory
which it creates in the working directory of the executing process (this will
later be enhanced to support storing objects in decentralised file storage like
[Swarm](http://swarm-gateways.net/bzz:/theswarm.eth/) or
[IPFS](https://ipfs.io/)).

### Build

To build the CLI, run:

```
go build -o bin/meta ./cmd/meta
```

You can then run it by executing `bin/meta`.

### Usage

#### Import XML Schema

Import an XML Schema (xsd) document into the META store:

```
$ meta import xsd ds \
    http://www.w3.org/2000/09/xmldsig# \
    <(curl -fSL https://www.w3.org/TR/2002/REC-xmldsig-core-20020212/xmldsig-core-schema.xsd)
```

which outputs a CID of the resulting root object:

```
INFO [07-25|19:21:26] object created                           cid=zdpuAz5xgov4sKBWFXvXz5h9kLAX2Tqnt2yBD6Ea79Wq3exfu
```

#### Import XML document

Import an XML document using the CID of an object representing the XML Schema
as the context:

```
TODO
```

#### Print a META object

```
$ meta dump zdpuAz5xgov4sKBWFXvXz5h9kLAX2Tqnt2yBD6Ea79Wq3exfu
{"@context":{"CanonicalizationMethod":"ds:CanonicalizationMethod", ..., ,"ds":"http://www.w3.org/2000/09/xmldsig#"}}
```

```
$ meta dump zdpuAz5xgov4sKBWFXvXz5h9kLAX2Tqnt2yBD6Ea79Wq3exfu/@context/X509SKI
"ds:X509SKI"
```

## Testing

To run the tests, run:

```
$ go test ./...
```

## Vendoring

This project currently uses the [govendor](https://github.com/kardianos/govendor)
tool to manage vendored dependencies.
