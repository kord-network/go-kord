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

package cwr

import (
	"context"
	"database/sql"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ipfs/go-cid"
	"github.com/meta-network/go-meta"
)

// Indexer is a META indexer which indexes a stream of META objects
// representing registered works (represented in cwr files) into a SQLite3 database, getting the
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

// IndexRegisteredWorks indexes a stream of META object links which are expected to
// point at registered works.
func (i *Indexer) IndexRegisteredWorks(ctx context.Context, stream chan *cid.Cid) error {
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
			registeredWork := &RegisteredWork{}
			if err := obj.Decode(registeredWork); err != nil {
				return err
			}
			if err := i.indexRegisteredWork(cid.String(), registeredWork); err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// indexRegisteredWork indexes the given artist on its title,iswc,CompositeType and record_type
// properties.
func (i *Indexer) indexRegisteredWork(cid string, registeredWork *RegisteredWork) error {
	log.Info("indexing registered work", "title", registeredWork.Title, "iswc", registeredWork.ISWC, "CompositeType", registeredWork.CompositeType, "Record Type", registeredWork.RecordType)
	_, err := i.indexDB.Exec(
		`INSERT INTO registered_work (object_id, title, iswc, composite_type,record_type) VALUES ($1, $2, $3, $4, $5)`,
		cid, registeredWork.Title, registeredWork.ISWC, registeredWork.CompositeType, registeredWork.RecordType,
	)
	if err != nil {
		return err
	}
	return nil
}
