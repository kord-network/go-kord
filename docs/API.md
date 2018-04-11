# KORD API

The KORD API supports the following actions:

* [`sendPatch`](#sendpatch)
* [`queryGraph`](#querygraph)
* [`queryPatch`](#queryPatch)

See the [Glossary](/glossary) for an overview of KORD definitions.

## sendPatch

Send a patch to modify the state of an object.

### Parameters

1. `patch` - the patch object

### Result

`true` - the patch was successfully applied

### Errors

1. invalid patch format
2. invalid signature

### Example

Send a patch to create a new object:

```
patch = {
  "@context": "...",
  "@id": "kord://0x1234/patch/1",
  "objectID": "kord://0x1234/alice",
  "patch": [
    {
      "op": "add",
      "path": "",
      "value": {
        "@context": "http://schema.org/",
        "@type": "Person",
        "name": "Alice",
        "email": "alice@example.com"
      }
    },
  ],
  "signature": {
    "type": KORDSignatureV1",
    "creator": "kord://0x1234/keys/1",
    "signatureValue": "0x4567"
  }
}

success, error = sendPatch(patch)
```

Send a patch to update an existing object:

```
patch = {
  "@context": "...",
  "@id": "kord://0x1234/patch/2",
  "objectID": "kord://0x1234/alice",
  "patch": [
    {
      "op": "replace",
      "path": "/email",
      "value": "alice@example.net"
    },
    {
      "op": "add",
      "path": "/birthDate",
      "value": "1980-01-14"
    }
  ],
  "signature": {
    "type": KORDSignatureV1",
    "creator": "kord://0x1234/keys/1",
    "signatureValue": "0x6789"
  }
}

success, error = sendPatch(patch)
```

## queryGraph

Query a KORD graph.

### Parameters

1. `graph` - the KORD graph to query
2. `query` - the Cayley query to run
3. `lang`  - the Cayley query language (default: "graphql")

## Result

`Object` - the query results

### Errors

1. invalid graph ID
2. invalid query
3. invalid language

### Example

```
query = "
nodes(type: <schema:Person>) {
  id
  name
  email
}
"

result, error = queryGraph("0x1234", query, "graphql")

> result:
>
> {
>   "data": {
>     "nodes": [
>       {
>         "id": "kord://0x1234/alice",
>         "name": "Alice",
>         "email": "alice@example.net"
>       }
>     ]
>   }
> }
```

## queryPatch

TODO

## Glossary

### Statement

A statement is an RDF triple which consists of three elements:

```
(subject, predicate, object)
```

Elements of a triple are either URIs or literals.

See [RDF Concepts](https://www.w3.org/TR/rdf-concepts/) for more details.


### Object

A KORD object is a set of statements which all have the same subject URI, and
that subject URI is defined as the object's ID.

Objects are represented using [JSON-LD](https://json-ld.org/), for example:

```
{
  "@context": "http://schema.org/",
  "@id":      "http://example.com/alice",
  "@type":    "Person",
  "name":     "Alice",
  "email":    "alice@example.com"
}
```

### KORD ID

A KORD ID is an Ethereum account.

### Graph

A KORD graph is a collection of objects which all have IDs that have a `kord`
scheme, a single KORD ID as the authority and a unique path to identify the
object.

For example, a set of objects with the following IDs:

```
kord://0x1234/object1
kord://0x1234/object2
kord://0x1234/object3
```

would define a graph containing three objects for the KORD ID `0x1234`.

A graph can be considered an object with ID `kord://0x1234`.

### Patch

A patch is an object which describes changes to an object.

A patch includes a linked data signature from the graph's KORD ID.

```
{
  "@context": "...",
  "@id": "kord://0x1234/patch/1",
  "@type": "KORDPatch"
  "objectID": "kord://0x1234/object1",
  "patch": [
    { "op": "add", "path": "/hello", "value": ["world"] }
  ],
  "signature": {
    "type": KORDSignatureV1",
    "creator": "kord://0x1234/keys/1",
    "signatureValue": "0x4567"
  }
}
```
