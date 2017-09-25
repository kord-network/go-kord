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
	"strings"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ipfs/go-cid"
	"github.com/meta-network/go-meta"
)

// Indexer is a META indexer which indexes a stream of META objects
// representing records (represented in cwr files) into a SQLite3 database, getting the
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

// Index indexes a stream of cwr records META objects
func (i *Indexer) Index(ctx context.Context, stream chan *cid.Cid) error {
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
			if err := i.index(obj); err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// index indexes a cwr records based on the record type (NWR/REV/SPU)
func (i *Indexer) index(cwrRecord *meta.Object) error {

	recordType, err := cwrRecord.GetString("record_type")
	if err != nil {
		return err
	}
	if strings.HasPrefix(recordType, "NWR") || strings.HasPrefix(recordType, "REV") {

		registeredWork := &RegisteredWork{}
		if err := cwrRecord.Decode(registeredWork); err != nil {
			return err
		}
		if err := i.indexRegisteredWork(cwrRecord.Cid().String(), registeredWork); err != nil {
			return err
		}
	}
	if strings.HasPrefix(recordType, "SPU") {

		publisherControlledBySubmitter := &PublisherControllBySubmitter{}
		if err := cwrRecord.Decode(publisherControlledBySubmitter); err != nil {
			return err
		}
		if err := i.indexPublisherControlledBySubmiter(cwrRecord.Cid().String(), publisherControlledBySubmitter); err != nil {
			return err
		}
	}

	return nil
}

// indexRegisteredWork indexes the given registeredWork record on its title,iswc,CompositeType and record_type
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

// indexPublisherControlledBySubmiter indexes the given SPU record on its publisher_sequence_n and record_type
// properties.
func (i *Indexer) indexPublisherControlledBySubmiter(cid string, publisherControlledBySubmitter *PublisherControllBySubmitter) error {
	log.Info("indexing publisherControlledBySubmitter ", "publisher_sequence_n", publisherControlledBySubmitter.PublisherSequenceNumber, "Record Type", publisherControlledBySubmitter.RecordType)
	_, err := i.indexDB.Exec(
		`INSERT INTO publisher_control (object_id, publisher_sequence_n, record_type) VALUES ($1, $2, $3)`,
		cid, publisherControlledBySubmitter.PublisherSequenceNumber, publisherControlledBySubmitter.RecordType,
	)
	if err != nil {
		return err
	}
	return nil
}
