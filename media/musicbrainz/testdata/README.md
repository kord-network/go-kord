# MusicBrainz indexer testdata

This directory contains `artists.json` and `recording-work-links.json` which
contain data exported from a MusicBrainz PostgreSQL database by running the
SQL commands in `artists.sql` and `recording-work-links.sql` and fixing double
quoting:

```
$ psql musicbrainz_db < artists.sql | sed -e 's/\\\\/\\/g' > artists.json

$ psql musicbrainz_db < recording-work-links.sql > recording-work-links.json
```

The explicit artist IDs were chosen as they have either multiple IPI, ISNI,
aliases or annotations, and the recording work link query generates links where
each ISWC maps to 3 ISRCs.
