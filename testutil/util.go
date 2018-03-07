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

package testutil

import (
	"io/ioutil"
	"os"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/swarm/storage"
	"github.com/kord-network/go-kord/registry"
)

type TestDPA struct {
	*storage.DPA

	Dir string
}

func NewTestDPA() (*TestDPA, error) {
	dir, err := ioutil.TempDir("", "kord-testutil")
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

type Registry struct {
	mtx    sync.Mutex
	hashes map[common.Address]common.Hash
	subs   map[common.Address]map[*RegistrySubscription]struct{}
}

func NewTestRegistry() *Registry {
	return &Registry{
		hashes: make(map[common.Address]common.Hash),
		subs:   make(map[common.Address]map[*RegistrySubscription]struct{}),
	}
}

func (r *Registry) Graph(kordID common.Address) (common.Hash, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	return r.hashes[kordID], nil
}

func (r *Registry) SetGraph(hash common.Hash, sig []byte) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	pub, err := crypto.SigToPub(hash[:], sig)
	if err != nil {
		return err
	}
	kordID := crypto.PubkeyToAddress(*pub)
	r.hashes[kordID] = hash
	if subs, ok := r.subs[kordID]; ok {
		for sub := range subs {
			sub.updates <- hash
		}
	}
	return nil
}

func (r *Registry) SubscribeGraph(kordID common.Address, updates chan common.Hash) (registry.Subscription, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	subs, ok := r.subs[kordID]
	if !ok {
		subs = make(map[*RegistrySubscription]struct{})
		r.subs[kordID] = subs
	}
	sub := &RegistrySubscription{r, kordID, updates}
	subs[sub] = struct{}{}
	return sub, nil
}

type RegistrySubscription struct {
	registry *Registry
	kordID   common.Address
	updates  chan common.Hash
}

func (r *RegistrySubscription) Close() error {
	r.registry.mtx.Lock()
	defer r.registry.mtx.Unlock()
	delete(r.registry.subs[r.kordID], r)
	return nil
}

func (r *RegistrySubscription) Err() error {
	return nil
}
