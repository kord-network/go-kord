// This file is part of the go-meta library.
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

package testutil

import (
	"io/ioutil"
	"os"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/swarm/storage"
	"github.com/meta-network/go-meta/ens"
)

type TestDPA struct {
	*storage.DPA

	Dir string
}

func NewTestDPA() (*TestDPA, error) {
	dir, err := ioutil.TempDir("", "meta-testutil")
	if err != nil {
		return nil, err
	}
	localStore, err := storage.NewLocalStore(
		storage.MakeHashFunc("SHA3"),
		&storage.StoreParams{
			ChunkDbPath:   dir,
			DbCapacity:    5000000,
			CacheCapacity: 5000,
			Radius:        0,
		},
	)
	if err != nil {
		os.RemoveAll(dir)
		return nil, err
	}
	chunker := storage.NewTreeChunker(storage.NewChunkerParams())
	dpa := &storage.DPA{
		Chunker:    chunker,
		ChunkStore: localStore,
	}
	dpa.Start()
	return &TestDPA{dpa, dir}, nil
}

func (t *TestDPA) Cleanup() {
	t.Stop()
	os.RemoveAll(t.Dir)
}

type ENS struct {
	mtx    sync.Mutex
	hashes map[string]common.Hash
	subs   map[string]map[*ENSSubscription]struct{}
}

func NewTestENS() *ENS {
	return &ENS{
		hashes: make(map[string]common.Hash),
		subs:   make(map[string]map[*ENSSubscription]struct{}),
	}
}

func (e *ENS) Content(name string) (common.Hash, error) {
	e.mtx.Lock()
	defer e.mtx.Unlock()
	return e.hashes[name], nil
}

func (e *ENS) SetContent(name string, hash common.Hash) error {
	e.mtx.Lock()
	defer e.mtx.Unlock()
	e.hashes[name] = hash
	if subs, ok := e.subs[name]; ok {
		for sub := range subs {
			sub.updates <- hash
		}
	}
	return nil
}

func (e *ENS) SubscribeContent(name string, updates chan common.Hash) (ens.Subscription, error) {
	e.mtx.Lock()
	defer e.mtx.Unlock()
	subs, ok := e.subs[name]
	if !ok {
		subs = make(map[*ENSSubscription]struct{})
		e.subs[name] = subs
	}
	sub := &ENSSubscription{e, name, updates}
	subs[sub] = struct{}{}
	return sub, nil
}

type ENSSubscription struct {
	ens     *ENS
	name    string
	updates chan common.Hash
}

func (e *ENSSubscription) Close() error {
	e.ens.mtx.Lock()
	defer e.ens.mtx.Unlock()
	delete(e.ens.subs[e.name], e)
	return nil
}

func (e *ENSSubscription) Err() error {
	return nil
}
