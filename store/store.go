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

package store

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cayleygraph/cayley/graph"
	"github.com/cayleygraph/cayley/quad"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	meta "github.com/meta-network/go-meta"
	"github.com/meta-network/go-meta/db"
)

type Client interface {
	ApplyTransaction(name string, tx *meta.SignedTx) (common.Hash, error)
}

type ClientStore struct {
	graph.QuadStore

	address common.Address
	signer  meta.TxSigner
	client  Client
	name    string
}

func NewClientStore(address common.Address, signer meta.TxSigner, client Client, name string) (*ClientStore, error) {
	qs, err := graph.NewQuadStore("meta", name, graph.Options{})
	if err != nil {
		return nil, err
	}
	return &ClientStore{
		QuadStore: qs,
		address:   address,
		signer:    signer,
		client:    client,
		name:      name,
	}, nil
}

func (c *ClientStore) ApplyDeltas(in []graph.Delta, opts graph.IgnoreOpts) error {
	// label each quad with the client's address
	deltas := make([]graph.Delta, len(in))
	for i, delta := range in {
		delta.Quad.Label = quad.String(c.address.Hex())
		deltas[i] = delta
	}

	// create a signed transaction
	tx, err := c.signer.SignTx(c.address, &meta.Tx{
		Name:   c.name,
		Deltas: deltas,
	})
	if err != nil {
		return err
	}

	// apply the signed transaction
	hash, err := c.client.ApplyTransaction(c.name, tx)
	if err != nil {
		return err
	}

	// update open databases with the new hash
	return db.Update(c.name, hash)
}

type ServerStore struct {
	graph.QuadStore

	name string
}

func NewServerStore(name string) (*ServerStore, error) {
	qs, err := graph.NewQuadStore("meta", name, graph.Options{})
	if err != nil {
		return nil, err
	}
	return &ServerStore{
		QuadStore: qs,
		name:      name,
	}, nil
}

func (s *ServerStore) HandleTx(signedTx *meta.SignedTx) (common.Hash, error) {
	// check the signature
	hash := crypto.Keccak256(signedTx.Tx)
	pubKey, err := crypto.SigToPub(hash, signedTx.Signature)
	if err != nil {
		return common.Hash{}, err
	}
	address := crypto.PubkeyToAddress(*pubKey)
	if address != signedTx.Address {
		return common.Hash{}, errors.New("invalid signature")
	}

	// decode the transaction
	var tx meta.Tx
	if err := json.Unmarshal(signedTx.Tx, &tx); err != nil {
		return common.Hash{}, err
	}

	// check the database name
	if tx.Name != s.name {
		return common.Hash{}, fmt.Errorf("invalid database name: %s", tx.Name)
	}

	// check all quads labelled with address
	for _, delta := range tx.Deltas {
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
	if err := s.ApplyDeltas(tx.Deltas, graph.IgnoreOpts{}); err != nil {
		return common.Hash{}, err
	}

	// TODO: rollback ApplyDeltas if saving fails
	return db.Commit(tx.Name)
}
