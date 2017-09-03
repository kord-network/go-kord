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

## Testing

To run the tests, run:

```
$ go test .
```

## Vendoring

This project currently uses the [govendor](https://github.com/kardianos/govendor)
tool to manage vendored dependencies.
