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

	"github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-ipld-cbor"
	"github.com/ipfs/go-ipld-format"
)

// Object is a META object which uses IPLD DAG CBOR as the byte representation,
// and IPLD CID as the object identifier.
type Object struct {
	typ   string
	block blocks.Block
	node  *cbornode.Node
}

// NewObject returns an Object represented by an IPLD CID and the IPLD DAG CBOR
// byte representation of the object.
func NewObject(id *cid.Cid, rawData []byte) (*Object, error) {
	block, err := NewBlock(id, rawData)
	if err != nil {
		return nil, err
	}
	return NewObjectFromBlock(block)
}

// NewObjectFromBlock returns an Object represented by an IPFS block containing
// the IPLD DAG CBOR byte representation of the object.
func NewObjectFromBlock(block *Block) (*Object, error) {
	obj := &Object{block: block}

	if block.Codec() != cid.DagCBOR {
		return nil, ErrInvalidCodec{block.Codec()}
	}

	node, err := cbornode.DecodeBlock(block)
	if err != nil {
		return nil, err
	}
	obj.node = node

	if typ, _, err := obj.node.Resolve([]string{"@type"}); err == nil {
		typString, ok := typ.(string)
		if !ok {
			return nil, ErrInvalidType{typ}
		}
		obj.typ = typString
	}

	return obj, nil
}

// MustObject is like NewObject but panics if the given CID and raw bytes do
// not represent a valid Object.
func MustObject(id *cid.Cid, rawData []byte) *Object {
	obj, err := NewObject(id, rawData)
	if err != nil {
		panic(err)
	}
	return obj
}

// Cid returns the object's CID.
func (o *Object) Cid() *cid.Cid {
	return o.block.Cid()
}

// RawData returns the IPLD DAG CBOR representation of the object.
func (o *Object) RawData() []byte {
	return o.block.RawData()
}

// Type returns the object's type which is stored in the @type property.
func (o *Object) Type() string {
	return o.typ
}

// GetString looks up the property with the given key, type asserts it as a
// string and returns it.
func (o *Object) GetString(key string) (string, error) {
	v, err := o.Get(key)
	if err != nil {
		return "", err
	}
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("key %q has type %T, not string", key, v)
	}
	return s, nil
}

// GetLink looks up the property with the given key, type asserts it as a
// link and returns it.
func (o *Object) GetLink(key string) (*format.Link, error) {
	v, err := o.Get(key)
	if err != nil {
		return nil, err
	}
	l, ok := v.(*format.Link)
	if !ok {
		return nil, fmt.Errorf("key %q has type %T, not *format.Link", key, v)
	}
	return l, nil
}

// GetList looks up the property with the given key, type asserts it as a
// generic list and returns it.
func (o *Object) GetList(key string) ([]interface{}, error) {
	v, err := o.Get(key)
	if err != nil {
		return nil, err
	}
	l, ok := v.([]interface{})
	if !ok {
		return nil, fmt.Errorf("key %q has type %T, not []interface{}", key, v)
	}
	return l, nil
}

// Get returns the property with the given key.
func (o *Object) Get(key string) (interface{}, error) {
	v, rest, err := o.node.Resolve([]string{key})
	if err != nil {
		return nil, fmt.Errorf("error getting key %q: %s", key, err)
	} else if len(rest) > 0 {
		return nil, fmt.Errorf("error getting key %q: cannot resolve through link", key)
	}
	return v, nil
}

// MarshalJSON implements the json.Marshaler interface by encoding the
// underlying CBOR node.
func (o *Object) MarshalJSON() ([]byte, error) {
	return o.node.MarshalJSON()
}

// Graph is used to traverse an object graph using a store and starting from
// a particular root object.
type Graph struct {
	store *Store
	root  *Object
}

// NewGraph returns a new Graph
func NewGraph(store *Store, root *Object) *Graph {
	return &Graph{store, root}
}

// Root returns the root object of the graph
func (g *Graph) Root() *Object {
	return g.root
}

// Get gets the object at the given path.
func (g *Graph) Get(path ...string) (interface{}, error) {
	if len(path) == 1 && path[0] == "" {
		return g.root, nil
	}
	v, rest, err := g.root.node.Resolve(path)
	if err != nil {
		if err == cbornode.ErrNoSuchLink {
			err = ErrPathNotFound{path}
		}
		return nil, err
	}
	if len(rest) == 0 {
		return v, nil
	}

	link, ok := v.(*format.Link)
	if !ok {
		return nil, fmt.Errorf("meta: expected link object, got %T", v)
	}

	obj, err := g.store.Get(link.Cid)
	if err != nil {
		return nil, err
	}

	return NewGraph(g.store, obj).Get(rest...)
}

// Store provides storage for objects.
type Store struct {
	store datastore.Datastore
}

// NewStore returns a new Store which uses an underlying datastore.
func NewStore(store datastore.Datastore) *Store {
	return &Store{store}
}

// Get gets an object from the store.
func (s *Store) Get(cid *cid.Cid) (*Object, error) {
	data, err := s.store.Get(s.key(cid))
	if err != nil {
		return nil, err
	}
	return NewObject(cid, data.([]byte))
}

// Put stores an object in the store.
func (s *Store) Put(obj *Object) error {
	return s.store.Put(s.key(obj.Cid()), obj.RawData())
}

// key generates the key to use to store and retrieve the object with the
// given CID.
func (s *Store) key(cid *cid.Cid) datastore.Key {
	return datastore.NewKey(cid.String())
}

// cidV1 is the number which identifies a CID as being CIDv1.
//
// TODO: move this to the github.com/ipfs/go-cid.
const cidV1 = 1

// Block wraps a raw byte slice and validates it against a CID.
type Block struct {
	blocks.BasicBlock

	prefix *cid.Prefix
}

// NewBlock returns a new block.
func NewBlock(cid *cid.Cid, data []byte) (*Block, error) {
	prefix := cid.Prefix()

	if prefix.Version != cidV1 {
		return nil, ErrInvalidCidVersion{prefix.Version}
	}

	expectedCid, err := prefix.Sum(data)
	if err != nil {
		return nil, err
	}
	if !cid.Equals(expectedCid) {
		return nil, ErrCidMismatch{Expected: expectedCid, Actual: cid}
	}

	block, err := blocks.NewBlockWithCid(data, cid)
	if err != nil {
		return nil, err
	}

	return &Block{
		BasicBlock: *block,
		prefix:     &prefix,
	}, nil
}

// Codec returns the codec of the underlying data (e.g. IPLD DAG CBOR).
func (b *Block) Codec() uint64 {
	return b.prefix.Codec
}
