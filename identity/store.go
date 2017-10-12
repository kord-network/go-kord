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

package identity

import (
	"fmt"
	"sync"
)

type Store interface {
	Load(name string) (*Identity, error)
	Save(*Identity) error
}

type MemoryStore struct {
	mtx        sync.RWMutex
	identities map[string]*Identity
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		identities: make(map[string]*Identity),
	}
}

func (m *MemoryStore) Load(name string) (*Identity, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	i, ok := m.identities[name]
	if !ok {
		return nil, fmt.Errorf("identity not found: %s", name)
	}
	return i, nil
}

func (m *MemoryStore) Save(identity *Identity) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.identities[identity.Name] = identity
	return nil
}
