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
	"bytes"
	"fmt"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipld-cbor"
	"github.com/multiformats/go-multihash"
	"github.com/whyrusleeping/cbor/go"
)

func Encode(properties Properties) (*Object, error) {
	p := make(map[string]interface{}, len(properties))
	for key, val := range properties {
		switch val.(type) {
		case string, []byte, *cid.Cid, map[string]string, map[string]*cid.Cid, []*cid.Cid, *Object:
			p[key] = val
		default:
			return nil, fmt.Errorf("meta: unsupported property value: %T", val)
		}
	}

	var buf bytes.Buffer
	enc := cbor.NewEncoder(&buf)
	enc.SetFilter(cbornode.EncoderFilter)
	if err := enc.Encode(p); err != nil {
		return nil, err
	}
	data := buf.Bytes()

	cid, err := cid.Prefix{
		Version:  1,
		Codec:    cid.DagCBOR,
		MhType:   multihash.SHA2_256,
		MhLength: -1,
	}.Sum(data)
	if err != nil {
		return nil, err
	}

	return NewObject(cid, data)
}

func MustEncode(properties Properties) *Object {
	obj, err := Encode(properties)
	if err != nil {
		panic(err)
	}
	return obj
}
