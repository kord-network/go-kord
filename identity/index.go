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
	"database/sql"

	"github.com/ethereum/go-ethereum/log"
	"github.com/meta-network/go-meta"
)

// Indexer is a META indexer which indexes identities and claims
// into a SQLite3 database
type Indexer struct {
	index *meta.Index
}

// NewIndexer returns an Indexer which updates the indexes in the given SQLite3
// database connection
func NewIndexer(index *meta.Index) (*Indexer, error) {
	// migrate the db to ensure it has an up-to-date schema
	if err := migrations.Run(index.DB); err != nil {
		return nil, err
	}

	return &Indexer{
		index: index,
	}, nil
}

// IndexIdentity indexes an identity based on its name , owner and id fields
func (i *Indexer) IndexIdentity(identity *Identity) (err error) {
	return i.index.Update(func(tx *sql.Tx) error {
		log.Info("indexing identity ", "ID ", identity.ID, "OWNER", identity.Owner, "Signature", identity.Sig)

		_, err = tx.Exec(`INSERT INTO identity (id,owner,signature) VALUES ($1, $2, $3)`,
			identity.ID, identity.Owner.String(), identity.Sig)
		return err
	})
}

// IndexClaim indexes a Claim based on its fields
func (i *Indexer) IndexClaim(claim *Claim) (err error) {
	return i.index.Update(func(tx *sql.Tx) error {
		log.Info("indexing claim ", "Issuer ", claim.Issuer, "Subject", claim.Subject, "Claim", claim.Claim, "value", claim.Signature)

		_, err = tx.Exec(`INSERT INTO claim (issuer,subject,claim,signature,id) VALUES ($1, $2, $3, $4, $5)`,
			claim.Issuer, claim.Subject, claim.Claim, claim.Signature, claim.ID)
		return err
	})
}
