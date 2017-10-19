// This file is part of the go-meta library.
//
// Copyright (C) 2017 JAAK MUSIC LTD
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
//
// If you have any questions please contact yo@jaak.io

package musicbrainz

import (
	"context"
	"database/sql"

	"github.com/ethereum/go-ethereum/log"
	"github.com/meta-network/go-meta"
)

// Indexer is a META indexer which indexes a stream of META objects
// representing MusicBrainz Artists into a SQLite3 database, getting the
// associated META objects from a META store.
type Indexer struct {
	index *meta.Index
	store *meta.Store
}

// NewIndexer returns an Indexer which updates the indexes in the given SQLite3
// database connection, getting META objects from the given META store.
func NewIndexer(index *meta.Index, store *meta.Store) (*Indexer, error) {
	// migrate the db to ensure it has an up-to-date schema
	if err := migrations.Run(index.DB); err != nil {
		return nil, err
	}

	return &Indexer{
		index: index,
		store: store,
	}, nil
}

// IndexArtists indexes a stream of META object links which are expected to
// point at MusicBrainz Artists.
func (i *Indexer) IndexArtists(ctx context.Context, stream *meta.StreamReader) error {
	return i.index.Update(func(tx *sql.Tx) error {
		for {
			select {
			case cid, ok := <-stream.Ch():
				if !ok {
					return stream.Err()
				}
				obj, err := i.store.Get(cid)
				if err != nil {
					return err
				}
				artist := &Artist{}
				if err := obj.Decode(artist); err != nil {
					return err
				}
				if err := i.indexArtist(tx, cid.String(), artist); err != nil {
					return err
				}
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	})
}

// indexArtist indexes the given artist on its Name, Type, MBID, IPI and ISNI
// properties.
func (i *Indexer) indexArtist(tx *sql.Tx, cid string, artist *Artist) error {
	log.Info("indexing artist", "id", artist.ID, "name", artist.Name, "mbid", artist.MBID)

	_, err := tx.Exec(
		`INSERT INTO artist (object_id, name, type, mbid) VALUES ($1, $2, $3, $4)`,
		cid, artist.Name, artist.Type, artist.MBID,
	)
	if err != nil {
		return err
	}

	for _, ipi := range artist.IPI {
		_, err := tx.Exec(
			`INSERT INTO artist_ipi (object_id, ipi) VALUES ($1, $2)`,
			cid, ipi,
		)
		if err != nil {
			return err
		}
	}

	for _, isni := range artist.ISNI {
		_, err := tx.Exec(
			`INSERT INTO artist_isni (object_id, isni) VALUES ($1, $2)`,
			cid, isni,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// IndexRecordingWorkLinks indexes a stream of META object links which are expected to
// point at MusicBrainz RecordingWorkLinks.
func (i *Indexer) IndexRecordingWorkLinks(ctx context.Context, stream *meta.StreamReader) error {
	return i.index.Update(func(tx *sql.Tx) error {
		for {
			select {
			case cid, ok := <-stream.Ch():
				if !ok {
					return stream.Err()
				}
				obj, err := i.store.Get(cid)
				if err != nil {
					return err
				}
				var link RecordingWorkLink
				if err := obj.Decode(&link); err != nil {
					return err
				}
				log.Info("indexing recording work link", "isrc", link.ISRC, "iswc", link.ISWC)
				_, err = tx.Exec(
					`INSERT INTO recording_work (object_id, isrc, iswc) VALUES ($1, $2, $3)`,
					cid.String(), link.ISRC, link.ISWC,
				)
				if err != nil {
					return err
				}
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	})
}
