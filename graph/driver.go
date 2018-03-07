// This file is part of the go-kord library.
//
// Copyright (C) 2018 JAAK MUSIC LTD
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

package graph

import (
	"sync"

	"github.com/cayleygraph/cayley/graph"
	cayleysql "github.com/cayleygraph/cayley/graph/sql"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/swarm/storage"
	"github.com/kord-network/go-kord/db"
	"github.com/kord-network/go-kord/registry"
)

type Driver struct {
	name     string
	db       *db.Driver
	registry registry.Registry

	stores   map[string]graph.QuadStore
	storeMtx sync.Mutex
}

func NewDriver(name string, dpa *storage.DPA, registry registry.Registry, tmpDir string) *Driver {
	// create a Swarm backed SQLite database driver
	db := db.NewDriver(name, dpa, registry, tmpDir)

	// register the db driver as a Cayley SQL backend
	cayleysql.Register(name, db.GraphRegistration())

	// return a graph driver
	return &Driver{
		name:     name,
		db:       db,
		registry: registry,
		stores:   make(map[string]graph.QuadStore),
	}
}

// Create creates a new graph.
func (d *Driver) Create(name string) (common.Hash, error) {
	if err := graph.InitQuadStore(d.name, name, graph.Options{}); err != nil {
		return common.Hash{}, err
	}
	return d.Commit(name)
}

func (d *Driver) SetGraph(hash common.Hash, sig []byte) error {
	return d.registry.SetGraph(hash, sig)
}

func (d *Driver) Get(name string) (graph.QuadStore, error) {
	d.storeMtx.Lock()
	defer d.storeMtx.Unlock()
	if store, ok := d.stores[name]; ok {
		return store, nil
	}
	store, err := graph.NewQuadStore(d.name, name, graph.Options{})
	if err != nil {
		return nil, err
	}
	d.stores[name] = store
	return store, nil
}

func (d *Driver) Commit(name string) (common.Hash, error) {
	return d.db.Commit(name)
}
