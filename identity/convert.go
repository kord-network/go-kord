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
	"github.com/ipfs/go-cid"
	"github.com/meta-network/go-meta"
)

// Converter converts Identity data in meta object
// objects.
type Converter struct {
	store *meta.Store
}

// NewConverter returns a Converter which reads data from the given identity
// and stores META object in the given META store.
func NewConverter(store *meta.Store) *Converter {
	return &Converter{
		store: store,
	}
}

// ConvertIdentity converts the given Identity into a META object and
// returns the CID of the META object.
func (c *Converter) ConvertIdentity(identity *Identity) (*cid.Cid, error) {
	obj, err := c.store.Put(identity)
	if err != nil {
		return nil, err
	}
	return obj.Cid(), nil
}

// ConvertClaim converts the given Claim into a META object and
// returns the CID of the META object.
func (c *Converter) ConvertClaim(claim *Claim) (*cid.Cid, error) {
	obj, err := c.store.Put(claim)
	if err != nil {
		return nil, err
	}
	return obj.Cid(), nil
}
