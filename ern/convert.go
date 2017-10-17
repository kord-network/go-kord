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

package ern

import (
	"io"

	"github.com/ipfs/go-cid"
	"github.com/meta-network/go-meta"
	"github.com/meta-network/go-meta/xml"
	"github.com/meta-network/go-meta/xmlschema"
)

// Converter converts DDEX ERN XML files into META objects.
type Converter struct {
	*metaxml.Converter
}

// NewConverter returns a Converter which stores META objects in the given META
// store.
func NewConverter(store *meta.Store) *Converter {
	return &Converter{metaxml.NewConverter(store)}
}

// ConvertERN converts the given ERN XML file into a META object graph with the
// given source and returns the CID of the graph's root META object.
func (c *Converter) ConvertERN(xml io.Reader, source string) (*cid.Cid, error) {
	// use the DDEX ERN/382 and AVS XML schemas as the JSON-LD context
	context := []*cid.Cid{
		xmlschema.DDEX_Ern382.Cid,
		xmlschema.DDEX_Avs.Cid,
	}
	obj, err := c.ConvertXML(xml, context, source)
	if err != nil {
		return nil, err
	}
	return obj.Cid(), nil
}
