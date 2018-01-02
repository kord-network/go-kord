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

package meta

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/cayleygraph/cayley/graph"
	"github.com/cayleygraph/cayley/quad"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tent/canonical-json-go"

	metasql "github.com/meta-network/go-meta/sql"
)

type Backend interface {
	Apply(tx *Tx) (common.Hash, error)
}

type Signer interface {
	SignHash(address common.Address, hash []byte) ([]byte, error)
}

type Client struct {
	graph.QuadStore

	Address common.Address
	Name    string
	Backend Backend
	Signer  Signer
	Driver  *metasql.Driver
}

func NewClient(address common.Address, name string, backend Backend, signer Signer, driver *metasql.Driver) (*Client, error) {
	if err := graph.InitQuadStore("meta", name, graph.Options{}); err != nil {
		return nil, err
	}
	store, err := graph.NewQuadStore("meta", name, graph.Options{})
	if err != nil {
		return nil, err
	}
	return &Client{
		QuadStore: store,
		Address:   address,
		Name:      name,
		Backend:   backend,
		Signer:    signer,
		Driver:    driver,
	}, nil
}

func (c *Client) ApplyDeltas(in []graph.Delta, opts graph.IgnoreOpts) error {
	deltas := make([]graph.Delta, len(in))
	for i, delta := range in {
		delta.Quad.Label = quad.String(c.Address.Hex())
		deltas[i] = delta
	}
	req := &Request{
		Name:   c.Name,
		Deltas: deltas,
	}
	sig, err := c.Signer.SignHash(c.Address, req.Hash())
	if err != nil {
		return err
	}
	tx := &Tx{
		Address: c.Address,
		Data:    req.Bytes(),
		Sig:     sig,
	}
	hash, err := c.Backend.Apply(tx)
	if err != nil {
		return err
	}
	return c.Driver.Update(c.Name, hash)
}

type State struct {
	mtx    sync.Mutex
	stores map[string]graph.QuadStore
	driver *metasql.Driver
}

func NewState(driver *metasql.Driver) *State {
	return &State{
		stores: make(map[string]graph.QuadStore),
		driver: driver,
	}
}

// Apply applies a transaction.
func (s *State) Apply(tx *Tx) (common.Hash, error) {
	// check the signature
	hash := crypto.Keccak256(tx.Data)
	pubKey, err := crypto.SigToPub(hash, tx.Sig)
	if err != nil {
		return common.Hash{}, err
	}
	address := crypto.PubkeyToAddress(*pubKey)
	if address != tx.Address {
		return common.Hash{}, errors.New("invalid signature")
	}

	// decode the request
	var req Request
	if err := json.Unmarshal(tx.Data, &req); err != nil {
		return common.Hash{}, err
	}

	// check all quads labelled with address
	for _, delta := range req.Deltas {
		label := delta.Quad.Label
		if label == nil {
			return common.Hash{}, errors.New("quad label must be set")
		}
		labelString, ok := label.Native().(string)
		if !ok {
			return common.Hash{}, fmt.Errorf("invalid label type: %T", label.Native())
		}
		if !common.IsHexAddress(labelString) {
			return common.Hash{}, fmt.Errorf("invalid address in quad label: %s", label)
		}
		addr := common.HexToAddress(labelString)
		if addr != address {
			return common.Hash{}, fmt.Errorf("invalid address in quad label: %s", addr)
		}
	}

	// ApplyDeltas to the store
	store, err := s.store(req.Name)
	if err != nil {
		return common.Hash{}, err
	}
	if err := store.ApplyDeltas(req.Deltas, graph.IgnoreOpts{}); err != nil {
		return common.Hash{}, err
	}

	// TODO: rollback ApplyDeltas if saving fails
	return s.driver.Save(req.Name)
}

func (s *State) Close() {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	for _, store := range s.stores {
		store.Close()
	}
}

func (s *State) store(name string) (graph.QuadStore, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if store, ok := s.stores[name]; ok {
		return store, nil
	}
	store, err := graph.NewQuadStore("meta", name, graph.Options{})
	if err != nil {
		return nil, err
	}
	s.stores[name] = store
	return store, nil
}

type Request struct {
	Name   string        `json:"name"`
	Deltas []graph.Delta `json:"deltas"`
}

func (r *Request) Bytes() []byte {
	data, _ := cjson.Marshal(r)
	return data
}

func (r *Request) Hash() []byte {
	return crypto.Keccak256(r.Bytes())
}

// Tx is a transaction which can be applied to the META state.
type Tx struct {
	Address common.Address
	Data    []byte
	Sig     []byte
}

type txJSON struct {
	Address string `json:"address"`
	Data    string `json:"data"`
	Sig     string `json:"sig"`
}

func (tx *Tx) MarshalJSON() ([]byte, error) {
	return json.Marshal(txJSON{
		Address: tx.Address.Hex(),
		Data:    hexutil.Encode(tx.Data),
		Sig:     hexutil.Encode(tx.Sig),
	})
}

func (tx *Tx) UnmarshalJSON(data []byte) error {
	var v txJSON
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	txData, err := hexutil.Decode(v.Data)
	if err != nil {
		return err
	}
	txSig, err := hexutil.Decode(v.Sig)
	if err != nil {
		return err
	}
	*tx = Tx{
		Address: common.HexToAddress(v.Address),
		Data:    txData,
		Sig:     txSig,
	}
	return nil
}
