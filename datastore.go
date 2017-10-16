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
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/sync/syncmap"

	"github.com/ethereum/go-ethereum/swarm/api/client"
	multihash "github.com/multiformats/go-multihash"
)

// Datastore represents storage for any key-value pair.
type Datastore interface {
	// put stores `data`.
	put(data []byte) (mhash multihash.Multihash, err error)
	// get retrieves the `value` named by `key`.
	//`key` should be the hex encoding of the hash part of the multihash returned from put
	get(key string) (value []byte, err error)
}

// MapDatastore uses a standard Go map for internal storage.
type MapDatastore struct {
	values syncmap.Map
}

// newMapDatastore constructs a MapDatastore
func newMapDatastore() (d *MapDatastore) {
	return &MapDatastore{
		values: syncmap.Map{},
	}
}

// put implements Datastore.put
func (d *MapDatastore) put(data []byte) (mhash multihash.Multihash, err error) {
	mhash, err = multihash.Sum(data, multihash.SHA2_256, -1)
	if err != nil {
		return nil, err
	}
	hash, err := multihash.Decode(mhash)
	if err != nil {
		return nil, err
	}
	d.values.Store(hex.EncodeToString(hash.Digest), data)
	return
}

// get implements Datastore.get
func (d *MapDatastore) get(key string) (value []byte, err error) {
	val, found := d.values.Load(key)
	if !found {
		return nil, fmt.Errorf("MapDatastore get: value for key %q not found ", key)
	}
	return val.([]byte), nil
}

const multihashSwarmCode = 0x30

func init() {
	multihash.Codes[multihashSwarmCode] = "swarm-hash-v1"
}

// SwarmDatastore struct
type SwarmDatastore struct {
	client *client.Client
}

// newSwarmDatastore returns a new swarm Datastore
func newSwarmDatastore(serverURL string) *SwarmDatastore {
	return &SwarmDatastore{
		client: client.NewClient(serverURL),
	}
}

// put stores the given value and return its hash
func (ds *SwarmDatastore) put(data []byte) (mhash multihash.Multihash, err error) {
	h, err := ds.client.UploadRaw(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, err
	}
	hash, err := hex.DecodeString(h)
	if err != nil {
		return nil, err
	}
	mhash, err = multihash.Encode(hash, multihashSwarmCode)
	return
}

// get returns the value for given key
func (ds *SwarmDatastore) get(key string) (value []byte, err error) {
	res, err := ds.client.DownloadRaw(key)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	return ioutil.ReadAll(res)
}

var ObjectKeySuffix = ".dsobject"

// FSDatastore use a file per key to store values.
type FSDatastore struct {
	path string
}

// newFSDatastore returns a new fs Datastore at given `path`
func newFSDatastore(path string) (*FSDatastore, error) {
	if !isDir(path) {
		return nil, fmt.Errorf("Failed to find directory at: %v (file? perms?)", path)
	}

	return &FSDatastore{path: path}, nil
}

// KeyFilename returns the filename associated with `key`
func (d *FSDatastore) KeyFilename(key string) string {
	return filepath.Join(d.path, key, ObjectKeySuffix)
}

// put stores the given value.
func (d *FSDatastore) put(data []byte) (mhash multihash.Multihash, err error) {
	mhash, err = multihash.Sum(data, multihash.SHA2_256, -1)
	if err != nil {
		return nil, err
	}
	hash, err := multihash.Decode(mhash)
	if err != nil {
		return nil, err
	}
	fn := d.KeyFilename(hex.EncodeToString(hash.Digest))
	// mkdirall above.
	err = os.MkdirAll(filepath.Dir(fn), 0755)
	if err != nil {
		return nil, err
	}
	return mhash, ioutil.WriteFile(fn, data, 0666)
}

// get returns the value for given key
func (d *FSDatastore) get(key string) (value []byte, err error) {
	fn := d.KeyFilename(key)
	if !isFile(fn) {
		return nil, fmt.Errorf("FSDatastore get: file %q not found ", key)
	}
	return ioutil.ReadFile(fn)
}

// isDir returns whether given path is a directory
func isDir(path string) bool {
	finfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return finfo.IsDir()
}

// isFile returns whether given path is a file
func isFile(path string) bool {
	finfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !finfo.IsDir()
}
