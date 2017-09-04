package musicbrainz

import (
	"context"
	"database/sql"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ipfs/go-cid"
	"github.com/meta-network/go-meta"
)

// Indexer is a META indexer which indexes a stream of META objects
// representing MusicBrainz Artists into a SQLite3 database, getting the
// associated META objects from a META store.
type Indexer struct {
	indexDB *sql.DB
	store   *meta.Store
}

// NewIndexer returns an Indexer which updates the indexes in the given SQLite3
// database connection, getting META objects from the given META store.
func NewIndexer(indexDB *sql.DB, store *meta.Store) (*Indexer, error) {
	// migrate the db to ensure it has an up-to-date schema
	if err := migrations.Run(indexDB); err != nil {
		return nil, err
	}

	return &Indexer{
		indexDB: indexDB,
		store:   store,
	}, nil
}

// IndexArtists indexes a stream of META object links which are expected to
// point at MusicBrainz Artists.
func (i *Indexer) IndexArtists(ctx context.Context, stream chan *cid.Cid) error {
	for {
		select {
		case cid, ok := <-stream:
			if !ok {
				return nil
			}
			obj, err := i.store.Get(cid)
			if err != nil {
				return err
			}
			artist := &Artist{}
			if err := obj.Decode(artist); err != nil {
				return err
			}
			if err := i.indexArtist(cid.String(), artist); err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}

// indexArtist indexes the given artist on its Name, Type, MBID, IPI and ISNI
// properties.
func (i *Indexer) indexArtist(cid string, artist *Artist) error {
	log.Info("indexing artist", "id", artist.ID, "name", artist.Name, "mbid", artist.MBID)

	_, err := i.indexDB.Exec(
		`INSERT INTO artist (object_id, name, type, mbid) VALUES ($1, $2, $3, $4)`,
		cid, artist.Name, artist.Type, artist.MBID,
	)
	if err != nil {
		return err
	}

	for _, ipi := range artist.IPI {
		_, err := i.indexDB.Exec(
			`INSERT INTO artist_ipi (object_id, ipi) VALUES ($1, $2)`,
			cid, ipi,
		)
		if err != nil {
			return err
		}
	}

	for _, isni := range artist.ISNI {
		_, err := i.indexDB.Exec(
			`INSERT INTO artist_isni (object_id, isni) VALUES ($1, $2)`,
			cid, isni,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
