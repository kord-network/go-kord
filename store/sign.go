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
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	meta "github.com/meta-network/go-meta"
)

func NewPrivateKeySigner(key *ecdsa.PrivateKey) meta.TxSigner {
	return &privateKeySigner{key}
}

type privateKeySigner struct {
	key *ecdsa.PrivateKey
}

func (p *privateKeySigner) SignTx(address common.Address, tx *meta.Tx) (*meta.SignedTx, error) {
	if address != crypto.PubkeyToAddress(p.key.PublicKey) {
		return nil, fmt.Errorf("unknown address: %s", address)
	}
	data := tx.Bytes()
	hash := crypto.Keccak256(data)
	sig, err := crypto.Sign(hash, p.key)
	if err != nil {
		return nil, err
	}
	return &meta.SignedTx{
		Address:   address,
		Tx:        data,
		Signature: sig,
	}, nil
}
