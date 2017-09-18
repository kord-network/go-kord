# CWR META Indexer

The `cwr` package implements a META indexer for a
[CWR Files Format ](http://members.cisac.org/CisacPortal/consulterDocument.do?id=26603) and exposes
a GraphQL API for querying the index.

## Overview

This module convert CWR transactions ,as appears in CWR formatted files, to META objects and index these with
certain keys.

The `Converter` type is used to read a CWR file , convert to META objects and append to a META stream.
It is done in 2 phases:
Phase 1 -
  Convert CWR format files to standard JSON files (CWR-JSON).
  This is done using [CWR-DataApi python project](http://cwr-dataapi.readthedocs.io/en/latest/)
  There is an [updated fixed forked version](https://github.com/orenyodfat/CWR-DataApi) as a git submodule in this project.
  (there is an open pull request against CWR-DataApi original repo .We should switch back when it is merged) 
  See [readme](https://github.com/orenyodfat/CWR-DataApi/blob/develop/README.rst) for python prerequisites & installing .
Phase 2 -
  Extract and convert CWR transactions from CWR-JSON to META objects.

The `Indexer` type reads META objects from a stream and indexes them in
a SQLite3 database.

The `Resolver` type defines GraphQL resolver functions to execute GraphQL
API queries.

## CLI

### Conversion

To run the conversion on a local `cwr` file:

```
$ meta cwr convert <cwrfile> > registeredwork.meta
```

### Indexing

To index the META stream stored in `registeredwork.meta` into `registeredwork.db`:

```
$ meta cwr index registeredwork.db < registeredwork.meta
```

You can then query the index with the `sqlite3` CLI and dump the resulting
META objects using `meta dump`, for example searching for "PUNK CLUB":

```
$ sqlite3 registeredwork.db "SELECT object_id FROM artist WHERE title = 'PUNK CLUB'"
zdpuAoVMEcareeS4TXPr7YAYNztY1ybbvobV8t7XMkzS9rMeq

$ meta dump zdpuAoVMEcareeS4TXPr7YAYNztY1ybbvobV8t7XMkzS9rMeq/iswc
"T0710203705"

### GraphQL API

To run the GraphQL API at `http://localhost:5000/cwr/graphql`:

```
$ meta server --cwr-index registeredwork.db
```

Then send GraphQL queries as a POST request with a JSON body with a `query`
key:

```
$ curl \
    -X POST \
    -H "Content-Type: application/json" \
    --data '{"query": "{ registered_work(title:\"PUNK CLUB\") { iswc } }"}' \
    http://localhost:5000/cwr/graphql

{"data":{"registered_work":[{"iswc":"T0710203705"}]}}
```

There is also a browser based GraphQL explorer at `http://localhost:5000/cwr/`.


