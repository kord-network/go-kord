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

package eidr

import (
	"io"

	"github.com/ipfs/go-cid"
	"github.com/meta-network/go-meta"
	metaxml "github.com/meta-network/go-meta/xml"
	"github.com/meta-network/go-meta/xmlschema"
)

// Converter converts DDEX ERN XML files into META objects.
type Converter struct {
	store *meta.Store
}

// NewConverter returns a Converter which stores META objects in the given META
// store.
func NewConverter(store *meta.Store) *Converter {
	return &Converter{
		store: store,
	}
}

func (c *Converter) ConvertEIDRXML(src io.Reader) (*cid.Cid, error) {
	context := []*cid.Cid{
		xmlschema.EIDR_common.Cid,
		xmlschema.EIDR_md.Cid,
	}
	obj, err := metaxml.EncodeXML(src, context, c.store.Put)
	if err != nil {
		return nil, err
	}

	return obj.Cid(), nil
}
