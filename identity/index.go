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

package identity

import (
	"context"
	"database/sql"

	"github.com/ethereum/go-ethereum/log"
	"github.com/meta-network/go-meta"
)

// Indexer is a META indexer which indexes a stream of META objects
// representing records into a SQLite3 database, getting the
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

// Index indexes a stream of META object links which are expected to
// point at IDs.
func (i *Indexer) Index(ctx context.Context, stream *meta.StreamReader) error {
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
				if err := i.indexIdentity(tx, obj); err != nil {
					return err
				}
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	})
}

// IndexClaim indexes a stream of META object links which are expected to
// point at Claims.
func (i *Indexer) IndexClaim(ctx context.Context, stream *meta.StreamReader) error {
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
				if err := i.indexClaim(tx, obj); err != nil {
					return err
				}
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	})
}

// indexIdentity indexes an identity based on its name , owner and id fields
func (i *Indexer) indexIdentity(tx *sql.Tx, identityObject *meta.Object) (err error) {
	obj, err := i.store.Get(identityObject.Cid())
	if err != nil {
		return err
	}
	identity := &Identity{}

	if err := obj.Decode(identity); err != nil {
		return err
	}
	log.Info("indexing identity ", "ID ", identity.ID, "OWNER", identity.Owner, "Signature", identity.Sig)

	_, err = tx.Exec(`INSERT INTO identity (object_id,id,owner,signature) VALUES ($1, $2, $3, $4)`,
		obj.Cid().String(), identity.ID, identity.Owner, identity.Sig)
	return err
}

// indexClaim indexes a Claim based on its fields
func (i *Indexer) indexClaim(tx *sql.Tx, identityObject *meta.Object) (err error) {
	obj, err := i.store.Get(identityObject.Cid())
	if err != nil {
		return err
	}
	claim := &Claim{}

	if err := obj.Decode(claim); err != nil {
		return err
	}
	log.Info("indexing claim ", "Issuer ", claim.Issuer, "Holder", claim.Holder, "Claim", claim.Claim, "value", claim.Signature)

	_, err = tx.Exec(`INSERT INTO claim (object_id,issuer,holder,claim,signature,id) VALUES ($1, $2, $3, $4, $5, $6)`,
		obj.Cid().String(), claim.Issuer, claim.Holder, claim.Claim, claim.Signature, claim.ID)
	return err
}
