# MusicBrainz META Indexer

The `musicbrainz` package implements a META indexer for the
[MusicBrainz dataset](https://musicbrainz.org/) and exposes
a GraphQL API for querying the index.

## Overview

The `Converter` type is used to read data from a MusicBrainz PostgreSQL
database, convert to META objects and append to a META stream.

The `Indexer` type reads META objects from a stream and indexes them in
a SQLite3 database.

The `Resolver` type defines GraphQL resolver functions to execute GraphQL
API queries.

## CLI

### Conversion

To run the conversion on a local `musicbrainz` PostgreSQL database:

```
$ meta musicbrainz convert postgres://localhost:5432/musicbrainz > artists.meta
```

### Indexing

To index the META stream stored in `artists.meta` into `artists.db`:

```
$ meta musicbrainz index artists.db < artists.meta
```

You can then query the index with the `sqlite3` CLI and dump the resulting
META objects using `meta dump`, for example searching for "Ludacris":

```
$ sqlite3 artists.db "SELECT object_id FROM artist WHERE name = 'Ludacris'"
zdpuAzxt6R4537XSihpkhkxR59QXcwWBPbuJ2vf3K8SWtPFmi

$ meta dump zdpuAzxt6R4537XSihpkhkxR59QXcwWBPbuJ2vf3K8SWtPFmi/disambiguation_comment
"US American rapper and actor"

$ meta dump zdpuAzxt6R4537XSihpkhkxR59QXcwWBPbuJ2vf3K8SWtPFmi/alias
["Ludicrous","Ludicrus","Luda","\"Chris Brian Bridges\"","\"Christopher Bridges\"","\"Christopher Brian Bridges\""]
```

### GraphQL API

To run the GraphQL API at `http://localhost:5000/musicbrainz/graphql`:

```
$ meta server --musicbrainz-index artists.db
```

Then send GraphQL queries as a POST request with a JSON body with a `query`
key:

```
$ curl \
    -X POST \
    -H "Content-Type: application/json" \
    --data '{"query": "{ artist(name:\"Ludacris\") { name type gender area } }"}' \
    http://localhost:5000/musicbrainz/graphql

{"data":{"artist":[{"name":"Ludacris","type":"Person","gender":"Male","area":"United States"}]}}
```

There is also a browser based GraphQL explorer at `http://localhost:5000/musicbrainz/`.
