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
	"io"

	"github.com/ipfs/go-cid"
	"github.com/meta-network/go-meta"
)

// Converter converts CWR data from a CWR file to META
// objects.
type Converter struct {
	store *meta.Store
}

// NewConverter returns a Converter which reads data from the given CWR file
// and stores META object in the given META store.
func NewConverter(store *meta.Store) *Converter {

	return &Converter{
		store: store,
	}
}

// ConvertRegisteredWork reads all registeredWork , converts them to META
// objects, stores them in the META store and sends their CIDs to the given
// stream.
func (c *Converter) ConvertRegisteredWork(ctx context.Context, outStream chan *cid.Cid, cwrFileReader io.Reader, cwr2JsonPython string) error {
	// get all artists from the db
	registeredWorks, err := ParseCWRFile(cwrFileReader, cwr2JsonPython)
	if err != nil {
		return err
	}
	for _, registerdWork := range registeredWorks {
		// convert the registerdWork to a META object
		obj, err := meta.Encode(registerdWork)
		if err != nil {
			return err
		}

		if err := c.store.Put(obj); err != nil {
			return err
		}
		//send the object's CID to the output stream
		select {
		case outStream <- obj.Cid():
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}
