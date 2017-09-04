# MusicBrainz indexer testdata

This directory contains `artists.json` which contains data exported from a
MusicBrainz PostgreSQL database by running the commands in `artists.sql` and
fixing double quoting:

```
$ psql musicbrainz_db < artists.sql | sed -e 's/\\\\/\\/g' > artists.json
```

The explicit artist IDs were chosen as they have either multiple IPI, ISNI,
aliases or annotations.
