# MusicBrainz META Indexer

The `musicbrainz` package implements a META indexer for the
[MusicBrainz dataset](https://musicbrainz.org/).

## Overview

The `Converter` type is used to read data from a MusicBrainz PostgreSQL
database, convert to META objects and append to a META stream.

The `Indexer` type reads META objects from a stream and indexes them in
a SQLite3 database.

## CLI

To run the conversion on a local `musicbrainz` PostgreSQL database:

```
$ meta musicbrainz convert postgres://localhost:5432/musicbrainz > artists.meta
```

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
