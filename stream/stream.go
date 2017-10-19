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

package stream

import "github.com/ipfs/go-cid"

// Stream provides readers and writers for an underlying META stream.
type Stream interface {
	NewReader() (StreamReader, error)
	NewWriter() (StreamWriter, error)
}

// StreamWriter writes CIDs to an underlying META stream.
type StreamWriter interface {
	Write(...*cid.Cid) error
	Close() error
}

// StreamReader reads CIDs from an underlying META stream.
type StreamReader interface {
	Ch() <-chan *cid.Cid
	Close() error
	Err() error
}
