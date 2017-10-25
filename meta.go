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
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	swarmapi "github.com/ethereum/go-ethereum/swarm/api"
	"github.com/ethereum/go-ethereum/swarm/storage"
	"github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipld-format"
	"github.com/lmars/go-ipld-cbor"
	multihash "github.com/multiformats/go-multihash"
)

// BaseObject contains the fields that all META objects should have.
type BaseObject struct {
	Source string `json:"@source"`
}

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
	if err := isValid(id, rawData); err != nil {
		panic(err)
	}
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
		if l, ok := v.(*format.Link); ok {
			v = l.Cid
		}
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

// Store provides storage for objects using a local Swarm chunk database
// and uses ENS to resolve names to hashes.
type Store struct {
	api       *swarmapi.Api
	dpa       *storage.DPA
	ens       ENS
	streamDir string
}

// NewStore returns a new store which maintains a local Swarm chunk databsse
// in the given directory and uses the given ENS to resolve names to hashes.
func NewStore(dir string, ens ENS) (*Store, error) {
	streamDir := filepath.Join(dir, "streams")
	if err := os.MkdirAll(streamDir, 0755); err != nil {
		return nil, err
	}
	dpa, err := storage.NewLocalDPA(dir, "", 20000000)
	if err != nil {
		return nil, err
	}
	dpa.Start()
	return &Store{
		api:       swarmapi.NewApi(dpa, ens),
		dpa:       dpa,
		ens:       ens,
		streamDir: streamDir,
	}, nil
}

// Close stops the underlying Swarm chunk database storage and retrieval loops.
func (s *Store) Close() {
	s.dpa.Stop()
}

// Get gets an object from the store.
func (s *Store) Get(cid *cid.Cid) (*Object, error) {
	hash, err := multihash.Decode(cid.Hash())
	if err != nil {
		return nil, err
	}
	reader := s.dpa.Retrieve(hash.Digest)
	size, err := reader.Size(nil)
	if err != nil {
		return nil, err
	}
	data := make([]byte, size)
	if _, err := io.ReadFull(reader, data); err != nil {
		return nil, err
	}
	return NewObject(cid, data)
}

const multihashSwarmCode = 0x30

func init() {
	multihash.Codes[multihashSwarmCode] = "swarm-hash-v1"
}

// Put encodes and stores object in the store.
func (s *Store) Put(v interface{}) (*Object, error) {
	data, err := encode(v)
	if err != nil {
		return nil, err
	}
	hash, err := s.api.Store(
		bytes.NewReader(data),
		int64(len(data)),
		&sync.WaitGroup{},
	)
	mhash, err := multihash.Encode(hash, multihashSwarmCode)
	if err != nil {
		return nil, err
	}
	cid := cid.NewCidV1(cid.DagCBOR, mhash)
	return NewObject(cid, data)
}

// MustPut is like Put but panics if v cannot be encoded or stored
func (s *Store) MustPut(v interface{}) *Object {
	obj, err := s.Put(v)
	if err != nil {
		panic(err)
	}
	return obj
}

// SwarmAPI returns the underlying Swarm Api instance to use when starting a
// HTTP server.
func (s *Store) SwarmAPI() *swarmapi.Api {
	return s.api
}

// StreamWriter returns a StreamWriter which writes CIDs to a local file.
func (s *Store) StreamWriter(name string) (*StreamWriter, error) {
	return NewStreamWriter(s.streamPath(name))
}

// StreamReader returns a StreamReader which reads CIDs from a local file.
func (s *Store) StreamReader(name string, opts ...StreamOpts) (*StreamReader, error) {
	return NewStreamReader(s.streamPath(name), opts...)
}

func (s *Store) streamPath(name string) string {
	return filepath.Join(s.streamDir, name)
}

// OpenIndex opens the META index with the given name by fetching it from
// Swarm to a temp file and opening it as a SQLite3 database.
func (s *Store) OpenIndex(name string) (index *Index, err error) {
	hash, err := s.ens.Resolve(name)
	if err == ErrNameNotExist {
		// if the name doesn't exist, create a new empty index (which
		// mimics the behaviour of opening a new SQLite3 file)
		hash, err = s.createIndex(name)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	reader := s.dpa.Retrieve(hash[:])
	size, err := reader.Size(nil)
	if err != nil {
		return nil, fmt.Errorf("index %s (%s) not found", name, hash)
	}
	tmp, err := ioutil.TempFile("", "meta-index")
	if err != nil {
		return nil, err
	}
	defer func() {
		tmp.Close()
		if err != nil {
			os.Remove(tmp.Name())
		}
	}()
	n, err := io.Copy(tmp, io.LimitReader(reader, size))
	if err != nil {
		return nil, err
	} else if n != size {
		return nil, io.ErrShortWrite
	}
	db, err := sql.Open("sqlite3", tmp.Name())
	if err != nil {
		return nil, err
	}
	return &Index{
		DB:    db,
		name:  name,
		store: s,
		path:  tmp.Name(),
	}, nil
}

// createIndex creates a new, empty index and points name at it.
func (s *Store) createIndex(name string) (common.Hash, error) {
	// store an empty file in Swarm
	hash, err := s.api.Store(bytes.NewReader(nil), 0, &sync.WaitGroup{})
	if err != nil {
		return common.Hash{}, err
	}
	// point name at the hash of the empty file
	if err := s.ens.SetContentHash(name, common.BytesToHash(hash)); err != nil {
		return common.Hash{}, err
	}
	return common.BytesToHash(hash), nil
}

// Index wraps an open SQLite3 database and provides an Update function which
// can be called to commit changes to the index and upload the result to Swarm.
type Index struct {
	*sql.DB

	name  string
	store *Store
	path  string
}

// Path returns the filesystem path of the SQLite3 database file.
func (i *Index) Path() string {
	return i.path
}

// Close closes the SQLite3 database and deletes the file.
func (i *Index) Close() error {
	defer os.Remove(i.path)
	return i.DB.Close()
}

// Update starts a transaction to update the index, passes it to the given
// function, then uploads the updated index to Swarm.
func (i *Index) Update(fn func(*sql.Tx) error) error {
	tx, err := i.Begin()
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	tmp, err := os.Open(i.path)
	if err != nil {
		return err
	}
	defer tmp.Close()
	info, err := tmp.Stat()
	if err != nil {
		return err
	}
	hash, err := i.store.dpa.Store(
		tmp,
		info.Size(),
		&sync.WaitGroup{},
		&sync.WaitGroup{},
	)
	if err != nil {
		return err
	}
	return i.store.ens.SetContentHash(i.name, common.BytesToHash(hash))
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

func isValid(cid *cid.Cid, data []byte) error {

	prefix := cid.Prefix()
	if prefix.Version != cidV1 {
		return ErrInvalidCidVersion{prefix.Version}
	}
	expectedCid, err := prefix.Sum(data)
	if err != nil {
		return err
	}
	if !cid.Equals(expectedCid) {
		return ErrCidMismatch{Expected: expectedCid, Actual: cid}
	}
	return nil
}

// NewBlock returns a new block.
func NewBlock(cid *cid.Cid, data []byte) (*Block, error) {

	prefix := cid.Prefix()

	if prefix.Version != cidV1 {
		return nil, ErrInvalidCidVersion{prefix.Version}
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
