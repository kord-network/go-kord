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
	"errors"

	"github.com/ipfs/go-cid"
	"github.com/lmars/cbor/go"
	"github.com/lmars/go-ipld-cbor"
	"github.com/multiformats/go-multihash"
)

// Encode returns the META object encoding of v.
func Encode(v interface{}) (*Object, error) {
	var buf bytes.Buffer
	enc := cbor.NewEncoder(&buf)
	enc.SetFilter(cbornode.EncoderFilter)
	if err := enc.Encode(v); err != nil {
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

// MustEncode is like Encode but panics if v cannot be encoded.
func MustEncode(v interface{}) *Object {
	obj, err := Encode(v)
	if err != nil {
		panic(err)
	}
	return obj
}

// Decode decodes the META object into the value pointed to by v.
func (o *Object) Decode(v interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(r.(string))
		}
	}()
	dec := cbor.NewDecoder(bytes.NewReader(o.RawData()))
	dec.TagDecoders[cbornode.CBORTagLink] = &cbornode.IpldLinkDecoder{}
	return dec.Decode(v)
}
