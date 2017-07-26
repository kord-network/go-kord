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

package meta

import (
	"fmt"
	"strings"

	"github.com/ipfs/go-cid"
)

type ErrPathNotFound struct {
	Path []string
}

func (e ErrPathNotFound) Error() string {
	return fmt.Sprintf("meta: path not found: %s", strings.Join(e.Path, "/"))
}

type ErrInvalidCidVersion struct {
	Version uint64
}

func (e ErrInvalidCidVersion) Error() string {
	return fmt.Sprintf("meta: invalid CID version: %d", e.Version)
}

type ErrInvalidCodec struct {
	Codec uint64
}

func (e ErrInvalidCodec) Error() string {
	return fmt.Sprintf("meta: invalid CID codec: %x", e.Codec)
}

type ErrInvalidType struct {
	Type interface{}
}

func (e ErrInvalidType) Error() string {
	return fmt.Sprintf("meta: field @type is not a string (has type %T)", e.Type)
}

type ErrCidMismatch struct {
	Expected *cid.Cid
	Actual   *cid.Cid
}

func (e ErrCidMismatch) Error() string {
	return fmt.Sprintf("meta: CID mismatch, expected %q, got %q", e.Expected, e.Actual)
}
