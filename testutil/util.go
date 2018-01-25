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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/swarm/storage"
	meta "github.com/meta-network/go-meta"
	"github.com/meta-network/go-meta/store"
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

type ENS struct{}

func (e *ENS) Content(name string) (common.Hash, error) {
	return common.Hash{}, nil
}

func NewTestSigner() (common.Address, meta.TxSigner, error) {
	key, err := crypto.GenerateKey()
	if err != nil {
		return common.Address{}, nil, err
	}
	address := crypto.PubkeyToAddress(key.PublicKey)
	return address, store.NewPrivateKeySigner(key), nil
}
