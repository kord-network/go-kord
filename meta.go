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

	"github.com/cayleygraph/cayley/graph"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/tent/canonical-json-go"
)

type TxSigner interface {
	SignTx(address common.Address, tx *Tx) (*SignedTx, error)
}

type ENS interface {
	Content(name string) (common.Hash, error)
}

type Tx struct {
	Name   string        `json:"name"`
	Deltas []graph.Delta `json:"deltas"`
}

func (tx *Tx) Bytes() []byte {
	data, _ := cjson.Marshal(tx)
	return data
}

// SignedTx is a transaction which can be applied to the META state.
type SignedTx struct {
	Address   common.Address
	Tx        []byte
	Signature []byte
}

type signedTxJSON struct {
	Address   string `json:"address"`
	Tx        string `json:"tx"`
	Signature string `json:"signature"`
}

func (tx *SignedTx) MarshalJSON() ([]byte, error) {
	return json.Marshal(signedTxJSON{
		Address:   tx.Address.Hex(),
		Tx:        hexutil.Encode(tx.Tx),
		Signature: hexutil.Encode(tx.Signature),
	})
}

func (tx *SignedTx) UnmarshalJSON(data []byte) error {
	var v signedTxJSON
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	txData, err := hexutil.Decode(v.Tx)
	if err != nil {
		return err
	}
	txSig, err := hexutil.Decode(v.Signature)
	if err != nil {
		return err
	}
	*tx = SignedTx{
		Address:   common.HexToAddress(v.Address),
		Tx:        txData,
		Signature: txSig,
	}
	return nil
}
