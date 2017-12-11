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
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/meta-network/go-meta"
)

type Index struct {
	*meta.Index
}

func NewIndex(store *meta.Store) (*Index, error) {
	index, err := store.OpenIndex("id.meta")
	if err != nil {
		return nil, err
	}
	// migrate the db to ensure it has an up-to-date schema
	if err := migrations.Run(index.DB); err != nil {
		index.Close()
		return nil, err
	}
	return &Index{index}, nil
}

func (i *Index) CreateIdentity(identity *Identity) error {
	if !verifyIdentity(identity) {
		return errors.New("identity: invalid identity")
	}
	return i.Update(func(tx *sql.Tx) error {
		_, err := tx.Exec(
			"INSERT INTO identity (id, username, owner, signature) VALUES ($1, $2, $3, $4)",
			identity.ID().String(),
			identity.Username,
			identity.Owner.String(),
			hexutil.Encode(identity.Signature),
		)
		return err
	})
}

func (i *Index) Identity(id string) (*Identity, error) {
	var (
		username  string
		owner     string
		signature string
	)
	row := i.QueryRow("SELECT username, owner, signature FROM identity WHERE id = ?", id)
	if err := row.Scan(&username, &owner, &signature); err != nil {
		return nil, err
	}
	return &Identity{
		Username:  username,
		Owner:     common.HexToAddress(owner),
		Signature: common.FromHex(signature),
	}, nil
}

type IdentityFilter struct {
	ID       *string
	Username *string
	Owner    *string
}

func (i *Index) Identities(filter IdentityFilter) ([]*Identity, error) {
	where := make([]string, 0, 3)
	params := make([]interface{}, 0, 3)
	if filter.ID != nil {
		where = append(where, "id = ?")
		params = append(params, *filter.ID)
	}
	if filter.Username != nil {
		where = append(where, "username = ?")
		params = append(params, *filter.Username)
	}
	if filter.Owner != nil {
		where = append(where, "owner = ?")
		params = append(params, *filter.Owner)
	}
	if len(where) == 0 {
		return nil, errors.New("missing id, username or owner")
	}
	rows, err := i.Query(
		fmt.Sprintf("SELECT username, owner, signature FROM identity WHERE %s", strings.Join(where, " AND ")),
		params...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var identities []*Identity
	for rows.Next() {
		var username, owner, signature string
		if err := rows.Scan(&username, &owner, &signature); err != nil {
			return nil, err
		}
		identities = append(identities, &Identity{
			Username:  username,
			Owner:     common.HexToAddress(owner),
			Signature: common.FromHex(signature),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return identities, nil
}

func (i *Index) CreateClaim(claim *Claim) error {
	return i.Update(func(tx *sql.Tx) error {
		var owner string
		err := i.QueryRow("SELECT owner FROM identity WHERE id = ?", claim.Issuer.String()).Scan(&owner)
		if err != nil {
			return err
		}
		if !verifyClaim(claim, common.HexToAddress(owner)) {
			return errors.New("identity: invalid claim")
		}
		_, err = tx.Exec(
			"INSERT INTO claim (id, issuer, subject, property, claim, signature) VALUES ($1, $2, $3, $4, $5, $6)",
			claim.ID().String(),
			claim.Issuer.String(),
			claim.Subject.String(),
			claim.Property,
			claim.Claim,
			hexutil.Encode(claim.Signature),
		)
		return err
	})
}

type ClaimFilter struct {
	Issuer   *string
	Subject  *string
	Property *string
	Claim    *string
}

func (i *Index) Claims(filter ClaimFilter) ([]*Claim, error) {
	where := make([]string, 0, 4)
	params := make([]interface{}, 0, 4)
	if filter.Issuer != nil {
		where = append(where, "issuer = ?")
		params = append(params, *filter.Issuer)
	}
	if filter.Subject != nil {
		where = append(where, "subject = ?")
		params = append(params, *filter.Subject)
	}
	if filter.Property != nil {
		where = append(where, "property = ?")
		params = append(params, *filter.Property)
	}
	if filter.Claim != nil {
		where = append(where, "claim = ?")
		params = append(params, *filter.Claim)
	}
	if len(where) == 0 {
		return nil, errors.New("missing issuer, subject, property or claim")
	}
	rows, err := i.Query(
		fmt.Sprintf("SELECT issuer, subject, property, claim, signature FROM claim WHERE %s", strings.Join(where, " AND ")),
		params...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var claims []*Claim
	for rows.Next() {
		var (
			issuer    string
			subject   string
			property  string
			claim     string
			signature string
		)
		if err := rows.Scan(
			&issuer,
			&subject,
			&property,
			&claim,
			&signature,
		); err != nil {
			return nil, err
		}
		claims = append(claims, &Claim{
			Issuer:    HexToID(issuer),
			Subject:   HexToID(subject),
			Property:  property,
			Claim:     claim,
			Signature: common.FromHex(signature),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return claims, nil
}

func verifyIdentity(identity *Identity) bool {
	id := identity.ID()
	return verifySignature(identity.Owner, id.Hash[:], identity.Signature)
}

func verifyClaim(claim *Claim, owner common.Address) bool {
	id := claim.ID()
	return verifySignature(owner, id[:], claim.Signature)
}

func verifySignature(owner common.Address, msg, signature []byte) bool {
	recoveredPub, err := crypto.Ecrecover(msg, signature)
	if err != nil {
		return false
	}
	pubKey := crypto.ToECDSAPub(recoveredPub)
	return crypto.PubkeyToAddress(*pubKey) == owner
}
