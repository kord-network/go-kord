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
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// Identity structure
type Identity struct {
	ID    string `json:"id"`
	Owner string `json:"owner"`
	Sig   string `json:"signature"`
}

// NewIdentity create and returns new Identity.
// It validate the ownership of the identity account using Ecrecover with a signature
// which is beeing sent by the owner.See https://github.com/ethereum/go-ethereum/wiki/Management-APIs#personal_sign
func NewIdentity(username string, owner common.Address, signature []byte) (identity *Identity, err error) {

	recoveredPub, err := crypto.Ecrecover(crypto.Keccak256([]byte(username)), signature)
	if err != nil {
		return nil, err
	}
	pubKey := crypto.ToECDSAPub(recoveredPub)
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)

	if owner != recoveredAddr {
		return nil, fmt.Errorf("NewIdentity: sender fail to prove its account ownership")
	}

	return &Identity{
		Owner: strings.ToLower(owner.String()),
		ID:    crypto.Keccak256Hash([]byte(username)).String(),
		Sig:   hex.EncodeToString(signature),
	}, nil
}
